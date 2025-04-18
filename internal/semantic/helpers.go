package semantic

import (
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
