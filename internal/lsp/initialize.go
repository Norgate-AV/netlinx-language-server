package lsp

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

type SemanticTokensClientCapabilities struct {
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
	Requests            struct {
		Range bool `json:"range"`
		Full  bool `json:"full"`
	} `json:"requests"`
	TokenTypes              []string `json:"tokenTypes"`
	TokenModifiers          []string `json:"tokenModifiers"`
	Formats                 []string `json:"formats"`
	OverlappingTokenSupport bool     `json:"overlappingTokenSupport,omitempty"`
	MultilineTokenSupport   bool     `json:"multilineTokenSupport,omitempty"`
}
