package protocol

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ServerCapabilities struct {
	TextDocumentSync       TextDocumentSyncKind `json:"textDocumentSync,omitempty"`
	HoverProvider          bool                 `json:"hoverProvider"`
	DefinitionProvider     bool                 `json:"definitionProvider"`
	CodeActionProvider     bool                 `json:"codeActionProvider"`
	CompletionProvider     map[string]any       `json:"completionProvider"`
	DocumentSymbolProvider bool                 `json:"documentSymbolProvider"`
	DiagnosticProvider     *DiagnosticOptions   `json:"diagnosticProvider,omitempty"`
}

type DiagnosticOptions struct {
	InterFileDependencies bool `json:"interFileDependencies"`
	WorkspaceDiagnostics  bool `json:"workspaceDiagnostics"`
}
