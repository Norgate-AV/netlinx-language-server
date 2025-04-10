package lsp

type CodeAction struct {
	Title       string         `json:"title"`
	Kind        string         `json:"kind,omitempty"`
	Diagnostics []Diagnostic   `json:"diagnostics,omitempty"`
	Edit        *WorkspaceEdit `json:"edit,omitempty"`
	Command     *Command       `json:"command,omitempty"`
}

type CodeActionParams struct {
    TextDocument TextDocumentIdentifier `json:"textDocument"`
    Range        Range                  `json:"range"`
    Context      CodeActionContext      `json:"context"`
}

type CodeActionContext struct {
    Diagnostics []Diagnostic `json:"diagnostics"`
    Only        []string     `json:"only,omitempty"`
}
