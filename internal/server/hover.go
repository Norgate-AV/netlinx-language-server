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

	_, ok := s.state.GetDocument(params.TextDocument.URI)
	if !ok {
		s.logger.Warn("Document not found for hover", logrus.Fields{
			"uri": params.TextDocument.URI,
		})

		conn.Reply(ctx, req.ID, nil)

		return
	}

	response, err := s.GetHoverInfo(params.TextDocument.URI, params.Position)
	if err != nil {
		s.sendError(ctx, conn, req.ID, createError(jsonrpc2.CodeInternalError, fmt.Sprintf("Hover error: %v", err)))
		return
	}

	if response == nil {
		s.logger.Warn("No hover information found", logrus.Fields{
			"uri":      params.TextDocument.URI,
			"position": params.Position,
		})

		conn.Reply(ctx, req.ID, nil)

		return
	}

	if err := conn.Reply(ctx, req.ID, response); err != nil {
		s.logger.Error("Failed to send hover response", logrus.Fields{
			"error": err.Error(),
		})
	}
}

func (s *Server) GetHoverInfo(uri string, position lsp.Position) (*lsp.Hover, error) {
	_, ok := s.state.GetDocument(uri)
	if !ok {
		return nil, nil // Document not found
	}

	// Here you'd implement actual hover logic based on document content and position
	// For now, return the simple message from the original implementation
	return &lsp.Hover{
		Contents: lsp.MarkupContent{
			Kind:  lsp.MarkupKindMarkdown,
			Value: "**NetLinx Language Server**\n\nConnection working correctly!",
		},
	}, nil
}
