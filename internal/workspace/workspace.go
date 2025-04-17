package workspace

import (
	"fmt"
	"sync"

	"github.com/Norgate-AV/netlinx-language-server/internal/logger"
	"github.com/Norgate-AV/netlinx-language-server/internal/lsp"
	"github.com/Norgate-AV/netlinx-language-server/internal/semantic"
	"github.com/Norgate-AV/netlinx-language-server/parser"
	"github.com/sirupsen/logrus"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type Document struct {
	Content       string
	Tree          *tree_sitter.Tree
	SemanticModel *semantic.Document
}

type Workspace struct {
	Documents  map[lsp.DocumentUri]*Document
	TreeSitter *parser.TreeSitter
	mutex      sync.RWMutex
	Logger     logger.Logger
}

type Options struct {
	Logger     logger.Logger
	TreeSitter *parser.TreeSitter
}

func NewWorkspace(options *Options) *Workspace {
	return &Workspace{
		Documents:  make(map[lsp.DocumentUri]*Document),
		TreeSitter: options.TreeSitter,
		Logger:     options.Logger,
	}
}

func (w *Workspace) AnalyzeDocument(uri string) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	doc, exists := w.Documents[uri]
	if !exists || doc.Tree == nil {
		return fmt.Errorf("document not found or has no syntax tree")
	}

	analyzer := semantic.NewAnalyzer(w.TreeSitter, w.Logger)
	semanticDoc, err := analyzer.Analyze(uri, doc.Content, doc.Tree)
	if err != nil {
		return err
	}

	doc.SemanticModel = semanticDoc
	return nil
}

func (w *Workspace) AddDocument(uri string, content string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.AddDocumentLocked(uri, content)
}

func (w *Workspace) AddDocumentLocked(uri string, content string) {
	tree := w.TreeSitter.Parser.Parse([]byte(content), nil)

	doc := &Document{
		Content: content,
		Tree:    tree,
	}

	w.Documents[uri] = doc

	// Analyze document to build semantic model
	analyzer := semantic.NewAnalyzer(w.TreeSitter, w.Logger)
	semanticDoc, err := analyzer.Analyze(uri, content, tree)
	if err != nil {
		w.Logger.Error("Failed to analyze document", logrus.Fields{
			"uri":   uri,
			"error": err,
		})
	} else {
		doc.SemanticModel = semanticDoc
	}
}

func (w *Workspace) GetDocument(uri string) (*Document, bool) {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	document, exists := w.Documents[uri]

	if !exists {
		return nil, false
	}

	return document, true
}

func (w *Workspace) UpdateDocument(uri string, content string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if document, exists := w.Documents[uri]; exists {
		document.Content = content
		document.Tree = w.TreeSitter.Parser.Parse([]byte(document.Content), document.Tree)

		return
	}

	// If the document doesn't exist, create a new one
	w.AddDocumentLocked(uri, content)
}

func (w *Workspace) CloseDocument(uri string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if document, exists := w.Documents[uri]; exists && document.Tree != nil {
		document.Tree.Close()
		document.Tree = nil
	}

	delete(w.Documents, uri)
}

func (w *Workspace) GetSyntaxTree(uri string) (*tree_sitter.Tree, bool) {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	document, exists := w.Documents[uri]

	if !exists || document.Tree == nil {
		return nil, false
	}

	return document.Tree, true
}
