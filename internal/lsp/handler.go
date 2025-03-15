package lsp

import (
	"context"
	"encoding/json"

	"github.com/Norgate-AV/netlinx-language-server/internal/jsonrpc"
	"github.com/Norgate-AV/netlinx-language-server/internal/log"
)

// Handler implements the JSON-RPC handler for the language server
type Handler struct {
	server *Server
	logger *log.Logger
}

// Handle processes incoming JSON-RPC requests
func (h *Handler) Handle(ctx context.Context, conn *jsonrpc.Conn, req *jsonrpc.Request) {
	switch req.Method {
	case "initialize":
		h.handleInitialize(ctx, conn, req)
	case "initialized":
		// Nothing to do here
	case "shutdown":
		conn.Reply(ctx, req.ID, nil)
	case "exit":
		conn.Close()
	case "textDocument/didOpen":
		h.handleTextDocumentDidOpen(ctx, conn, req)
	case "textDocument/didChange":
		h.handleTextDocumentDidChange(ctx, conn, req)
	case "textDocument/didClose":
		h.handleTextDocumentDidClose(ctx, conn, req)
	case "textDocument/completion":
		h.handleCompletion(ctx, conn, req)
	case "textDocument/hover":
		h.handleHover(ctx, conn, req)
	default:
		if req.IsNotify() {
			h.logger.Debug("Ignoring notification", "method", req.Method)
			return
		}
		conn.ReplyWithError(ctx, req.ID, &jsonrpc.Error{
			Code:    jsonrpc.CodeMethodNotFound,
			Message: "method not supported: " + req.Method,
		})
	}
}

func (h *Handler) handleInitialize(ctx context.Context, conn *jsonrpc.Conn, req *jsonrpc.Request) {
	var params InitializeParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		h.logger.Error("Failed to unmarshal initialize params", err)
		conn.ReplyWithError(ctx, req.ID, &jsonrpc.Error{
			Code:    jsonrpc.CodeParseError,
			Message: err.Error(),
		})
		return
	}

	result := InitializeResult{
		Capabilities: ServerCapabilities{
			TextDocumentSync: TextDocumentSyncKindFull,
			CompletionProvider: &CompletionOptions{
				ResolveProvider:   false,
				TriggerCharacters: []string{".", "#"},
			},
			HoverProvider: true,
		},
	}

	h.logger.Info("Initialized language server", "clientName", params.ClientInfo.Name)
	conn.Reply(ctx, req.ID, result)
}

// Add handler implementations for textDocument notifications and requests here
