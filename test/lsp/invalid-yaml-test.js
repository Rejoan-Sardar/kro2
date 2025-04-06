const net = require('net');
const fs = require('fs');
const path = require('path');

// Configuration
const SERVER_HOST = '127.0.0.1';
const SERVER_PORT = 5001;
const TEST_FILE_PATH = path.resolve(__dirname, 'invalid-yaml.yaml');

// Create the test file with an invalid YAML but valid CEL expressions
const invalidContent = `apiVersion: kro.run/v1alpha1
kind: ResourceGraphDefinition
metadata:
  name: web-application
spec:
  parameters:
    appName:
      type: string
      description: "The name of the application"
      required: true
  resources:
    webserver:
      type: aws.ec2.Instance
      template:
      name: "{{ params.appName }}-webserver"
        instanceType: "t3.micro"
        tags:
          - key: "Name"
            value: "{{ params.appName }}"`;

fs.writeFileSync(TEST_FILE_PATH, invalidContent);
console.log(`Created invalid test file at ${TEST_FILE_PATH}`);

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

        // Wait for diagnostics - our file has invalid YAML indentation
        // so we should receive diagnostics for it
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

// Buffer to handle partial messages
let messageBuffer = '';

// Handle incoming messages
function handleMessage(client, data) {
    // Append the new data to our buffer
    messageBuffer += data.toString();
    
    // Process complete messages from the buffer
    while (messageBuffer.length > 0) {
        // Look for Content-Length header
        const headerMatch = messageBuffer.match(/Content-Length: (\d+)\r\n\r\n/);
        if (!headerMatch) {
            // If we don't have a complete header yet, wait for more data
            return;
        }

        const headerEnd = messageBuffer.indexOf('\r\n\r\n') + 4;
        const contentLength = parseInt(headerMatch[1]);
        
        // Check if we have enough data for the complete message
        if (messageBuffer.length < headerEnd + contentLength) {
            // Need more data, wait for the next chunk
            return;
        }
        
        // Extract the message content
        const content = messageBuffer.substring(headerEnd, headerEnd + contentLength);
        
        // Remove this message from the buffer
        messageBuffer = messageBuffer.substring(headerEnd + contentLength);
        
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
                
                // We can disconnect after receiving diagnostics - we don't need to wait
                setTimeout(() => {
                    client.end();
                }, 500);
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
}

// Send JSON-RPC message with Content-Length header
function sendJsonRpcMessage(client, message) {
    const content = JSON.stringify(message);
    const contentLength = Buffer.byteLength(content, 'utf8');
    const header = `Content-Length: ${contentLength}\r\n\r\n`;
    
    console.log(`Sending ${message.method || 'response'}...`);
    client.write(header + content);
}