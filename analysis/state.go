package analysis

import "sync"

// State holds the state of the language server
type State struct {
	Documents map[string]string
	mu        sync.RWMutex // Protect concurrent access to Documents
}

// NewState creates a new State instance
func NewState() State {
	return State{
		Documents: make(map[string]string),
	}
}

// OpenDocument adds a document to the state
func (s *State) OpenDocument(uri string, content string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Documents[uri] = content
}

// GetDocument returns the content of a document
func (s *State) GetDocument(uri string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	content, ok := s.Documents[uri]
	return content, ok
}

// UpdateDocument updates the content of a document
func (s *State) UpdateDocument(uri string, content string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Documents[uri] = content
}

// CloseDocument removes a document from the state
func (s *State) CloseDocument(uri string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Documents, uri)
}
