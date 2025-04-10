package protocol

type DocumentHighlight struct {
	Range Range `json:"range"`
	Kind  int   `json:"kind,omitempty"`
}

const (
	DocumentHighlightKindText  = 1
	DocumentHighlightKindRead  = 2
	DocumentHighlightKindWrite = 3
)

type DocumentHighlightParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}
