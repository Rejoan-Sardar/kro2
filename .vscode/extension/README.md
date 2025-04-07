# Kro Language Server Protocol (LSP) Extension

A VSCode extension providing advanced language support for Kro ResourceGraphDefinitions (.kro files).

## Features

- Syntax highlighting for Kro ResourceGraphDefinition files
- Schema validation
- Code completion for ResourceGraphDefinition fields
- Hover documentation
- CEL expression validation
- Template syntax support with {{ }} expressions

## Architecture

This extension consists of two main components:

1. **LSP Server**: A Go language implementation that provides:
   - Validation of ResourceGraphDefinition schemas
   - Template expression validation
   - Intelligent code completions
   - Hover documentation

2. **VSCode Client**: A TypeScript implementation that:
   - Communicates with the LSP server
   - Provides syntax highlighting using TextMate grammar
   - Handles VSCode-specific integration

## Installation

### Prerequisites

- Go 1.19+
- Node.js 14+
- npm

### Building the Extension

1. Clone the repository:
   ```
   git clone https://github.com/your-org/kro.git
   cd kro/.vscode/extension
   ```

2. Run the build script:
   ```
   ./build.sh
   ```

This will:
- Build the Go LSP server
- Install Node.js dependencies
- Compile the TypeScript client

### Installing in VSCode

Method 1: Symlink the extension directory
```
ln -s /path/to/kro/.vscode/extension ~/.vscode/extensions/kro-extension
```

Method 2: Launch VSCode with the extension
```
code --extensionDevelopmentPath=/path/to/kro/.vscode/extension
```

## Testing

### Testing the LSP Server Directly

A test script is provided to verify the LSP server is working correctly:

```
cd server
./test-client.sh
```

This script:
1. Builds and starts the LSP server
2. Verifies the health endpoint is working
3. Sends an LSP initialize request
4. Checks that the server responds with expected capabilities

### Testing with VSCode

1. Open a Kro ResourceGraphDefinition file (`.kro` or YAML file with `kind: ResourceGraphDefinition`)
2. Verify syntax highlighting is working
3. Test auto-completion by typing fields like `apiVersion`, `kind`, `metadata`, etc.
4. Test hover documentation by hovering over fields
5. Test validation by introducing errors in the file

## Development

### LSP Server

The Go-based LSP server is in the `server` directory:

- `main.go`: Entry point, sets up HTTP and TCP servers
- `server.go`: Core server implementation
- `handler.go`: Request/notification handlers
- `document.go`: Document representation and methods
- `document_store.go`: Document storage
- `validation.go`: Schema and expression validation
- `completion.go`: Code completion implementation
- `diagnostics.go`: Diagnostic processing

### VSCode Client

The TypeScript client is in the `client` directory:

- `src/extension.ts`: Main extension code
- `syntaxes/kro.tmLanguage.json`: TextMate grammar for syntax highlighting
- `package.json`: Extension metadata and configuration

## Extension Settings

- `kro.server.path`: Path to the LSP server executable (default: bundled server)
- `kro.trace.server`: Trace communication between VSCode and the language server

## Contributing

See [CONTRIBUTING.md](../../CONTRIBUTING.md) for details on how to contribute to this extension.

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](../../LICENSE) file for details.
