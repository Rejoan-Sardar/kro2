package main

import (
        "sync"
        
        "go.lsp.dev/protocol"
        "go.lsp.dev/uri"
)

// DocumentStore manages the collection of documents
type DocumentStore struct {
        documents map[string]*Document
        mu        sync.RWMutex
}

// NewDocumentStore creates a new document store
func NewDocumentStore() *DocumentStore {
        return &DocumentStore{
                documents: make(map[string]*Document),
        }
}

// GetDocument retrieves a document by URI
func (s *DocumentStore) GetDocument(uri string) *Document {
        s.mu.RLock()
        defer s.mu.RUnlock()
        return s.documents[uri]
}

// UpsertDocument adds or updates a document in the store
func (s *DocumentStore) UpsertDocument(uri string, doc *Document) {
        s.mu.Lock()
        defer s.mu.Unlock()
        s.documents[uri] = doc
}

// DeleteDocument removes a document from the store
func (s *DocumentStore) DeleteDocument(uri string) {
        s.mu.Lock()
        defer s.mu.Unlock()
        delete(s.documents, uri)
}

// GetDocumentCount returns the number of documents in the store
func (s *DocumentStore) GetDocumentCount() int {
        s.mu.RLock()
        defer s.mu.RUnlock()
        return len(s.documents)
}

// GetAllDocuments returns a copy of all documents as a map
func (s *DocumentStore) GetAllDocuments() map[string]*Document {
        s.mu.RLock()
        defer s.mu.RUnlock()
        
        // Create a copy of the map
        result := make(map[string]*Document, len(s.documents))
        for k, v := range s.documents {
                result[k] = v
        }
        
        return result
}

// GetAllDocumentsAsList returns all documents in the store as a slice
func (s *DocumentStore) GetAllDocumentsAsList() []*Document {
        s.mu.RLock()
        defer s.mu.RUnlock()
        
        docs := make([]*Document, 0, len(s.documents))
        for _, doc := range s.documents {
                docs = append(docs, doc)
        }
        
        return docs
}

// HasDocument checks if a document exists in the store
func (s *DocumentStore) HasDocument(uri string) bool {
        s.mu.RLock()
        defer s.mu.RUnlock()
        _, exists := s.documents[uri]
        return exists
}

// ValidateAllDocuments validates all documents and returns diagnostics
func (s *DocumentStore) ValidateAllDocuments() map[uri.URI][]protocol.Diagnostic {
        s.mu.RLock()
        docs := make(map[string]*Document, len(s.documents))
        for k, v := range s.documents {
                docs[k] = v
        }
        s.mu.RUnlock()
        
        result := make(map[uri.URI][]protocol.Diagnostic)
        for _, doc := range docs {
                diagnostics := doc.Validate()
                result[doc.URI] = diagnostics
        }
        
        return result
}