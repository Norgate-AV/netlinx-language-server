package semantic

import (
	"github.com/Norgate-AV/netlinx-language-server/internal/lsp"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

// SymbolType represents the semantic type of a symbol
type SymbolType int

const (
	UnknownSymbol SymbolType = iota
	DeviceSymbol
	ConstantSymbol
	VariableSymbol
	FunctionSymbol
	ParameterSymbol
	// Add other symbol types as needed
)

// VariableKind represents the kind of variable (volatile, non-volatile, etc.)
type VariableKind int

const (
	DefaultVar VariableKind = iota
	VolatileVar
	NonVolatileVar
)

// Symbol represents a NetLinx symbol with semantic information
type Symbol struct {
	Name         string
	Type         SymbolType
	DataType     string // INTEGER, CHAR, etc.
	Value        string // Initial value if any
	VariableKind VariableKind
	Range        lsp.Range
	Section      *Section // Parent section
	Node         *tree_sitter.Node
	Children     []*Symbol // For hierarchical symbols
}

// Section represents a NetLinx section (DEFINE_DEVICE, etc.)
type Section struct {
	Kind  string
	Range lsp.Range
	Node  *tree_sitter.Node
}

// Document represents semantically analyzed NetLinx document
type Document struct {
	URI      string
	Sections []*Section
	Symbols  map[string]*Symbol // All symbols by name
	// Could add more document-level semantic info
}
