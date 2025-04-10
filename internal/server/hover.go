package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Norgate-AV/netlinx-language-server/internal/lsp"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/jsonrpc2"
)

func (s *Server) Hover(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params lsp.HoverParams

	if err := json.Unmarshal(*req.Params, &params); err != nil {
		s.logger.Error("Failed to unmarshal hover params", logrus.Fields{
			"error": err.Error(),
		})

		s.sendError(ctx, conn, req.ID, createError(jsonrpc2.CodeParseError, fmt.Sprintf("Invalid hover params: %v", err)))

		return
	}

	response := lsp.Hover{
		Contents: lsp.MarkupContent{
			Kind:  lsp.MarkupKindMarkdown,
			Value: "**NetLinx Language Server**\n\nConnection working correctly!",
		},
	}

	if err := conn.Reply(ctx, req.ID, response); err != nil {
		s.logger.Error("Failed to send hover response", logrus.Fields{
			"error": err.Error(),
		})
	}
}
