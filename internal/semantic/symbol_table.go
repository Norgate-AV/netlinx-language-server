package semantic

import (
	"github.com/Norgate-AV/netlinx-language-server/internal/logger"
	"github.com/Norgate-AV/netlinx-language-server/internal/parser"
	"github.com/Norgate-AV/netlinx-language-server/internal/queries"
	"github.com/sirupsen/logrus"

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

func GetSymbolTable(tree *tree_sitter.Tree, content []byte, log logger.Logger) *SymbolTable {
	table := NewSymbolTable()

	query, err := queries.GetQuery("symbols.scm")
	if err != nil {
		log.Error("Error getting query", logrus.Fields{
			"error": err.Error(),
		})

		return table
	}

	q, err := parser.CreateQuery(query)
	if err != nil {
		log.Error("Error creating query", logrus.Fields{
			"error": err.Error(),
		})

		return table
	}

	defer q.Close()

	cursor := tree_sitter.NewQueryCursor()
	defer cursor.Close()

	matches := cursor.Matches(q, tree.RootNode(), nil)

	sp := NewSymbolProcessor(q, content, table, log)

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

func GetSymbolTable2(matches tree_sitter.QueryMatches, sp *SymbolProcessor) *SymbolTable {
	table := NewSymbolTable()

	// Track current section
	var section string = ""

	for match := matches.Next(); match != nil; match = matches.Next() {
		if isSectionMatch(match, sp.query) {
			section = getSectionType(match, sp.query)
			continue
		}

		sp.ProcessEntity(match, section)
	}

	return table
}
