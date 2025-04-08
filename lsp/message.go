package lsp

// We can simplify these message types since jsonrpc2 will handle most of the structure for us
// These are now mainly for documentation and type checking

type Request struct {
	// jsonrpc2 will handle the "jsonrpc", "id", and "method" fields
	// We only need to define the types for the parameters in specific requests
}

type Response struct {
	// jsonrpc2 will handle the "jsonrpc", "id", and "error" fields
	// We only need to define the types for the results in specific responses
}

type Notification struct {
	// jsonrpc2 will handle the "jsonrpc" and "method" fields
	// We only need to define the types for the parameters in specific notifications
}

// LSPError represents an LSP error object
type LSPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
