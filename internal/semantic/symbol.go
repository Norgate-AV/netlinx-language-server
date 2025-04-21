package semantic

import (
	"fmt"

	"github.com/Norgate-AV/netlinx-language-server/internal/lsp"
	"github.com/Norgate-AV/netlinx-language-server/parser"
	"github.com/Norgate-AV/netlinx-language-server/queries"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type StorageType = string

const (
	StorageTypeConstant    StorageType = "constant"
	StorageTypeVolatile    StorageType = "volatile"
	StorageTypeNonVolatile StorageType = "non_volatile"
	StorageTypePersistent  StorageType = "persistent"
)

type SymbolKind = string

const (
	SymbolKindDevice   SymbolKind = "device"
	SymbolKindConstant SymbolKind = "constant"
	SymbolKindVariable SymbolKind = "variable"
	SymbolKindFunction SymbolKind = "function"
	SymbolKindStruct   SymbolKind = "struct"
)

type DataType = string

const (
	DataTypeChar     DataType = "char"
	DataTypeWideChar DataType = "widechar"
	DataTypeInteger  DataType = "integer"
	DataTypeSinteger DataType = "sinteger"
	DataTypeLong     DataType = "long"
	DataTypeSlong    DataType = "slong"
	DataTypeFloat    DataType = "float"
	DataTypeDouble   DataType = "double"
	DataTypeDev      DataType = "dev"
	DataTypeDevChan  DataType = "devchan"
	DataTypeDevLev   DataType = "devlev"
)

type Symbol struct {
	Name  string
	Kind  SymbolKind
	Range lsp.Range
	Node  *tree_sitter.Node

	StorageType StorageType // Volatile, Non-volatile, etc.
	DataType    DataType    // INTEGER, CHAR, etc.
	Value       string      // Initial value if any
	Scope       string      // Scope of the symbol (e.g., global, local)

	Size       uint
	Dimensions uint

	Parameters []Parameter
	ReturnType DataType

	Referenced bool // Whether the symbol is referenced in the code
}

type SymbolTable struct {
	Symbols map[string][]*Symbol
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		Symbols: make(map[string][]*Symbol),
	}
}

func (st *SymbolTable) AddSymbol(symbol *Symbol) {
	st.Symbols[symbol.Name] = append(st.Symbols[symbol.Name], symbol)
}

func GetSymbolTable(tree *tree_sitter.Tree, content []byte) *SymbolTable {
	table := NewSymbolTable()

	query, err := queries.GetQuery("symbols2.scm")
	if err != nil {
		fmt.Println("Error getting query:", err)
		return table
	}

	q, err := parser.CreateQuery(query)
	if err != nil {
		fmt.Println("Error creating query:", err)
		return table
	}

	count := GetMatchCount(q, tree)
	if count == 0 {
		fmt.Println("No matches found")
		return table
	}

	fmt.Println("Match count: ", count)

	cursor := tree_sitter.NewQueryCursor()
	defer cursor.Close()

	matches := cursor.Matches(q, tree.RootNode(), nil)
	// captures := cursor.Captures(q, tree.RootNode(), nil)
	// for match, index := captures.Next(); match != nil; match, index = captures.Next() {
	// 	fmt.Printf(
	// 		"Capture %d: %s\n",
	// 		index,
	// 		match.Captures[index].Node.Utf8Text(content),
	// 	)
	// }
	// Track current section
	// var currentSection string

	for match := matches.Next(); match != nil; match = matches.Next() {
		for _, capture := range match.Captures {
			fmt.Printf(
				"Match %d, Capture %d (%s): %s\n",
				match.PatternIndex,
				capture.Index,
				q.CaptureNames()[capture.Index],
				capture.Node.Utf8Text(content),
			)

			// captureName := q.CaptureNames()[capture.Index]
			// if captureName == "device_section" {
			// 	currentSection = SectionDefineDevice
			// 	break
			// }

			// switch q.CaptureNames()[capture.Index] {
			// case "define_device_section":
		}

		// Skip processing if we're not in a device section
		// if currentSection != SectionDefineDevice {
		// 	continue
		// }

		// Extract identifier, value, qualifier, and type
		var identNode *tree_sitter.Node
		var valueNode *tree_sitter.Node
		var qualifierNode *tree_sitter.Node
		var typeNode *tree_sitter.Node

		for _, capture := range match.Captures {
			captureName := q.CaptureNames()[capture.Index]
			switch captureName {
			case "identifier":
				identNode = &capture.Node
			case "value":
				valueNode = &capture.Node
			case "qualifier":
				qualifierNode = &capture.Node
			case "type":
				typeNode = &capture.Node
			}
		}

		// Skip if any of the nodes are nil
		if identNode == nil {
			continue
		}

		// Set default values
		storageType := StorageTypeConstant
		value := ""
		dataType := DataTypeDev

		if valueNode != nil {
			value = valueNode.Utf8Text(content)
		}

		if qualifierNode != nil {
			storageType = qualifierNode.Utf8Text(content)
		}

		if typeNode != nil {
			dataType = typeNode.Utf8Text(content)
		}

		// Create and add the symbol
		table.AddSymbol(&Symbol{
			Name:        identNode.Utf8Text(content),
			Kind:        SymbolKindDevice,
			Range:       GetNodeRange(identNode),
			Node:        identNode,
			StorageType: storageType,
			DataType:    dataType,
			Value:       value,
			Size:        SizeOfDataType(dataType),
			Dimensions:  0,
		})
	}

	return table
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

const (
	SectionDefineDevice   string = "define_device_section"
	SectionDefineConstant string = "define_constant_section"
	SectionDefineVariable string = "define_variable_section"
)

// IsInSubscriptExpression checks if a node is part of a subscript expression
func IsInSubscriptExpression(node *tree_sitter.Node) bool {
	current := node
	for current != nil && current.Parent() != nil {
		parent := current.Parent()
		if parent.Kind() == "subscript_expression" {
			for i := uint(0); i < parent.NamedChildCount(); i++ {
				if parent.NamedChild(i).Equals(*current) ||
					(parent.ChildByFieldName("argument") != nil &&
						parent.ChildByFieldName("argument").Equals(*current)) {
					return true
				}
			}
		}
		current = parent
	}
	return false
}

// GetBaseIdentifier gets the base identifier in a subscript expression
func GetBaseIdentifier(node *tree_sitter.Node) *tree_sitter.Node {
	if node.Kind() == "identifier" {
		return node
	}

	current := node
	for current != nil && current.Kind() == "subscript_expression" {
		argument := current.ChildByFieldName("argument")
		if argument != nil {
			if argument.Kind() == "identifier" {
				return argument
			}
			current = argument
		} else {
			break
		}
	}

	return nil
}

// GetOutermostSubscriptExpression finds the outermost subscript_expression
func GetOutermostSubscriptExpression(node *tree_sitter.Node) *tree_sitter.Node {
	if node.Kind() != "subscript_expression" && node.Parent() == nil {
		return nil
	}

	var current *tree_sitter.Node

	// Find the first subscript expression
	temp := node
	for temp != nil {
		if temp.Kind() == "subscript_expression" {
			current = temp
			break
		}
		if temp.Parent() != nil {
			temp = temp.Parent()
		} else {
			break
		}
	}

	// Now find the outermost one
	for current != nil && current.Parent() != nil {
		parent := current.Parent()
		if parent.Kind() == "subscript_expression" {
			current = parent
		} else {
			break
		}
	}

	return current
}

// IsSubscriptLeftHandSide checks if a subscript_expression is the left-hand side of an assignment
func IsSubscriptLeftHandSide(assignmentNode, subscriptNode *tree_sitter.Node) bool {
	leftNode := assignmentNode.ChildByFieldName("left")
	if leftNode == nil {
		return false
	}

	// Check if leftNode is the subscriptNode or contains it
	return leftNode.Equals(*subscriptNode) || HasChildNode(leftNode, subscriptNode)
}

// HasChildNode checks if parent contains child node
func HasChildNode(parent, child *tree_sitter.Node) bool {
	for i := uint(0); i < parent.ChildCount(); i++ {
		if parent.Child(i).Equals(*child) {
			return true
		}
	}

	return false
}

// FindParentOfType finds the nearest parent node of the specified type
func FindParentOfType(node *tree_sitter.Node, parentType string) *tree_sitter.Node {
	current := node
	for current != nil && current.Parent() != nil {
		parent := current.Parent()
		if parent.Kind() == parentType {
			return parent
		}
		current = parent
	}
	return nil
}
