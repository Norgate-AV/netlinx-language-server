package semantic

import (
	"fmt"
	"strings"

	"github.com/Norgate-AV/netlinx-language-server/internal/lsp"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

func SizeOfDataType(dataType SymbolDataType) uint {
	dataType = strings.ToLower(dataType)

	switch dataType {
	case SymbolDataTypeChar:
		return 1
	case SymbolDataTypeWideChar:
		return 2
	case SymbolDataTypeInteger:
		return 2
	case SymbolDataTypeSinteger:
		return 2
	case SymbolDataTypeLong:
		return 4
	case SymbolDataTypeSlong:
		return 4
	case SymbolDataTypeFloat:
		return 8
	case SymbolDataTypeDouble:
		return 8
	case SymbolDataTypeDev:
		return 6
	case SymbolDataTypeDevChan:
		return 8
	case SymbolDataTypeDevLev:
		return 8
	default:
		// How can I work out the size of a custom data type?
		return 0
	}
}

func GetNodeRange(node *tree_sitter.Node) lsp.Range {
	return lsp.Range{
		Start: lsp.Position{
			Line:      node.StartPosition().Row,
			Character: node.StartPosition().Column,
		},
		End: lsp.Position{
			Line:      node.EndPosition().Row,
			Character: node.EndPosition().Column,
		},
	}
}

func GetNodeValue(node *tree_sitter.Node, content []byte) string {
	var value string

	if node.ChildCount() >= 3 { // Left side + equals + right side
		valueNode := node.Child(2) // Third child is the value
		if valueNode != nil {
			value = valueNode.Utf8Text(content)
		}
	}

	return value
}

func IsLeftHandSide(parent *tree_sitter.Node, node *tree_sitter.Node) bool {
	// Check if the node is a left-hand side assignment
	if parent.ChildCount() > 0 {
		firstChild := parent.Child(0)
		return (firstChild.StartByte() == node.StartByte() &&
			firstChild.EndByte() == node.EndByte())
	}

	return false
}

func IsArray(node *tree_sitter.Node) bool {
	// Check if this node is inside an array_declarator
	current := node

	for current != nil {
		if current.Kind() == "array_declarator" {
			return true
		}

		current = current.Parent()
	}

	return false
}

func GetArrayDimensions(node *tree_sitter.Node) uint {
	// Find the array_declarator node
	var arrayDecl *tree_sitter.Node

	current := node
	for current != nil {
		if current.Kind() == "array_declarator" {
			arrayDecl = current
			break
		}
		current = current.Parent()
	}

	if arrayDecl == nil {
		return 0
	}

	// Count the number of bracket pairs
	dimensions := uint(0)
	for i := uint(0); i < arrayDecl.ChildCount(); i++ {
		if arrayDecl.Child(i).Kind() == "[" {
			dimensions++
		}
	}

	return dimensions
}

func PrintSymbolTable(table *SymbolTable) {
	for name, symbol := range table.Symbols {
		for i := range symbol {
			fmt.Printf("Symbol: %s = %v\n", name, symbol[i])
		}
	}
}

func GetMatchCount(query *tree_sitter.Query, tree *tree_sitter.Tree) int {
	count := 0

	cursor := tree_sitter.NewQueryCursor()
	defer cursor.Close()

	matches := cursor.Matches(query, tree.RootNode(), nil)

	for matches.Next() != nil {
		count++
	}

	return count
}

func isSectionMatch(match *tree_sitter.QueryMatch, query *tree_sitter.Query) bool {
	for _, capture := range match.Captures {
		name := query.CaptureNames()[capture.Index]

		if strings.HasPrefix(name, "section.") {
			return true
		}
	}

	return false
}

func getSectionType(match *tree_sitter.QueryMatch, query *tree_sitter.Query) string {
	for _, capture := range match.Captures {
		name := query.CaptureNames()[capture.Index]

		if strings.HasPrefix(name, "section.") {
			// May want to trim the prefix here?
			return name
		}
	}

	return ""
}
