package semantic

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
