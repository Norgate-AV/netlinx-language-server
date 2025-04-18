package semantic

import (
	"github.com/Norgate-AV/netlinx-language-server/internal/lsp"

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
	DateType    DataType    // INTEGER, CHAR, etc.
	Value       string      // Initial value if any

	Size       uint
	Dimensions uint
}

type SymbolTable struct {
	Symbols map[string]*Symbol
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		Symbols: make(map[string]*Symbol),
	}
}

func (st *SymbolTable) AddSymbol(symbol *Symbol) {
	st.Symbols[symbol.Name] = symbol
}

func GetSymbolTable(root *tree_sitter.Node, content []byte) *SymbolTable {
	table := NewSymbolTable()

	var section string

	var walk func(node *tree_sitter.Node)
	walk = func(node *tree_sitter.Node) {
		if node == nil {
			return
		}

		kind := node.Kind()

		switch kind {
		case "define_device_section":
			section = SectionDefineDevice
		case "define_constant_section":
			section = SectionDefineConstant
		case "define_variable_section":
			section = SectionDefineVariable
		}

		if kind == "identifier" {
			name := node.Utf8Text(content)
			parent := node.Parent()

			if parent != nil && section != "" {
				switch section {
				case SectionDefineDevice:
					// In DEFINE_DEVICE section, device definitions are parsed
					// as assignment expressions.
					// Eg. identifier = device_literal | expression
					// Eg. dvTP = 10001:1:0
					if parent.Kind() == "assignment_expression" {
						if IsLeftHandSide(parent, node) {
							table.AddSymbol(&Symbol{
								Name:        name,
								Kind:        SymbolKindDevice,
								Range:       GetNodeRange(node),
								Node:        node,
								StorageType: StorageTypeConstant,
								DateType:    DataTypeDev,
								Value:       GetNodeValue(parent, content),
								Size:        SizeOfDataType(DataTypeDev),
								Dimensions:  0,
							})
						}
					}

				case SectionDefineConstant:
					// In DEFINE_CONSTANT section, constants can be parsed
					// as either a declaration or an assignment expression
					// depending on how it is declared.
					// StorageType and DataType are optional
					// If not specified, it will be parsed as an assignment expression
					// Otherwise, it will be parsed as a declaration
					// StorageType for a constant can only be "constant" if specified
					// If DataType is not specified, it is implicitly INTEGER for non array values
					// Otherwise, it is implicitly CHAR for array values
					// Eg. [constant] [data_type] identifier = value
					// Eg. [constant] [data_type] identifier[[size]] = value

					switch parent.Kind() {
					case "assignment_expression":
						if IsLeftHandSide(parent, node) {
							// Always parse as a constant
							storageType := StorageTypeConstant
							dataType := DataTypeInteger

							table.AddSymbol(&Symbol{
								Name:        name,
								Kind:        SymbolKindConstant,
								Range:       GetNodeRange(node),
								Node:        node,
								StorageType: storageType,
								DateType:    dataType,
								Value:       GetNodeValue(parent, content),
								Size:        SizeOfDataType(dataType),
								Dimensions:  0,
							})
						}
					case "declaration":
						// Parse declaration here
					}
				}
			}
		}

		for i := uint(0); i < node.ChildCount(); i++ {
			walk(node.Child(i))
		}
	}

	walk(root)
	return table
}

const (
	SectionDefineDevice   string = "define_device_section"
	SectionDefineConstant string = "define_constant_section"
	SectionDefineVariable string = "define_variable_section"
)
