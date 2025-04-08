package lsp

type DidOpenTextDocumentNotification struct {
	// Notification fields are handled by jsonrpc2
	Params DidOpenTextDocumentParams `json:"params"`
}

type DidOpenTextDocumentParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}
