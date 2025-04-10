package protocol

type (
	DocumentUri = string
	URI         = string
)

type LSPAny = any

type Position struct {
	Line      uint32 `json:"line"`
	Character uint32 `json:"character"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type MarkupContent struct {
	Kind  MarkupKind `json:"kind"`
	Value string     `json:"value"`
}

type MarkupKind = string

const (
	MarkupKindPlaintext MarkupKind = "plaintext"
	MarkupKindMarkdown  MarkupKind = "markdown"
)
