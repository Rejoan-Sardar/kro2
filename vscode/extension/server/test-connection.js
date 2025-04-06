#!/usr/bin/env node

const net = require('net');
const client = new net.Socket();

// LSP message for initialize request
const initMessage = {
  jsonrpc: '2.0',
  id: 1,
  method: 'initialize',
  params: {
    processId: process.pid,
    clientInfo: {
      name: 'Test Client',
      version: '1.0.0'
    },
    rootUri: null,
    capabilities: {
      textDocument: {
        synchronization: {
          dynamicRegistration: false,
          willSave: false,
          willSaveWaitUntil: false,
          didSave: false
        },
        completion: {
          dynamicRegistration: false,
          completionItem: {
            snippetSupport: false,
            commitCharactersSupport: false,
            documentationFormat: ['plaintext'],
            deprecatedSupport: false,
            preselectSupport: false
          },
          completionItemKind: {
            valueSet: []
          },
          contextSupport: false
        }
      }
    }
  }
};

// Format the message according to LSP specification
function formatLSPMessage(message) {
  const content = JSON.stringify(message);
  const contentLength = Buffer.byteLength(content, 'utf8');
  const header = `Content-Length: ${contentLength}\r\n\r\n`;
  return header + content;
}

client.connect(5001, '127.0.0.1', () => {
  console.log('Connected to LSP server at 127.0.0.1:5001');
  
  // Send initialize request
  const message = formatLSPMessage(initMessage);
  client.write(message);
  console.log('Sent initialize request');
});

// Buffer to accumulate data from the server
let buffer = '';
let contentLength = -1;

client.on('data', (data) => {
  buffer += data.toString();
  
  // Process complete messages
  while (true) {
    if (contentLength < 0) {
      // Looking for Content-Length header
      const match = buffer.match(/Content-Length: (\d+)\r\n\r\n/);
      if (!match) break;
      
      contentLength = parseInt(match[1], 10);
      buffer = buffer.substring(match[0].length);
    }
    
    // Check if we have a complete message
    if (buffer.length >= contentLength) {
      const message = buffer.substring(0, contentLength);
      buffer = buffer.substring(contentLength);
      contentLength = -1;
      
      // Parse and display the message
      try {
        const parsed = JSON.parse(message);
        console.log('Received response:');
        console.log(JSON.stringify(parsed, null, 2));
        
        // Exit after receiving the initialize response
        if (parsed.id === 1 && parsed.result) {
          console.log('LSP server responded successfully to initialize request');
          client.end();
          process.exit(0);
        }
      } catch (e) {
        console.error('Error parsing JSON response:', e);
        client.end();
        process.exit(1);
      }
    } else {
      // Need more data
      break;
    }
  }
});

client.on('close', () => {
  console.log('Connection closed');
});

client.on('error', (err) => {
  console.error('Connection error:', err.message);
  process.exit(1);
});

// Set a timeout to exit if no response is received
setTimeout(() => {
  console.error('Timeout: No response received from the server');
  client.end();
  process.exit(1);
}, 5000);