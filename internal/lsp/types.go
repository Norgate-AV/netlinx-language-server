package lsp

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

type Command struct {
	Title     string   `json:"title"`
	Command   string   `json:"command"`
	Arguments []LSPAny `json:"arguments,omitempty"`
}

type WorkspaceEdit struct {
	Changes         map[DocumentUri][]TextEdit `json:"changes,omitempty"`
	DocumentChanges []TextDocumentEdit         `json:"documentChanges,omitempty"`
}

type TextEdit struct {
	Range   Range  `json:"range"`
	NewText string `json:"newText"`
}

type TextDocumentEdit struct {
	TextDocument VersionedTextDocumentIdentifier `json:"textDocument"`
	Edits        []TextEdit                      `json:"edits"`
}

type FormattingOptions struct {
	TabSize      int  `json:"tabSize"`
	InsertSpaces bool `json:"insertSpaces"`
}

type RenameOptions struct {
	PrepareProvider bool `json:"prepareProvider,omitempty"`
}

type WorkspaceFolder struct {
	URI  DocumentUri `json:"uri"`
	Name string      `json:"name"`
}

type WorkDoneProgressParams struct {
	WorkDoneToken LSPAny `json:"workDoneToken,omitempty"`
}
