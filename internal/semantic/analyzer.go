package semantic

import (
	"github.com/Norgate-AV/netlinx-language-server/internal/logger"
	"github.com/Norgate-AV/netlinx-language-server/internal/parser"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type Analyzer struct {
	parser   *parser.TreeSitter
	logger   logger.Logger
	document *Document
	content  []byte
}

func NewAnalyzer(parser *parser.TreeSitter, logger logger.Logger) *Analyzer {
	return &Analyzer{
		parser: parser,
		logger: logger,
	}
}

// Analyze performs semantic analysis on a NetLinx document
func (a *Analyzer) Analyze(uri string, content string, tree *tree_sitter.Tree) (*Document, error) {
	if tree == nil || tree.RootNode() == nil {
		return &Document{URI: uri, Symbols: make(map[string]*Symbol)}, nil
	}

	// Initialize document and state
	a.document = &Document{
		URI:           uri,
		Symbols:       make(map[string]*Symbol),
		SymbolsByKind: make(map[SymbolType][]*Symbol),
	}
	a.content = []byte(content)

	// Create global scope
	// a.document.GlobalScope = &Scope{
	// 	Symbols:   make(map[string]*Symbol),
	// 	Node:      tree.RootNode(),
	// 	Range:     nodeToRange(tree.RootNode()),
	// 	ScopeType: "global",
	// }

	a.document.Scopes = append(a.document.Scopes, a.document.GlobalScope)

	return a.document, nil
}
