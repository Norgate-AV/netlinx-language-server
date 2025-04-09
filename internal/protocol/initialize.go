package protocol

func NewInitializeResponse(_ int) InitializeResult {
	return InitializeResult{
		Capabilities: ServerCapabilities{
			// Use incremental sync for better performance
			TextDocumentSync: TextDocumentSyncKindIncremental,
			HoverProvider:    true,

			// Add document symbol support for testing structure parsing
			// DocumentSymbolProvider: true,

			// // For showing parse errors
			// DiagnosticProvider: &DiagnosticOptions{
			// 	InterFileDependencies: false,
			// 	WorkspaceDiagnostics:  false,
			// },
		},
		ServerInfo: ServerInfo{
			Name:    "netlinx-language-server",
			Version: "0.1.0",
		},
	}
}

const (
	TextDocumentSyncKindNone = iota
	TextDocumentSyncKindFull
	TextDocumentSyncKindIncremental
)

type DiagnosticOptions struct {
	InterFileDependencies bool `json:"interFileDependencies"`
	WorkspaceDiagnostics  bool `json:"workspaceDiagnostics"`
}
