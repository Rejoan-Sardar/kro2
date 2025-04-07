#!/bin/bash

# Simple script to test the Kro LSP server by sending some basic LSP messages

echo "Testing Kro LSP Server..."

if [ ! -f "./kro-lsp" ]; then
  echo "Error: kro-lsp server binary not found. Please build it first."
  exit 1
fi

# Start the server in the background
echo "Starting the LSP server in the background..."
./kro-lsp --port 5000 &
SERVER_PID=$!

# Give the server a moment to start
sleep 1

# Check if health endpoint is working
echo "Checking health endpoint..."
HEALTH_RESULT=$(curl -s http://localhost:5000/health)
if [ "$HEALTH_RESULT" != "Kro LSP Server is running" ]; then
  echo "Error: Health check failed"
  kill $SERVER_PID
  exit 1
fi
echo "Health check successful!"

# Simple test using netcat to send LSP requests
echo "Sending LSP requests to the server..."

# Initialize request
cat > init_request.json << EOF2
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "processId": 123,
    "rootUri": "file:///tmp",
    "capabilities": {}
  }
}
EOF2

# Add content length header to the request
CONTENT=$(cat init_request.json)
LENGTH=$(echo -n "$CONTENT" | wc -c)
REQUEST="Content-Length: $LENGTH\r\n\r\n$CONTENT"

# Send the request and receive the response
echo "Sending initialize request..."
echo -e "$REQUEST" | nc localhost 5001 > response.txt &
NETCAT_PID=$!

# Wait a moment for the response
sleep 2

# Check if we got a response
if [ -s response.txt ]; then
  echo "Received response from server!"
  
  # Extract the JSON part (skip headers)
  RESPONSE_JSON=$(cat response.txt | sed '1,/^\r$/d')
  
  # Check if the response contains server capabilities
  if echo "$RESPONSE_JSON" | grep -q "capabilities"; then
    echo "LSP server properly responded with capabilities!"
  else
    echo "LSP response doesn't contain capabilities section."
  fi
else
  echo "No response received from the server."
fi

# Clean up
echo "Cleaning up..."
kill $SERVER_PID
kill $NETCAT_PID 2> /dev/null
rm init_request.json response.txt 2> /dev/null

echo "Test complete."
