package semantic

import (
	"fmt"
	"strings"

	"github.com/Norgate-AV/netlinx-language-server/internal/lsp"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

func SizeOfDataType(dataType DataType) uint {
	switch dataType {
	case DataTypeChar:
		return 1
	case DataTypeWideChar:
		return 2
	case DataTypeInteger:
		return 2
	case DataTypeSinteger:
		return 2
	case DataTypeLong:
		return 4
	case DataTypeSlong:
		return 4
	case DataTypeFloat:
		return 8
	case DataTypeDouble:
		return 8
	case DataTypeDev:
		return 6
	case DataTypeDevChan:
		return 8
	case DataTypeDevLev:
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

// IsInSubscriptExpression checks if a node is part of a subscript expression
// func IsInSubscriptExpression(node *tree_sitter.Node) bool {
// 	current := node
// 	for current != nil && current.Parent() != nil {
// 		parent := current.Parent()
// 		if parent.Kind() == "subscript_expression" {
// 			for i := uint(0); i < parent.NamedChildCount(); i++ {
// 				if parent.NamedChild(i).Equals(*current) ||
// 					(parent.ChildByFieldName("argument") != nil &&
// 						parent.ChildByFieldName("argument").Equals(*current)) {
// 					return true
// 				}
// 			}
// 		}
// 		current = parent
// 	}
// 	return false
// }

// GetBaseIdentifier gets the base identifier in a subscript expression
// func GetBaseIdentifier(node *tree_sitter.Node) *tree_sitter.Node {
// 	if node.Kind() == "identifier" {
// 		return node
// 	}

// 	current := node
// 	for current != nil && current.Kind() == "subscript_expression" {
// 		argument := current.ChildByFieldName("argument")
// 		if argument != nil {
// 			if argument.Kind() == "identifier" {
// 				return argument
// 			}
// 			current = argument
// 		} else {
// 			break
// 		}
// 	}

// 	return nil
// }

// GetOutermostSubscriptExpression finds the outermost subscript_expression
// func GetOutermostSubscriptExpression(node *tree_sitter.Node) *tree_sitter.Node {
// 	if node.Kind() != "subscript_expression" && node.Parent() == nil {
// 		return nil
// 	}

// 	var current *tree_sitter.Node

// 	// Find the first subscript expression
// 	temp := node
// 	for temp != nil {
// 		if temp.Kind() == "subscript_expression" {
// 			current = temp
// 			break
// 		}
// 		if temp.Parent() != nil {
// 			temp = temp.Parent()
// 		} else {
// 			break
// 		}
// 	}

// 	// Now find the outermost one
// 	for current != nil && current.Parent() != nil {
// 		parent := current.Parent()
// 		if parent.Kind() == "subscript_expression" {
// 			current = parent
// 		} else {
// 			break
// 		}
// 	}

// 	return current
// }

// IsSubscriptLeftHandSide checks if a subscript_expression is the left-hand side of an assignment
// func IsSubscriptLeftHandSide(assignmentNode, subscriptNode *tree_sitter.Node) bool {
// 	leftNode := assignmentNode.ChildByFieldName("left")
// 	if leftNode == nil {
// 		return false
// 	}

// 	// Check if leftNode is the subscriptNode or contains it
// 	return leftNode.Equals(*subscriptNode) || HasChildNode(leftNode, subscriptNode)
// }

// HasChildNode checks if parent contains child node
// func HasChildNode(parent, child *tree_sitter.Node) bool {
// 	for i := uint(0); i < parent.ChildCount(); i++ {
// 		if parent.Child(i).Equals(*child) {
// 			return true
// 		}
// 	}

// 	return false
// }

// FindParentOfType finds the nearest parent node of the specified type
// func FindParentOfType(node *tree_sitter.Node, parentType string) *tree_sitter.Node {
// 	current := node
// 	for current != nil && current.Parent() != nil {
// 		parent := current.Parent()
// 		if parent.Kind() == parentType {
// 			return parent
// 		}
// 		current = parent
// 	}
// 	return nil
// }
