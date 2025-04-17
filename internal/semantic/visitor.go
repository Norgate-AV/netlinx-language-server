package semantic

import (
	"github.com/Norgate-AV/netlinx-language-server/internal/logger"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

// NodeVisitor processes a node and returns whether to continue visiting its children
type NodeVisitor interface {
	Visit(node *tree_sitter.Node, content []byte) bool
}

// TreeWalker walks the entire tree using a visitor
type TreeWalker struct {
	visitor NodeVisitor
	content []byte
	logger  logger.Logger
}

// NewTreeWalker creates a new tree walker
func NewTreeWalker(visitor NodeVisitor, content []byte, logger logger.Logger) *TreeWalker {
	return &TreeWalker{
		visitor: visitor,
		content: content,
		logger:  logger,
	}
}

// Walk traverses the entire tree applying the visitor
func (w *TreeWalker) Walk(node *tree_sitter.Node) {
	if node == nil {
		return
	}

	// Visit this node first - if it returns false, don't traverse children
	if !w.visitor.Visit(node, w.content) {
		return
	}

	// Recursively visit all children
	for i := uint(0); i < node.ChildCount(); i++ {
		w.Walk(node.Child(i))
	}
}
