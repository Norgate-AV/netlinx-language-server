package server

import (
	"context"
	"encoding/json"

	"github.com/Norgate-AV/netlinx-language-server/internal/lsp"

	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/jsonrpc2"
)

func (s *Server) WorkspaceDidChangeWatchedFiles(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params lsp.DidChangeWatchedFilesParams

	if err := json.Unmarshal(*req.Params, &params); err != nil {
		s.logger.Error("Failed to unmarshal document params", logrus.Fields{
			"method": "workspace/didChangeWatchedFiles",
			"error":  err.Error(),
		})

		return
	}

	for _, change := range params.Changes {
		// Convert FileChangeType to action string
		var action string
		switch change.Type {
		case lsp.Created:
			action = "created"
		case lsp.Changed:
			action = "changed"
		case lsp.Deleted:
			action = "deleted"
		default:
			action = "unknown"
		}

		s.logger.LogDocumentEvent(action, change.URI)
	}
}
