package lsp

import (
	"context"
	"encoding/json"

	"github.com/Norgate-AV/netlinx-language-server/internal/jsonrpc"
)

// CompletionItem represents a completion suggestion
type CompletionItem struct {
	Label         string             `json:"label"`
	Kind          CompletionItemKind `json:"kind,omitempty"`
	Detail        string             `json:"detail,omitempty"`
	Documentation string             `json:"documentation,omitempty"`
	InsertText    string             `json:"insertText,omitempty"`
}

// CompletionList represents a list of completion items
type CompletionList struct {
	IsIncomplete bool             `json:"isIncomplete"`
	Items        []CompletionItem `json:"items"`
}

func (h *Handler) handleCompletion(ctx context.Context, conn *jsonrpc.Conn, req *jsonrpc.Request) {
	var params CompletionParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		h.logger.Error("Failed to parse completion params", err)
		conn.ReplyWithError(ctx, req.ID, &jsonrpc.Error{
			Code:    jsonrpc.CodeParseError,
			Message: err.Error(),
		})
		return
	}

	_, ok := h.server.docs.Get(params.TextDocument.URI)
	if !ok {
		h.logger.Warn("Document not found for completion", "uri", params.TextDocument.URI)
		conn.Reply(ctx, req.ID, nil)
		return
	}

	// This is where you would implement NetLinx-specific completion
	// For now, just return a basic example

	items := []CompletionItem{
		{
			Label:         "DEFINE_DEVICE",
			Detail:        "Define a device in NetLinx",
			Documentation: "Used to define a device in the NetLinx system",
			InsertText:    "DEFINE_DEVICE",
		},
		{
			Label:         "DEFINE_VARIABLE",
			Detail:        "Define variables in NetLinx",
			Documentation: "Used to define global variables",
			InsertText:    "DEFINE_VARIABLE",
		},
	}

	conn.Reply(ctx, req.ID, CompletionList{
		IsIncomplete: false,
		Items:        items,
	})
}
