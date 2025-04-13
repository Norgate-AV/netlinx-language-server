package analysis

import (
	"sync"

	"github.com/Norgate-AV/netlinx-language-server/internal/logger"
	"github.com/Norgate-AV/netlinx-language-server/internal/lsp"
	"github.com/Norgate-AV/netlinx-language-server/parser"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type State struct {
	Documents  map[string]lsp.DocumentUri
	TreeSitter *parser.TreeSitter
	mutex      sync.RWMutex
	Logger     logger.Logger
}

type NewStateOptions struct {
	Logger     logger.Logger
	TreeSitter *parser.TreeSitter
}

func NewState(options *NewStateOptions) *State {
	return &State{
		Documents:  make(map[string]lsp.DocumentUri),
		TreeSitter: options.TreeSitter,
		Logger:     options.Logger,
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

func (s *State) GetSyntaxTree(uri string) (*tree_sitter.Tree, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	content, ok := s.Documents[uri]
	if !ok {
		return nil, false
	}

	tree := s.TreeSitter.Parser.Parse([]byte(content), nil)
	if tree == nil {
		return nil, false
	}

	return tree, true
}
