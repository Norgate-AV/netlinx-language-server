package protocol

func NewInitializeResponse(_ int) InitializeResult {
	return InitializeResult{
		Capabilities: ServerCapabilities{
			TextDocumentSync: 1,
		},
		ServerInfo: ServerInfo{
			Name:    "netlinx-language-server",
			Version: "0.1.0",
		},
	}
}
