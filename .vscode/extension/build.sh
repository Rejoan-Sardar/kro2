#!/bin/bash

# Fail on any error
set -e

# Check Go version
echo "Checking Go version..."
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.19"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo "Error: Go version $REQUIRED_VERSION or higher is required (found $GO_VERSION)"
    exit 1
fi

# Build the Kro Language Server
echo "Building Kro Language Server..."
cd server
go build -o kro-lsp
if [ ! -f "kro-lsp" ]; then
    echo "Error: Failed to build kro-lsp server"
    exit 1
fi
echo "Server built successfully"
cd ..

# Set up the VSCode extension
echo "Setting up VSCode extension..."
cd client
if ! command -v npm &> /dev/null; then
    echo "Error: npm is required but not installed"
    exit 1
fi

echo "Installing dependencies..."
npm install
if [ $? -ne 0 ]; then
    echo "Error: Failed to install npm dependencies"
    exit 1
fi

echo "Compiling TypeScript code..."
npm run compile
if [ $? -ne 0 ]; then
    echo "Error: Failed to compile TypeScript code"
    exit 1
fi
cd ..

echo "Build complete!"
echo "VSCode extension is ready in kro/.vscode/extension"
echo "To test the extension, copy the extension folder to your VSCode extensions directory"
echo "or use 'code --extensionDevelopmentPath=/path/to/kro/.vscode/extension' to start VSCode with the extension"
echo ""
echo "To test the LSP server directly, run:"
echo "cd server && ./test-client.sh"