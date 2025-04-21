package semantic

import (
	"errors"
	"strings"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type SymbolProcessor struct {
	query   *tree_sitter.Query
	content []byte
	table   *SymbolTable
}

// SymbolInfo contains all extracted data from a symbol match
type SymbolInfo struct {
	Declaration *tree_sitter.Node
	Identifier  *tree_sitter.Node
	Qualifier   *tree_sitter.Node
	Type        *tree_sitter.Node
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

func NewSymbolProcessor(query *tree_sitter.Query, content []byte, table *SymbolTable) *SymbolProcessor {
	return &SymbolProcessor{
		query:   query,
		content: content,
		table:   table,
	}
}

func (sp *SymbolProcessor) ProcessEntity(match *tree_sitter.QueryMatch, section string) {
	if len(match.Captures) == 0 {
		return
	}

	// Get the first capture
	capture := sp.query.CaptureNames()[match.Captures[0].Index]

	switch {
	case strings.HasPrefix(capture, "symbol."):
		sp.processSymbol(match, section)
		// case strings.HasPrefix(capture, "function."):
		// 	sp.processFunction(match, section)
		// case strings.HasPrefix(capture, "type."):
		// 	sp.processType(match, section)
	}
}

func (sp *SymbolProcessor) processSymbol(match *tree_sitter.QueryMatch, section string) {
	info := sp.extractSymbolInfo(match)

	if info.Identifier == nil {
		return
	}

	symbol, err := sp.createSymbol(info, section)
	if err != nil {
		return
	}

	sp.table.AddSymbol(symbol)
}

func (sp *SymbolProcessor) extractSymbolInfo(match *tree_sitter.QueryMatch) SymbolInfo {
	var info SymbolInfo

	for _, capture := range match.Captures {
		name := sp.query.CaptureNames()[capture.Index]
		node := capture.Node

		switch name {
		case CaptureSymbolDeclaration:
			info.Declaration = &node
		case CaptureSymbolQualifier:
			info.Qualifier = &node
		case CaptureSymbolType:
			info.Type = &node
		case CaptureSymbolIdentifier:
			info.Identifier = &node
		case CaptureSymbolSize:
			info.Size = &node
		case CaptureSymbolValue:
			info.Value = &node
		}
	}

	return info
}

func (sp *SymbolProcessor) createSymbol(info SymbolInfo, section string) (*Symbol, error) {
	var kind SymbolKind
	var storageType StorageType
	var dataType DataType
	var dimensions uint = 0

	if info.Declaration == nil {
		return nil, errors.New("declaration node is nil")
	}

	if info.Identifier == nil {
		return nil, errors.New("identifier node is nil")
	}

	// Set default values based on the section
	switch section {
	case CaptureSectionDefineDevice:
		kind = SymbolKindDevice
		storageType = StorageTypeConstant
		dataType = DataTypeDev
	case CaptureSectionDefineConstant:
		kind = SymbolKindConstant
		storageType = StorageTypeConstant

		// Data type is implicitly an integer for non-array types
		// For array types the data type is implicitly a char array
		if IsArray(info.Declaration) {
			dataType = DataTypeChar
			dimensions = GetArrayDimensions(info.Declaration)
		} else {
			dataType = DataTypeInteger
		}

	case CaptureSectionDefineVariable:
		kind = SymbolKindVariable
		storageType = StorageTypeNonVolatile

		// Data type is implicitly an integer for non-array types
		// For array types the data type is implicitly a char array
		if IsArray(info.Declaration) {
			dataType = DataTypeChar
			dimensions = GetArrayDimensions(info.Declaration)
		} else {
			dataType = DataTypeInteger
		}
	}

	if info.Qualifier != nil {
		storageType = info.Qualifier.Utf8Text(sp.content)
	}

	if info.Type != nil {
		dataType = info.Type.Utf8Text(sp.content)
	}

	value := ""
	if info.Value != nil {
		value = info.Value.Utf8Text(sp.content)
	}

	size := SizeOfDataType(strings.ToLower(dataType))

	return &Symbol{
		Name:        info.Identifier.Utf8Text(sp.content),
		Kind:        kind,
		Range:       GetNodeRange(info.Identifier),
		Node:        info.Identifier,
		Section:     section,
		StorageType: storageType,
		DataType:    dataType,
		Value:       value,
		Size:        size,
		Dimensions:  dimensions,
	}, nil
}
