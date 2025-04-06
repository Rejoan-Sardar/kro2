const net = require('net');
const fs = require('fs');
const path = require('path');

// Configuration
const SERVER_HOST = '127.0.0.1';
const SERVER_PORT = 5001;
const TEST_FILE_PATH = path.resolve(__dirname, 'sample-with-cel.yaml');

// Create the test file if it doesn't exist
const sampleContent = `apiVersion: kro.run/v1alpha1
kind: ResourceGraphDefinition
metadata:
  name: web-application
spec:
  parameters:
    appName:
      type: string
      description: "The name of the application"
      required: true
    region:
      type: string
      description: "The AWS region to deploy to"
      default: "us-west-2"
  resources:
    webserver:
      type: aws.ec2.Instance
      template:
        name: "{{ params.appName }}-webserver"
        instanceType: "t3.micro"
        region: "{{ params.region }}"
    database:
      type: aws.rds.Instance
      template:
        name: "{{ params.appName }}-db"
        engine: "mysql"
        region: "{{ params.region }}"
  relations:
    webserverToDb:
      type: NetworkConnection
      from: webserver
      to: database
      template:
        port: 3306`;

fs.writeFileSync(TEST_FILE_PATH, sampleContent);
console.log(`Created test file at ${TEST_FILE_PATH}`);

// Unique ID for JSON-RPC requests
let messageId = 0;

// Connect to the LSP server
const client = net.createConnection({ host: SERVER_HOST, port: SERVER_PORT }, () => {
    console.log('Connected to LSP server');

    // Send initialize request
    const initializeParams = {
        processId: process.pid,
        clientInfo: {
            name: "test-client"
        },
        rootUri: `file://${process.cwd()}`,
        capabilities: {
            textDocument: {
                synchronization: {
                    didSave: true,
                    dynamicRegistration: true
                }
            }
        }
    };

    sendJsonRpcMessage(client, {
        jsonrpc: '2.0',
        id: messageId++,
        method: 'initialize',
        params: initializeParams
    });

    // Read the file content
    const fileContent = fs.readFileSync(TEST_FILE_PATH, 'utf8');
    const fileUri = `file://${TEST_FILE_PATH}`;

    // Send document open notification after server is initialized
    client.once('data', (data) => {
        console.log('Received initialize response');

        // Send initialized notification
        sendJsonRpcMessage(client, {
            jsonrpc: '2.0',
            method: 'initialized',
            params: {}
        });

        // Send didOpen notification for the test file
        sendJsonRpcMessage(client, {
            jsonrpc: '2.0',
            method: 'textDocument/didOpen',
            params: {
                textDocument: {
                    uri: fileUri,
                    languageId: 'yaml',
                    version: 1,
                    text: fileContent
                }
            }
        });

        console.log('Sent didOpen notification');

        // Wait for diagnostics
        setTimeout(() => {
            // Request completion at a position in the file (after "template:")
            sendJsonRpcMessage(client, {
                jsonrpc: '2.0',
                id: messageId++,
                method: 'textDocument/completion',
                params: {
                    textDocument: { uri: fileUri },
                    position: { line: 18, character: 16 } // After "template:"
                }
            });

            // Request hover information for a cel expression
            sendJsonRpcMessage(client, {
                jsonrpc: '2.0',
                id: messageId++,
                method: 'textDocument/hover',
                params: {
                    textDocument: { uri: fileUri },
                    position: { line: 19, character: 24 } // Inside "{{ params.appName }}"
                }
            });
        }, 1000);
    });
});

// Handle incoming data
client.on('data', (data) => {
    handleMessage(client, data);
});

// Handle connection errors
client.on('error', (err) => {
    console.error('Connection error:', err);
});

// Handle connection close
client.on('close', () => {
    console.log('Connection closed');
});

// Parse headers for Content-Length
function parseHeaders(headerText) {
    const headers = {};
    const headerLines = headerText.split('\r\n');
    for (const line of headerLines) {
        if (line.trim() === '') continue;
        const [key, value] = line.split(':');
        headers[key.trim()] = value.trim();
    }
    return headers;
}

// Handle incoming messages
function handleMessage(client, data) {
    const message = data.toString();
    
    // Check for Content-Length header pattern
    const headerMatch = message.match(/Content-Length: (\d+)\r\n\r\n/);
    if (!headerMatch) {
        console.log('Received incomplete message:', message);
        return;
    }

    const headerEnd = message.indexOf('\r\n\r\n') + 4;
    const contentLength = parseInt(headerMatch[1]);
    
    // Extract headers and content
    const headers = parseHeaders(message.substring(0, headerEnd));
    const content = message.substring(headerEnd, headerEnd + contentLength);
    
    try {
        const jsonContent = JSON.parse(content);
        
        // Handle different response types
        if (jsonContent.method === 'textDocument/publishDiagnostics') {
            console.log('Received diagnostics:');
            const diagnostics = jsonContent.params.diagnostics;
            if (diagnostics.length === 0) {
                console.log('No diagnostics reported (document is valid)');
            } else {
                diagnostics.forEach(diag => {
                    console.log(`- [${diag.severity}] ${diag.message} (${diag.source}) at line ${diag.range.start.line + 1}, col ${diag.range.start.character + 1}`);
                });
            }
        } else if (jsonContent.id !== undefined) {
            // This is a response to a request
            console.log(`Received response for request ${jsonContent.id}:`, JSON.stringify(jsonContent, null, 2));
        } else if (jsonContent.method) {
            // This is a notification
            console.log(`Received notification: ${jsonContent.method}`);
        }
    } catch (e) {
        console.error('Error parsing JSON:', e);
        console.log('Content:', content);
    }
}

// Send JSON-RPC message with Content-Length header
function sendJsonRpcMessage(client, message) {
    const content = JSON.stringify(message);
    const contentLength = Buffer.byteLength(content, 'utf8');
    const header = `Content-Length: ${contentLength}\r\n\r\n`;
    
    console.log(`Sending ${message.method || 'response'}...`);
    client.write(header + content);
}