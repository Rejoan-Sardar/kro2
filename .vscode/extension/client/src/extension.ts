import * as path from 'path';
import * as vscode from 'vscode';
import * as net from 'net';
import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  StreamInfo
} from 'vscode-languageclient/node';

let client: LanguageClient;

export function activate(context: vscode.ExtensionContext) {
  // Output channel for logging
  const outputChannel = vscode.window.createOutputChannel('Kro Language Server');
  
  // Server settings from configuration
  const config = vscode.workspace.getConfiguration('kroLanguageServer');
  const debugEnabled = config.get<boolean>('debug', false);
  
  // Use TCP socket for communication with server
  const serverHost = '127.0.0.1';
  const serverPort = 5001;

  // Define server options - connect to LSP server over TCP
  const serverOptions: ServerOptions = () => {
    return new Promise<StreamInfo>((resolve, reject) => {
      const socket = new net.Socket();
      
      socket.connect(serverPort, serverHost, () => {
        if (debugEnabled) {
          outputChannel.appendLine(`Connected to Kro Language Server at ${serverHost}:${serverPort}`);
        }
        
        resolve({
          reader: socket,
          writer: socket
        });
      });
      
      socket.on('error', (err) => {
        outputChannel.appendLine(`Socket error: ${err.message}`);
        reject(err);
      });
    });
  };
  
  // Options to control the language client
  const clientOptions: LanguageClientOptions = {
    // Register the server for YAML/Kro documents
    documentSelector: [
      { scheme: 'file', language: 'yaml' },
      { scheme: 'file', language: 'kro' }
    ],
    outputChannel: outputChannel,
    synchronize: {
      // Synchronize the setting section 'kroLanguageServer' to the server
      configurationSection: 'kroLanguageServer'
    }
  };

  // Create the language client
  client = new LanguageClient(
    'kroLanguageServer',
    'Kro Language Server',
    serverOptions,
    clientOptions
  );

  // Add status bar item
  const statusBarItem = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Right, 100);
  statusBarItem.text = "$(sync) Kro LS";
  statusBarItem.tooltip = "Kro Language Server";
  statusBarItem.show();
  context.subscriptions.push(statusBarItem);

  // Log start message
  outputChannel.appendLine('Starting Kro Language Server extension');

  // Start the client. This will also launch the server
  client.start();
}

export function deactivate(): Thenable<void> | undefined {
  if (!client) {
    return undefined;
  }
  return client.stop();
}