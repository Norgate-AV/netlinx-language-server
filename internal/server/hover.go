package server

import (
	"context"
	"encoding/json"

	"github.com/Norgate-AV/netlinx-language-server/internal/lsp"
	"github.com/sourcegraph/jsonrpc2"
)

func (s *Server) handleHover(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params lsp.HoverParams

	if err := json.Unmarshal(*req.Params, &params); err != nil {
		s.logger.Printf("Error unmarshalling hover params: %v\n", err)
		sendError(ctx, conn, req.ID, createError(jsonrpc2.CodeParseError, err.Error()))

		return
	}

	response := lsp.Hover{
		Contents: lsp.MarkupContent{
			Kind:  lsp.MarkupKindMarkdown,
			Value: "**NetLinx Language Server**\n\nConnection working correctly!",
		},
	}

	if err := conn.Reply(ctx, req.ID, response); err != nil {
		s.logger.Printf("Error sending hover response: %v\n", err)
	}
}
