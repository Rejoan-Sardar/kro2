package main

import (
        "context"
        "encoding/json"
        "fmt"
        "log"
        "os"

        "go.lsp.dev/jsonrpc2"
        "go.lsp.dev/protocol"
        "go.lsp.dev/uri"
)

// Handler handles LSP protocol messages
type Handler struct {
        conn          jsonrpc2.Conn
        documentStore *DocumentStore
        capabilities  protocol.ClientCapabilities
        initialized   bool
        logger        *log.Logger
}

// NewHandler creates a new LSP message handler
func NewHandler(documentStore *DocumentStore) *Handler {
        return &Handler{
                documentStore: documentStore,
                initialized:   false,
                logger:        log.New(os.Stderr, "[handler] ", log.LstdFlags),
        }
}

// SetConnection sets the connection for sending notifications and responses
func (h *Handler) SetConnection(conn jsonrpc2.Conn) {
        h.conn = conn
}

// GetHandler returns a jsonrpc2.Handler function that can be used with conn.Go()
func (h *Handler) GetHandler() jsonrpc2.Handler {
        return func(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
                return h.Handle(ctx, reply, req)
        }
}

// Handle processes an LSP request/notification
func (h *Handler) Handle(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
        h.logger.Printf("Received request: %s", req.Method())

        switch req.Method() {
        case protocol.MethodInitialize:
                return h.handleInitialize(ctx, reply, req)
        case protocol.MethodInitialized:
                return h.handleInitialized(ctx, reply, req)
        case protocol.MethodTextDocumentDidOpen:
                return h.handleTextDocumentDidOpen(ctx, reply, req)
        case protocol.MethodTextDocumentDidChange:
                return h.handleTextDocumentDidChange(ctx, reply, req)
        case protocol.MethodTextDocumentDidClose:
                return h.handleTextDocumentDidClose(ctx, reply, req)
        case protocol.MethodTextDocumentCompletion:
                return h.handleTextDocumentCompletion(ctx, reply, req)
        case protocol.MethodShutdown:
                return h.handleShutdown(ctx, reply, req)
        case protocol.MethodExit:
                return h.handleExit(ctx, reply, req)
        default:
                h.logger.Printf("Unknown method: %s", req.Method())
                // Check if this is a notification (requests have an ID, notifications don't)
                if req.Params() == nil {
                        return nil // Treat as notification
                }
                return reply(ctx, nil, fmt.Errorf("unsupported method: %s", req.Method()))
        }
}

// handleInitialize processes the initialize request
func (h *Handler) handleInitialize(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
        h.logger.Printf("Handling initialize request")

        var params protocol.InitializeParams
        if err := json.Unmarshal(req.Params(), &params); err != nil {
                return reply(ctx, nil, fmt.Errorf("invalid initialize params: %v", err))
        }

        h.capabilities = params.Capabilities

        // Save client capabilities for future reference
        h.capabilities = params.Capabilities

        // Create server capabilities
        result := protocol.InitializeResult{
                Capabilities: protocol.ServerCapabilities{
                        TextDocumentSync: &protocol.TextDocumentSyncOptions{
                                OpenClose: true,
                                Change:    protocol.TextDocumentSyncKindFull, // Full sync for now
                        },
                        CompletionProvider: &protocol.CompletionOptions{
                                TriggerCharacters: []string{".", ":", " "},
                        },
                        DocumentFormattingProvider: false,
                        HoverProvider:              true,
                        DefinitionProvider:         false,
                        ReferencesProvider:         false,
                        DocumentSymbolProvider:     false,
                        RenameProvider:             false,
                        WorkspaceSymbolProvider:    false,
                },
                ServerInfo: &protocol.ServerInfo{
                        Name:    ServerName,
                        Version: ServerVersion,
                },
        }

        h.initialized = true
        return reply(ctx, result, nil)
}

// handleInitialized processes the initialized notification
func (h *Handler) handleInitialized(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
        h.logger.Printf("Handling initialized notification")
        
        // Show welcome message
        h.showMessage(ctx, protocol.MessageTypeInfo, "Kro LSP Server initialized")
        
        return nil
}

// handleTextDocumentDidOpen processes the textDocument/didOpen notification
func (h *Handler) handleTextDocumentDidOpen(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
        h.logger.Printf("Handling textDocument/didOpen notification")

        var params protocol.DidOpenTextDocumentParams
        if err := json.Unmarshal(req.Params(), &params); err != nil {
                return fmt.Errorf("invalid didOpen params: %v", err)
        }

        // Create and store the document
        document := NewDocument(
                uri.URI(params.TextDocument.URI),
                params.TextDocument.Text,
                string(params.TextDocument.LanguageID),
        )

        h.documentStore.UpsertDocument(string(params.TextDocument.URI), document)

        // Validate the document
        diagnostics := document.Validate()
        h.publishDiagnostics(ctx, uri.URI(params.TextDocument.URI), diagnostics)

        return nil
}

// handleTextDocumentDidChange processes the textDocument/didChange notification
func (h *Handler) handleTextDocumentDidChange(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
        h.logger.Printf("Handling textDocument/didChange notification")

        var params protocol.DidChangeTextDocumentParams
        if err := json.Unmarshal(req.Params(), &params); err != nil {
                return fmt.Errorf("invalid didChange params: %v", err)
        }

        docURI := string(params.TextDocument.URI)
        document := h.documentStore.GetDocument(docURI)
        if document == nil {
                // Document not found, create it
                h.logger.Printf("Document not found, creating new: %s", docURI)
                document = NewDocument(
                        uri.URI(params.TextDocument.URI),
                        "",
                        "",  // Language ID not provided in didChange
                )
                h.documentStore.UpsertDocument(docURI, document)
        }

        // Update content (assuming full content update, not incremental)
        if len(params.ContentChanges) > 0 {
                document.Content = params.ContentChanges[0].Text
        }

        // Validate the document
        diagnostics := document.Validate()
        h.publishDiagnostics(ctx, uri.URI(params.TextDocument.URI), diagnostics)

        return nil
}

// handleTextDocumentDidClose processes the textDocument/didClose notification
func (h *Handler) handleTextDocumentDidClose(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
        h.logger.Printf("Handling textDocument/didClose notification")

        var params protocol.DidCloseTextDocumentParams
        if err := json.Unmarshal(req.Params(), &params); err != nil {
                return fmt.Errorf("invalid didClose params: %v", err)
        }

        // Remove document from store
        h.documentStore.DeleteDocument(string(params.TextDocument.URI))

        // Clear diagnostics by publishing an empty array
        h.publishDiagnostics(ctx, uri.URI(params.TextDocument.URI), []protocol.Diagnostic{})

        return nil
}

// handleTextDocumentCompletion processes the textDocument/completion request
func (h *Handler) handleTextDocumentCompletion(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
        h.logger.Printf("Handling textDocument/completion request")

        var params protocol.CompletionParams
        if err := json.Unmarshal(req.Params(), &params); err != nil {
                return reply(ctx, nil, fmt.Errorf("invalid completion params: %v", err))
        }

        // Get the document
        docURI := string(params.TextDocument.URI)
        document := h.documentStore.GetDocument(docURI)
        if document == nil {
                return reply(ctx, nil, fmt.Errorf("document not found: %s", docURI))
        }

        // Get completion items for the position
        completionItems := document.GetCompletionItems(params.Position)

        return reply(ctx, completionItems, nil)
}

// handleShutdown processes the shutdown request
func (h *Handler) handleShutdown(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
        h.logger.Printf("Handling shutdown request")
        h.initialized = false
        return reply(ctx, nil, nil)
}

// handleExit processes the exit notification
func (h *Handler) handleExit(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
        h.logger.Printf("Handling exit notification")
        
        // If we haven't received a shutdown request, exit with error code 1
        if h.initialized {
                os.Exit(1)
        }
        
        os.Exit(0)
        return nil
}

// publishDiagnostics sends diagnostics to the client
func (h *Handler) publishDiagnostics(ctx context.Context, docURI uri.URI, diagnostics []protocol.Diagnostic) {
        // Ensure diagnostics is never null
        if diagnostics == nil {
                diagnostics = make([]protocol.Diagnostic, 0)
        }

        // Create params
        params := &protocol.PublishDiagnosticsParams{
                URI:         protocol.DocumentURI(docURI),
                Diagnostics: diagnostics,
        }

        // Notify the client
        if err := h.conn.Notify(ctx, protocol.MethodTextDocumentPublishDiagnostics, params); err != nil {
                h.logger.Printf("Failed to publish diagnostics: %v", err)
        }
}

// showMessage sends a notification to the client to display a message
func (h *Handler) showMessage(ctx context.Context, msgType protocol.MessageType, message string) {
        params := &protocol.ShowMessageParams{
                Type:    msgType,
                Message: message,
        }
        if err := h.conn.Notify(ctx, protocol.MethodWindowShowMessage, params); err != nil {
                h.logger.Printf("Failed to show message: %v", err)
        }
}

// logMessage sends a notification to the client to log a message
func (h *Handler) logMessage(ctx context.Context, msgType protocol.MessageType, message string) {
        params := &protocol.LogMessageParams{
                Type:    msgType,
                Message: message,
        }
        if err := h.conn.Notify(ctx, protocol.MethodWindowLogMessage, params); err != nil {
                h.logger.Printf("Failed to log message: %v", err)
        }
}
