package lsp

type SymbolInformation struct {
	Name          string   `json:"name"`
	Kind          int      `json:"kind"`
	Location      Location `json:"location"`
	ContainerName string   `json:"containerName,omitempty"`
}

type DocumentSymbol struct {
	Name           string           `json:"name"`
	Detail         string           `json:"detail,omitempty"`
	Kind           int              `json:"kind"`
	Deprecated     bool             `json:"deprecated,omitempty"`
	Range          Range            `json:"range"`
	SelectionRange Range            `json:"selectionRange"`
	Children       []DocumentSymbol `json:"children,omitempty"`
}

type DocumentSymbolParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type SymbolKind = int

const (
	SymbolKindFile        SymbolKind = 1
	SymbolKindModule      SymbolKind = 2
	SymbolKindNamespace   SymbolKind = 3
	SymbolKindPackage     SymbolKind = 4
	SymbolKindClass       SymbolKind = 5
	SymbolKindMethod      SymbolKind = 6
	SymbolKindProperty    SymbolKind = 7
	SymbolKindField       SymbolKind = 8
	SymbolKindConstructor SymbolKind = 9
	SymbolKindEnum        SymbolKind = 10
	SymbolKindInterface   SymbolKind = 11
	SymbolKindFunction    SymbolKind = 12
	SymbolKindVariable    SymbolKind = 13
	SymbolKindConstant    SymbolKind = 14
)
