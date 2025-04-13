package parser

import (
	"fmt"

	tree_sitter_netlinx "github.com/norgate-av/tree-sitter-netlinx/bindings/go"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type TreeSitter struct {
	Parser   *tree_sitter.Parser
	Language *tree_sitter.Language
}

var NetLinxLanguage = tree_sitter.NewLanguage(tree_sitter_netlinx.Language())

func NewTreeSitter() (*TreeSitter, error) {
	parser := tree_sitter.NewParser()
	if err := parser.SetLanguage(NetLinxLanguage); err != nil {
		return nil, fmt.Errorf("failed to set language: %w", err)
	}

	return &TreeSitter{
		Parser:   parser,
		Language: NetLinxLanguage,
	}, nil
}

func (ts *TreeSitter) Close() {
	if ts.Parser == nil {
		return
	}

	ts.Parser.Close()
}

func CreateQuery(query string) (*tree_sitter.Query, error) {
	q, err := tree_sitter.NewQuery(NetLinxLanguage, query)
	if err != nil {
		return nil, fmt.Errorf("failed to create query: %w", err)
	}

	return q, nil
}

func CreateQueryCursor() *tree_sitter.QueryCursor {
	return tree_sitter.NewQueryCursor()
}

// func (ts *TreeSitter) GetLanguage() *tree_sitter.Language {
// 	return ts.language
// }

// func findFirstErrorNode(node *tree_sitter.Node) *tree_sitter.Node {
// 	if !node.HasError() {
// 		return nil
// 	}

// 	count := node.ChildCount()

// 	for i := uint(0); i < count; i++ {
// 		child := node.Child(i)

// 		if child.HasError() {
// 			return findFirstErrorNode(child)
// 		}
// 	}

// 	return node
// }
