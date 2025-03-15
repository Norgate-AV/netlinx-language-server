package lsp

import (
	"sync"
)

// Document represents a text document open in the editor
type Document struct {
	URI        string
	LanguageID string
	Version    int
	Content    string
}

// DocumentManager keeps track of all open documents
type DocumentManager struct {
	mu    sync.RWMutex
	docs  map[string]*Document
}

// NewDocumentManager creates a new document manager
func NewDocumentManager() *DocumentManager {
	return &DocumentManager{
		docs: make(map[string]*Document),
	}
}

// Add adds a document to the manager
func (dm *DocumentManager) Add(doc *Document) {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	dm.docs[doc.URI] = doc
}

// Remove removes a document from the manager
func (dm *DocumentManager) Remove(uri string) {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	delete(dm.docs, uri)
}

// Get retrieves a document by URI
func (dm *DocumentManager) Get(uri string) (*Document, bool) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	doc, ok := dm.docs[uri]
	return doc, ok
}

// Update updates an existing document's content
func (dm *DocumentManager) Update(uri string, version int, content string) bool {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	
	if doc, ok := dm.docs[uri]; ok {
		if doc.Version < version {
			doc.Version = version
			doc.Content = content
		}
		return true
	}
	return false
}
