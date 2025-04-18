package parser

import (
	"fmt"
	"strings"

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

type PrettyPrintOptions struct {
	// Ranges indicates whether to include ranges in the output
	ShowRanges bool
}

// PrettyPrintSexp formats a tree as S-expression
func PrettyPrint(tree *tree_sitter.Tree, options PrettyPrintOptions) string {
	var result strings.Builder

	cursor := tree.Walk()
	defer cursor.Close()

	indentLevel := 0
	needsNewline := false
	didVisitChildren := false

	for {
		node := cursor.Node()
		isNamed := node.IsNamed()

		if didVisitChildren {
			if isNamed {
				result.WriteString(")")
				needsNewline = true
			}

			if cursor.GotoNextSibling() {
				didVisitChildren = false
			} else if cursor.GotoParent() {
				didVisitChildren = true
				indentLevel -= 1
			} else {
				break
			}
		} else {
			if isNamed {
				if needsNewline {
					result.WriteString("\n")
				}

				for range indentLevel {
					result.WriteString("  ")
				}

				fieldName := cursor.FieldName()
				if fieldName != "" {
					result.WriteString(fieldName)
					result.WriteString(": ")
				}

				result.WriteString("(")
				result.WriteString(node.Kind())

				if options.ShowRanges {
					start := node.StartPosition()
					end := node.EndPosition()
					result.WriteString(fmt.Sprintf(" [%d, %d] - [%d, %d]", start.Row, start.Column, end.Row, end.Column))
				}

				needsNewline = true
			}

			if cursor.GotoFirstChild() {
				didVisitChildren = false
				indentLevel += 1
			} else {
				didVisitChildren = true
			}
		}
	}

	result.WriteString("\n")
	return result.String()
}
