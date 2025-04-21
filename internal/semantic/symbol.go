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
	Name    string
	Kind    SymbolKind
	Range   lsp.Range
	Node    *tree_sitter.Node
	Section string // Which section this symbol was declared in

	StorageType StorageType // Volatile, Non-volatile, etc.
	DataType    DataType    // INTEGER, CHAR, etc.
	Value       string      // Initial value if any
	Scope       string      // Scope of the symbol (e.g., global, local)

	Size       uint
	Dimensions uint

	Parameters []Parameter
	ReturnType DataType

	Referenced bool // Whether the symbol is referenced in the code

	// Is this symbol inside a preprocessor directive?
	Preprocessor bool
}

func GetSymbolTable(tree *tree_sitter.Tree, content []byte) *SymbolTable {
	table := NewSymbolTable()

	query, err := queries.GetQuery("symbols.scm")
	if err != nil {
		fmt.Println("Error getting query:", err)
		return table
	}

	q, err := parser.CreateQuery(query)
	if err != nil {
		fmt.Println("Error creating query:", err)
		return table
	}

	defer q.Close()

	cursor := tree_sitter.NewQueryCursor()
	defer cursor.Close()

	matches := cursor.Matches(q, tree.RootNode(), nil)

	sp := NewSymbolProcessor(q, content, table)

	// Track current section
	var section string = ""

	for match := matches.Next(); match != nil; match = matches.Next() {
		if isSectionMatch(match, q) {
			section = getSectionType(match, q)
			// fmt.Printf("Entered Section: %s\n", strings.ToUpper(strings.TrimPrefix(section, "section.")))
			continue
		}

		sp.ProcessEntity(match, section)
	}

	return table
}
