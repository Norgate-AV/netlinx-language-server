package captures

import (
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

// SymbolInfo contains all extracted data from a symbol match
type SymbolInfo struct {
	Declaration *tree_sitter.Node
	Storage     *tree_sitter.Node
	Qualifier   *tree_sitter.Node
	Type        *tree_sitter.Node
	Identifier  *tree_sitter.Node
	Value       *tree_sitter.Node
	Size        *tree_sitter.Node
}

// FunctionInfo contains all extracted data from a function match
type FunctionInfo struct {
	Definition *tree_sitter.Node
	Name       *tree_sitter.Node
	ReturnType *tree_sitter.Node
	Parameters []*ParameterInfo
}

type ParameterInfo struct {
	Type *tree_sitter.Node
	Name *tree_sitter.Node
	Size *tree_sitter.Node
}
