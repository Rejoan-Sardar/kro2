package main

import (
        "context"
        "encoding/json"
        "errors"
        "fmt"
        "log"
        "net"
        "net/http"
        "os"
        "os/signal"
        "syscall"

        "go.lsp.dev/jsonrpc2"
)

// LSPServer represents the LSP server implementation
type LSPServer struct {
        handler       *Handler
        documentStore *DocumentStore
        logger        *log.Logger
        ctx           context.Context
        cancelFunc    context.CancelFunc
        
        // Server addresses
        lspAddr      string
        healthAddr   string
}

// NewLSPServer creates a new LSP server
func NewLSPServer(lspAddr, healthAddr string) *LSPServer {
        ctx, cancel := context.WithCancel(context.Background())
        documentStore := NewDocumentStore()
        handler := NewHandler(documentStore)
        
        return &LSPServer{
                handler:       handler,
                documentStore: documentStore,
                logger:        log.New(os.Stderr, "[server] ", log.LstdFlags),
                ctx:           ctx,
                cancelFunc:    cancel,
                lspAddr:       lspAddr,
                healthAddr:    healthAddr,
        }
}

// Start starts the LSP server
func (s *LSPServer) Start() error {
        // Start the health check server
        go s.startHealthServer()
        
        // Start the LSP server
        return s.startLSPServer()
}

// startHealthServer starts the health check HTTP server
func (s *LSPServer) startHealthServer() {
        if s.healthAddr == "" {
                s.logger.Println("Health server disabled (no address provided)")
                return
        }
        
        mux := http.NewServeMux()
        
        // Add health check endpoint
        mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
                w.Header().Set("Content-Type", "application/json")
                resp := map[string]interface{}{
                        "status":   "ok",
                        "version":  ServerVersion,
                        "serverId": ServerName,
                        "docs":     s.documentStore.GetDocumentCount(),
                }
                json.NewEncoder(w).Encode(resp)
        })
        
        // Start the HTTP server
        s.logger.Printf("Starting health server on %s", s.healthAddr)
        if err := http.ListenAndServe(s.healthAddr, mux); err != nil && !errors.Is(err, http.ErrServerClosed) {
                s.logger.Printf("Health server error: %v", err)
        }
}

// startLSPServer starts the LSP server
func (s *LSPServer) startLSPServer() error {
        // Start TCP listener
        s.logger.Printf("Starting LSP server on %s", s.lspAddr)
        listener, err := net.Listen("tcp", s.lspAddr)
        if err != nil {
                return fmt.Errorf("failed to listen on %s: %v", s.lspAddr, err)
        }
        defer listener.Close()
        
        // Setup signal handling for graceful shutdown
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
        
        go func() {
                sig := <-sigChan
                s.logger.Printf("Received signal %v, shutting down...", sig)
                s.cancelFunc()
                listener.Close()
        }()
        
        // Accept connections
        for {
                conn, err := listener.Accept()
                if err != nil {
                        select {
                        case <-s.ctx.Done():
                                return nil // Server is shutting down
                        default:
                                s.logger.Printf("Failed to accept connection: %v", err)
                                continue
                        }
                }
                
                s.logger.Printf("New connection from %s", conn.RemoteAddr())
                
                // Handle the connection in a goroutine
                go s.handleConnection(conn)
        }
}

// handleConnection handles a single client connection
func (s *LSPServer) handleConnection(netConn net.Conn) {
        defer netConn.Close()
        
        // Create JSON-RPC stream
        stream := jsonrpc2.NewStream(netConn)
        
        // Create JSON-RPC connection
        conn := jsonrpc2.NewConn(stream)
        s.handler.SetConnection(conn)
        
        // Start processing in a goroutine
        conn.Go(s.ctx, s.handler.GetHandler())
        
        // Wait for the connection to be closed
        <-conn.Done()
        
        if err := conn.Err(); err != nil {
                s.logger.Printf("Connection closed with error: %v", err)
        } else {
                s.logger.Printf("Connection closed")
        }
}

// Stop stops the LSP server
func (s *LSPServer) Stop() {
        s.cancelFunc()
}
