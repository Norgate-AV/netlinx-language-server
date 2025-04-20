package semantic

import (
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

	query, err := queries.GetQuery("device.scm")
	if err != nil {
		return table
	}

	q, err := parser.CreateQuery(query)
	if err != nil {
		return table
	}

	cursor := tree_sitter.NewQueryCursor()
	defer cursor.Close()

	matches := cursor.Matches(q, tree.RootNode(), nil)

	// Track current section
	var currentSection string

	for match := matches.Next(); match != nil; match = matches.Next() {
		for _, capture := range match.Captures {
			// fmt.Printf(
			// 	"Match %d, Capture %d (%s): %s\n",
			// 	match.PatternIndex,
			// 	capture.Index,
			// 	q.CaptureNames()[capture.Index],
			// 	capture.Node.Utf8Text(content),
			// )
			captureName := q.CaptureNames()[capture.Index]
			if captureName == "device_section" {
				currentSection = SectionDefineDevice
				break
			}
		}

		// Skip processing if we're not in a device section
		if currentSection != SectionDefineDevice {
			continue
		}

		// Extract identifier, value, qualifier, and type
		var identNode tree_sitter.Node
		var valueNode tree_sitter.Node
		var qualifierNode tree_sitter.Node
		var typeNode tree_sitter.Node

		for _, capture := range match.Captures {
			captureName := q.CaptureNames()[capture.Index]
			switch captureName {
			case "identifier":
				identNode = capture.Node
			case "value":
				valueNode = capture.Node
			case "qualifier":
				qualifierNode = capture.Node
			case "type":
				typeNode = capture.Node
			}
		}

		// Create and add the symbol
		table.AddSymbol(&Symbol{
			Name:        identNode.Utf8Text(content),
			Kind:        SymbolKindDevice,
			Range:       GetNodeRange(&identNode),
			Node:        &identNode,
			StorageType: qualifierNode.Utf8Text(content),
			DataType:    typeNode.Utf8Text(content),
			Value:       valueNode.Utf8Text(content),
			Size:        SizeOfDataType(typeNode.Utf8Text(content)),
			Dimensions:  0,
		})
	}

	return table
}

func processNode(node *tree_sitter.Node, content []byte, section *string, table *SymbolTable) {
	if node == nil {
		return
	}

	kind := node.Kind()

	switch kind {
	case "define_device_section":
		*section = SectionDefineDevice
	case "define_constant_section":
		*section = SectionDefineConstant
	case "define_variable_section":
		*section = SectionDefineVariable
	}

	if kind != "identifier" {
		return
	}

	name := node.Utf8Text(content)
	parent := node.Parent()

	if parent != nil && *section != "" {
		switch *section {
		case SectionDefineDevice:
			processDefineDevice(node, parent, name, content, table)
			// case SectionDefineConstant:
			// 	processDefineConstant(node, parent, name, content, table)
		}
	}
}

func processDefineDevice(node *tree_sitter.Node, parent *tree_sitter.Node, name string, content []byte, table *SymbolTable) {
	// In DEFINE_DEVICE section, device definitions are parsed
	// as assignment expressions.
	// Eg. identifier = device_literal | expression
	// Eg. dvTP = 10001:1:0

	if parent.Kind() == "assignment_expression" {
		if !IsLeftHandSide(parent, node) {
			return
		}

		table.AddSymbol(&Symbol{
			Name:        name,
			Kind:        SymbolKindDevice,
			Range:       GetNodeRange(node),
			Node:        node,
			StorageType: StorageTypeConstant,
			DataType:    DataTypeDev,
			Value:       GetNodeValue(parent, content),
			Size:        SizeOfDataType(DataTypeDev),
			Dimensions:  0,
		})
	}
}

// func processDefineConstant(node *tree_sitter.Node, parent *tree_sitter.Node, name string, content []byte, table *SymbolTable) {
// 	// In DEFINE_CONSTANT section, constants can be parsed
// 	// as either a declaration or an assignment expression
// 	// depending on how it is declared.
// 	// StorageType and DataType are optional
// 	// If not specified, it will be parsed as an assignment expression
// 	// Otherwise, it will be parsed as a declaration
// 	// StorageType for a constant can only be "constant" if specified
// 	// If DataType is not specified, it is implicitly INTEGER for non array values
// 	// Otherwise, it is implicitly CHAR for array values
// 	// Eg. [constant] [data_type] identifier = value
// 	// Eg. [constant] [data_type] identifier[[size]] = value

// 	switch parent.Kind() {
// 	case "assignment_expression":
// 		if !IsLeftHandSide(parent, node) {
// 			return
// 		}

// 		// Always parse as a constant
// 		storageType := StorageTypeConstant
// 		dataType := DataTypeInteger
// 		dimensions := uint(0)
// 		size := SizeOfDataType(dataType)

// 		table.AddSymbol(&Symbol{
// 			Name:        name,
// 			Kind:        SymbolKindConstant,
// 			Range:       GetNodeRange(node),
// 			Node:        node,
// 			StorageType: storageType,
// 			DataType:    dataType,
// 			Value:       GetNodeValue(parent, content),
// 			Size:        size,
// 			Dimensions:  dimensions,
// 		})
// 	case "subscript_expression":
// 		// In this case, the identifier is part of an array declaration
// 		// Eg. identifier[[size]] = value
// 		// We can ignore this case as it will be handled in the array declaration
// 		baseIdNode := GetBaseIdentifier(parent)

// 		if baseIdNode != nil && baseIdNode.Kind() == "identifier" {
// 			// Find the outermost subscript expression
// 			outerSubscript := GetOutermostSubscriptExpression(parent)

// 			// Find the assignment expression that contains this subscript
// 			assignmentNode := FindParentOfType(outerSubscript, "assignment_expression")

// 			if assignmentNode != nil && IsSubscriptLeftHandSide(assignmentNode, outerSubscript) {
// 				baseNodeName := baseIdNode.Utf8Text(content)

// 				// Now create the symbol using the base identifier
// 				storageType := StorageTypeConstant
// 				dataType := DataTypeChar
// 				dimensions := GetArrayDimensions(outerSubscript)
// 				size := SizeOfDataType(dataType) * dimensions

// 				table.AddSymbol(&Symbol{
// 					Name:        baseNodeName, // Use the base identifier name!
// 					Kind:        SymbolKindConstant,
// 					Range:       GetNodeRange(baseIdNode),
// 					Node:        baseIdNode,
// 					StorageType: storageType,
// 					DataType:    dataType,
// 					Value:       GetNodeValue(assignmentNode, content),
// 					Size:        size,
// 					Dimensions:  dimensions,
// 				})
// 			}
// 		}
// 	case "declaration":
// 		// Parse declaration here
// 	}
// }

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
