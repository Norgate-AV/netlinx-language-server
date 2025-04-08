// filepath: /Users/damienbutt/projects/netlinx-language-server/rpc/jsonrpc_handler.go
package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Norgate-AV/netlinx-language-server/analysis"
	"github.com/Norgate-AV/netlinx-language-server/lsp"
	"github.com/sourcegraph/jsonrpc2"
)

// LSPHandler implements the jsonrpc2.Handler interface for handling LSP messages.
type LSPHandler struct {
	logger *log.Logger
	state  *analysis.State
}

// NewLSPHandler creates a new LSPHandler instance.
func NewLSPHandler(logger *log.Logger, state *analysis.State) *LSPHandler {
	return &LSPHandler{
		logger: logger,
		state:  state,
	}
}

// Handle processes incoming JSON-RPC requests.
func (h *LSPHandler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	h.logger.Printf("Received method: %s\n", req.Method)

	switch req.Method {
	case "initialize":
		h.handleInitialize(ctx, conn, req)
	case "textDocument/didOpen":
		h.handleTextDocumentDidOpen(ctx, conn, req)
	default:
		h.logger.Printf("Method not implemented: %s\n", req.Method)
		if req.ID != (jsonrpc2.ID{}) { // Check if ID is not empty (not a notification)
			if err := conn.Reply(ctx, req.ID, nil); err != nil {
				h.logger.Printf("Error sending reply: %v\n", err)
			}
		}
	}
}

// handleInitialize handles the initialize request.
func (h *LSPHandler) handleInitialize(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params lsp.InitializeRequestParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		h.logger.Printf("Error unmarshalling initialize params: %v\n", err)
		sendError(ctx, conn, req.ID, createError(jsonrpc2.CodeParseError, fmt.Sprintf("Invalid initialize params: %v", err)))
		return
	}

	h.logger.Printf("Connected to: %s %s", params.ClientInfo.Name, params.ClientInfo.Version)

	response := lsp.NewInitializeResponse(0) // We'll ignore our own ID and use the one from the request
	if err := conn.Reply(ctx, req.ID, response); err != nil {
		h.logger.Printf("Error sending initialize response: %v\n", err)
	}
}

// handleTextDocumentDidOpen handles the textDocument/didOpen notification.
func (h *LSPHandler) handleTextDocumentDidOpen(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params lsp.DidOpenTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		h.logger.Printf("Error unmarshalling didOpen params: %v\n", err)
		return // No need to send a response for notifications
	}

	h.logger.Printf("Opened document: %s\n", params.TextDocument.URI)
	h.state.OpenDocument(params.TextDocument.URI, params.TextDocument.Text)
}

// createError creates a jsonrpc2.Error from a code and message
func createError(code int64, message string) *jsonrpc2.Error {
	return &jsonrpc2.Error{
		Code:    code,
		Message: message,
	}
}

// sendError sends an error response.
func sendError(ctx context.Context, conn *jsonrpc2.Conn, id jsonrpc2.ID, err *jsonrpc2.Error) {
	if replyErr := conn.ReplyWithError(ctx, id, err); replyErr != nil {
		log.Printf("Error sending error response: %v\n", replyErr)
	}
}
