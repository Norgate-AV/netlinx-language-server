package protocol

type TextDocumentUri = string

type Position struct {
	Line      uint32 `json:"line"`
	Character uint32 `json:"character"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type TextDocumentItem struct {
	URI        TextDocumentUri `json:"uri"`
	LanguageID string          `json:"languageId"`
	Version    int             `json:"version"`
	Text       string          `json:"text"`
}

type InitializeRequest struct {
	Params InitializeRequestParams `json:"params"`
}

type InitializeRequestParams struct {
	ClientInfo *ClientInfo `json:"clientInfo"`
	// ... there's tons more that goes here
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResponse struct {
	Result InitializeResult `json:"result"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverInfo"`
}

type ServerCapabilities struct {
	TextDocumentSync int  `json:"textDocumentSync,omitempty"`
	HoverProvider    bool `json:"hoverProvider"`
	// DefinitionProvider bool           `json:"definitionProvider"`
	// CodeActionProvider bool           `json:"codeActionProvider"`
	// CompletionProvider map[string]any `json:"completionProvider"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type DidChangeTextDocumentParams struct {
	TextDocument   VersionedTextDocumentIdentifier  `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

type VersionedTextDocumentIdentifier struct {
	URI     TextDocumentUri `json:"uri"`
	Version int             `json:"version"`
}

type TextDocumentContentChangeEvent struct {
	Range       *Range `json:"range,omitempty"`
	RangeLength *int   `json:"rangeLength,omitempty"`
	Text        string `json:"text"`
}

// DidCloseTextDocumentParams for textDocument/didClose
type DidCloseTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// HoverParams for textDocument/hover
type HoverParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}

// Hover response
type Hover struct {
	Contents MarkupContent `json:"contents"`
	Range    *Range        `json:"range,omitempty"`
}

// MarkupContent for formatted documentation
type MarkupContent struct {
	Kind  MarkupKind `json:"kind"`
	Value string     `json:"value"`
}

type MarkupKind string

const (
	MarkupKindPlaintext MarkupKind = "plaintext"
	MarkupKindMarkdown  MarkupKind = "markdown"
)
