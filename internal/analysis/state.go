package analysis

import (
	"fmt"
	"sync"

	"github.com/Norgate-AV/netlinx-language-server/internal/analysis/semantic"
	"github.com/Norgate-AV/netlinx-language-server/internal/logger"
	"github.com/Norgate-AV/netlinx-language-server/internal/lsp"
	"github.com/Norgate-AV/netlinx-language-server/parser"
	"github.com/sirupsen/logrus"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type Document struct {
	Content       string
	Tree          *tree_sitter.Tree
	SemanticModel *semantic.Document
}

type State struct {
	Documents  map[lsp.DocumentUri]*Document
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
		Documents:  make(map[lsp.DocumentUri]*Document),
		TreeSitter: options.TreeSitter,
		Logger:     options.Logger,
	}
}

func (s *State) AnalyzeDocument(uri string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	doc, exists := s.Documents[uri]
	if !exists || doc.Tree == nil {
		return fmt.Errorf("document not found or has no syntax tree")
	}

	analyzer := semantic.NewAnalyzer(s.TreeSitter, s.Logger)
	semanticDoc, err := analyzer.Analyze(uri, doc.Content, doc.Tree)
	if err != nil {
		return err
	}

	doc.SemanticModel = semanticDoc
	return nil
}

func (s *State) AddDocument(uri string, content string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.AddDocumentLocked(uri, content)
}

func (s *State) AddDocumentLocked(uri string, content string) {
	tree := s.TreeSitter.Parser.Parse([]byte(content), nil)

	doc := &Document{
		Content: content,
		Tree:    tree,
	}

	s.Documents[uri] = doc

	// Analyze document to build semantic model
	analyzer := semantic.NewAnalyzer(s.TreeSitter, s.Logger)
	semanticDoc, err := analyzer.Analyze(uri, content, tree)
	if err != nil {
		s.Logger.Error("Failed to analyze document", logrus.Fields{
			"uri":   uri,
			"error": err,
		})
	} else {
		doc.SemanticModel = semanticDoc
	}
}

func (s *State) GetDocument(uri string) (*Document, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	document, exists := s.Documents[uri]

	if !exists {
		return nil, false
	}

	return document, true
}

func (s *State) UpdateDocument(uri string, content string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if document, exists := s.Documents[uri]; exists {
		document.Content = content
		document.Tree = s.TreeSitter.Parser.Parse([]byte(document.Content), document.Tree)

		return
	}

	// If the document doesn't exist, create a new one
	s.AddDocumentLocked(uri, content)
}

func (s *State) CloseDocument(uri string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if document, exists := s.Documents[uri]; exists && document.Tree != nil {
		document.Tree.Close()
		document.Tree = nil
	}

	delete(s.Documents, uri)
}

func (s *State) GetSyntaxTree(uri string) (*tree_sitter.Tree, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	document, exists := s.Documents[uri]

	if !exists || document.Tree == nil {
		return nil, false
	}

	return document.Tree, true
}
