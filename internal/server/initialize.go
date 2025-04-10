package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Norgate-AV/netlinx-language-server/internal/lsp"
	"github.com/sirupsen/logrus"

	"github.com/sourcegraph/jsonrpc2"
)

func (s *Server) Initialize(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params lsp.InitializeRequestParams

	if err := json.Unmarshal(*req.Params, &params); err != nil {
		s.logger.Error("Failed to unmarshal initialize params", logrus.Fields{
			"error": err.Error(),
		})

		s.sendError(ctx, conn, req.ID, createError(jsonrpc2.CodeParseError, fmt.Sprintf("Invalid initialize params: %v", err)))

		return
	}

	s.logger.Info("Client connected", logrus.Fields{
		"client_name":    params.ClientInfo.Name,
		"client_version": params.ClientInfo.Version,
	})

	response := NewInitializeResponse()
	if err := conn.Reply(ctx, req.ID, response); err != nil {
		s.logger.Error("Failed to send initialize response", logrus.Fields{
			"error": err.Error(),
		})
	}
}

func (s *Server) Initialized(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	s.logger.LogServerEvent("Initialized")
}

func NewInitializeResponse() lsp.InitializeResult {
	b := func(v bool) *bool { return &v }

	return lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			TextDocumentSync:       lsp.TextDocumentSyncKindIncremental,
			HoverProvider:          b(true),
			DocumentSymbolProvider: b(true),
			DiagnosticProvider: &lsp.DiagnosticOptions{
				InterFileDependencies: false,
				WorkspaceDiagnostics:  false,
			},
			// SemanticTokensProvider: &SemanticTokensOptions{
		},
		ServerInfo: lsp.ServerInfo{
			Name:    "netlinx-language-server",
			Version: "0.1.0",
		},
	}
}
