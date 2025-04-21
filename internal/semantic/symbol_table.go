package semantic

import (
	"fmt"

	"github.com/Norgate-AV/netlinx-language-server/internal/parser"
	"github.com/Norgate-AV/netlinx-language-server/internal/queries"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

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

	query, err := queries.GetQuery("symbols.scm")
	if err != nil {
		// Should be using my logger here
		fmt.Println("Error getting query:", err)
		return table
	}

	q, err := parser.CreateQuery(query)
	if err != nil {
		// Should be using my logger here
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
			continue
		}

		sp.ProcessEntity(match, section)
	}

	return table
}
