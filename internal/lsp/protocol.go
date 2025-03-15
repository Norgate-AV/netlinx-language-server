package lsp

// TextDocumentSyncKind represents the sync capabilities
type TextDocumentSyncKind int

const (
	TextDocumentSyncKindNone TextDocumentSyncKind = iota
	TextDocumentSyncKindFull
	TextDocumentSyncKindIncremental
)

// Position in a text document
type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

// Range in a text document
type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

// ClientInfo represents information about the client
type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

// InitializeParams represents the parameters sent from client during initialization
type InitializeParams struct {
	ProcessID             int        `json:"processId"`
	ClientInfo            ClientInfo `json:"clientInfo,omitempty"`
	RootURI               string     `json:"rootUri"`
	InitializationOptions interface{} `json:"initializationOptions,omitempty"`
	Capabilities          interface{} `json:"capabilities"`
}

// CompletionOptions represents the server's completion capabilities
type CompletionOptions struct {
	ResolveProvider   bool     `json:"resolveProvider"`
	TriggerCharacters []string `json:"triggerCharacters,omitempty"`
}

// ServerCapabilities represents what the language server can do
type ServerCapabilities struct {
	TextDocumentSync    TextDocumentSyncKind `json:"textDocumentSync"`
	CompletionProvider  *CompletionOptions   `json:"completionProvider,omitempty"`
	HoverProvider       bool                 `json:"hoverProvider"`
}

// InitializeResult represents the result of the initialize request
type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
}

// TextDocumentItem represents an opened document
type TextDocumentItem struct {
	URI        string `json:"uri"`
	LanguageID string `json:"languageId"`
	Version    int    `json:"version"`
	Text       string `json:"text"`
}

// TextDocumentIdentifier identifies a document
type TextDocumentIdentifier struct {
	URI string `json:"uri"`
}

// VersionedTextDocumentIdentifier identifies a document with version info
type VersionedTextDocumentIdentifier struct {
	TextDocumentIdentifier
	Version int `json:"version"`
}

// DidOpenTextDocumentParams parameters for document open notification
type DidOpenTextDocumentParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}

// TextDocumentContentChangeEvent represents a change to a text document
type TextDocumentContentChangeEvent struct {
	Range       *Range `json:"range,omitempty"`
	RangeLength int    `json:"rangeLength,omitempty"`
	Text        string `json:"text"`
}

// DidChangeTextDocumentParams parameters for document change notification
type DidChangeTextDocumentParams struct {
	TextDocument   VersionedTextDocumentIdentifier  `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

// DidCloseTextDocumentParams parameters for document close notification
type DidCloseTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// CompletionParams parameters for completion request
type CompletionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}

// CompletionItemKind defines the type of a completion item
type CompletionItemKind int

// HoverParams parameters for hover request
type HoverParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}

// MarkedString represents a string that is marked as code
type MarkedString struct {
	Language string `json:"language"`
	Value    string `json:"value"`
}

// Hover represents the result of a hover request
type Hover struct {
	Contents []MarkedString `json:"contents"`
	Range    *Range         `json:"range,omitempty"`
}
