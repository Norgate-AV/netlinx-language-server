package semantic

import (
	"errors"
	"strings"

	"github.com/Norgate-AV/netlinx-language-server/internal/captures"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type SymbolProcessor struct {
	query   *tree_sitter.Query
	content []byte
	table   *SymbolTable
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

func (sp *SymbolProcessor) extractSymbolInfo(match *tree_sitter.QueryMatch) captures.SymbolInfo {
	var info captures.SymbolInfo

	for _, capture := range match.Captures {
		name := sp.query.CaptureNames()[capture.Index]
		node := capture.Node

		switch name {
		case captures.CaptureSymbolDeclaration:
			info.Declaration = &node
		case captures.CaptureSymbolStorage:
			info.Storage = &node
		case captures.CaptureSymbolQualifier:
			info.Qualifier = &node
		case captures.CaptureSymbolType:
			info.Type = &node
		case captures.CaptureSymbolIdentifier:
			info.Identifier = &node
		case captures.CaptureSymbolSize:
			info.Size = &node
		case captures.CaptureSymbolValue:
			info.Value = &node
		}
	}

	return info
}

func (sp *SymbolProcessor) createSymbol(info captures.SymbolInfo, section string) (*Symbol, error) {
	var kind SymbolKind
	var storage SymbolStorage
	var qualifier SymbolQualifier
	var dataType SymbolDataType
	var dimensions uint = 0

	if info.Declaration == nil {
		return nil, errors.New("declaration node is nil")
	}

	if info.Identifier == nil {
		return nil, errors.New("identifier node is nil")
	}

	// Set default values based on the section
	switch section {
	case captures.CaptureSectionDefineDevice:
		kind = SymbolKindDevice
		qualifier = SymbolQualifierConstant
		dataType = SymbolDataTypeDev
	case captures.CaptureSectionDefineConstant:
		kind = SymbolKindConstant
		qualifier = SymbolQualifierConstant

		// Data type is implicitly an integer for non-array types
		// For array types the data type is implicitly a char array
		if IsArray(info.Declaration) {
			dataType = SymbolDataTypeChar
			dimensions = GetArrayDimensions(info.Declaration)
		} else {
			dataType = SymbolDataTypeInteger
		}

	case captures.CaptureSectionDefineVariable:
		kind = SymbolKindVariable
		qualifier = SymbolQualifierNonVolatile

		// Data type is implicitly an integer for non-array types
		// For array types the data type is implicitly a char array
		if IsArray(info.Declaration) {
			dataType = SymbolDataTypeChar
			dimensions = GetArrayDimensions(info.Declaration)
		} else {
			dataType = SymbolDataTypeInteger
		}
	}

	if info.Storage != nil {
		storage = info.Storage.Utf8Text(sp.content)
	}

	if info.Qualifier != nil {
		qualifier = info.Qualifier.Utf8Text(sp.content)
	}

	if info.Type != nil {
		dataType = info.Type.Utf8Text(sp.content)
	}

	value := ""
	if info.Value != nil {
		value = info.Value.Utf8Text(sp.content)
	}

	size := SizeOfDataType(dataType)

	return &Symbol{
		Name:       info.Identifier.Utf8Text(sp.content),
		Kind:       kind,
		Range:      GetNodeRange(info.Identifier),
		Node:       info.Identifier,
		Section:    section,
		Storage:    storage,
		Qualifier:  qualifier,
		DataType:   dataType,
		Value:      value,
		Size:       size,
		Dimensions: dimensions,
		Scope:      SymbolScopeGlobal,
	}, nil
}
