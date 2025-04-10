package protocol

type InitializeRequest struct {
	Params InitializeRequestParams `json:"params"`
}

type InitializeRequestParams struct {
	ClientInfo *ClientInfo `json:"clientInfo"`
	// ... there's tons more that goes here
}

type InitializeResponse struct {
	Result InitializeResult `json:"result"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverInfo"`
}

func NewInitializeResponse(_ int) InitializeResult {
	return InitializeResult{
		Capabilities: ServerCapabilities{
			TextDocumentSync:       TextDocumentSyncKindIncremental,
			HoverProvider:          true,
			DocumentSymbolProvider: true,
			DiagnosticProvider: &DiagnosticOptions{
				InterFileDependencies: false,
				WorkspaceDiagnostics:  false,
			},
		},
		ServerInfo: ServerInfo{
			Name:    "netlinx-language-server",
			Version: "0.1.0",
		},
	}
}
