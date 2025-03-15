package lsp

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/Norgate-AV/netlinx-language-server/internal/jsonrpc"
)

func (h *Handler) handleHover(ctx context.Context, conn *jsonrpc.Conn, req *jsonrpc.Request) {
	var params HoverParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		h.logger.Error("Failed to parse hover params", err)
		conn.ReplyWithError(ctx, req.ID, &jsonrpc.Error{
			Code:    jsonrpc.CodeParseError,
			Message: err.Error(),
		})
		return
	}

	doc, ok := h.server.docs.Get(params.TextDocument.URI)
	if !ok {
		h.logger.Warn("Document not found for hover", "uri", params.TextDocument.URI)
		conn.Reply(ctx, req.ID, nil)
		return
	}

	// Get word at position - simplified implementation
	lines := strings.Split(doc.Content, "\n")
	if params.Position.Line >= len(lines) {
		conn.Reply(ctx, req.ID, nil)
		return
	}

	line := lines[params.Position.Line]
	if params.Position.Character >= len(line) {
		conn.Reply(ctx, req.ID, nil)
		return
	}

	// This is where you would implement NetLinx-specific hover
	// For demonstration, just check for some known keywords
	wordStart := params.Position.Character
	for wordStart > 0 && isIdentifierChar(rune(line[wordStart-1])) {
		wordStart--
	}

	wordEnd := params.Position.Character
	for wordEnd < len(line) && isIdentifierChar(rune(line[wordEnd])) {
		wordEnd++
	}

	if wordStart == wordEnd {
		conn.Reply(ctx, req.ID, nil)
		return
	}

	word := line[wordStart:wordEnd]

	var hoverText string
	switch strings.ToUpper(word) {
	case "DEFINE_DEVICE":
		hoverText = "DEFINE_DEVICE - Used to define NetLinx devices"
	case "DEFINE_VARIABLE":
		hoverText = "DEFINE_VARIABLE - Used to define global variables in NetLinx"
	default:
		conn.Reply(ctx, req.ID, nil)
		return
	}

	conn.Reply(ctx, req.ID, Hover{
		Contents: []MarkedString{
			{
				Language: "netlinx",
				Value:    hoverText,
			},
		},
	})
}

func isIdentifierChar(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}
