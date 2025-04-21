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
// type Symbol struct {
// 	Name         string
// 	Type         SymbolType
// 	DataType     string // INTEGER, CHAR, etc.
// 	Value        string // Initial value if any
// 	VariableKind VariableKind
// 	Range        lsp.Range
// 	Section      *Section // Parent section
// 	Node         *tree_sitter.Node
// 	Children     []*Symbol // For hierarchical symbols
// }

type Symbol struct {
	Name    string
	Kind    SymbolKind
	Range   lsp.Range
	Node    *tree_sitter.Node
	Section string // Which section this symbol was declared in

	Storage   SymbolStorage   // stack_var, local_var
	Qualifier SymbolQualifier // constant, volatile, etc.
	DataType  SymbolDataType  // INTEGER, CHAR, etc.
	Value     string          // Initial value if any
	Scope     SymbolScope     // Scope of the symbol (e.g., global, local)

	Size       uint
	Dimensions uint

	Parameters []Parameter
	ReturnType SymbolDataType

	Referenced bool // Whether the symbol is referenced in the code

	// Is this symbol inside a preprocessor directive?
	Preprocessor bool
}

// Section represents a NetLinx section (DEFINE_DEVICE, etc.)
type Section struct {
	Kind  string
	Range lsp.Range
	Node  *tree_sitter.Node
}

// Document is the root of the semantic model
type Document struct {
	URI           string
	ProgramName   string
	Sections      []*Section
	Functions     []*Function
	Symbols       map[string]*Symbol       // All symbols by name
	SymbolsByKind map[SymbolType][]*Symbol // Symbols organized by kind
	Scopes        []*Scope
	GlobalScope   *Scope
}

// Scope represents a lexical scope
type Scope struct {
	Parent    *Scope
	Children  []*Scope
	Symbols   map[string]*Symbol
	Node      *tree_sitter.Node
	Range     lsp.Range
	ScopeType string // "global", "function", "if-block", etc.
}

// Function represents a function declaration
type Function struct {
	Name           string
	ReturnType     string
	Parameters     []*Parameter
	Body           *tree_sitter.Node
	Scope          *Scope
	Range          lsp.Range
	Node           *tree_sitter.Node
	IsEventHandler bool
}

// Parameter represents a function parameter
type Parameter struct {
	Name  string
	Type  string
	Range lsp.Range
	Node  *tree_sitter.Node
}

// Reference represents a reference to a symbol
type Reference struct {
	Symbol  *Symbol
	Node    *tree_sitter.Node
	Range   lsp.Range
	IsWrite bool // true if this reference modifies the symbol
}
