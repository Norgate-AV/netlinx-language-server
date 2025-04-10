package analysis

import (
	"sync"

	"github.com/Norgate-AV/netlinx-language-server/internal/lsp"
	"github.com/Norgate-AV/netlinx-language-server/internal/parser"
)

type State struct {
	Documents map[string]lsp.DocumentUri
	Parser    *parser.Parser
	mutex     sync.RWMutex
}

func NewState() *State {
	return &State{
		Documents: make(map[string]lsp.DocumentUri),
		Parser:    parser.NewParser(),
	}
}

func (s *State) AddDocument(uri string, content string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.Documents[uri] = content
}

func (s *State) GetDocument(uri string) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	content, ok := s.Documents[uri]

	return content, ok
}

func (s *State) UpdateDocument(uri string, content string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.Documents[uri] = content
}

func (s *State) CloseDocument(uri string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.Documents, uri)
}

func (s *State) GetSyntaxTree(uri string) (*parser.Tree, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	content, ok := s.Documents[uri]
	if !ok {
		return nil, false
	}

	tree := s.Parser.Parse([]byte(content))

	return tree, true
}
