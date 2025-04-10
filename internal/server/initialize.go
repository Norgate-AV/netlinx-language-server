package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Norgate-AV/netlinx-language-server/internal/lsp"

	"github.com/sourcegraph/jsonrpc2"
)

func (s *Server) Initialize(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params lsp.InitializeRequestParams

	if err := json.Unmarshal(*req.Params, &params); err != nil {
		s.logger.Printf("Error unmarshalling initialize params: %v\n", err)
		sendError(ctx, conn, req.ID, createError(jsonrpc2.CodeParseError, fmt.Sprintf("Invalid initialize params: %v", err)))

		return
	}

	s.logger.Printf("Connected to: %s %s", params.ClientInfo.Name, params.ClientInfo.Version)

	response := NewInitializeResponse()
	if err := conn.Reply(ctx, req.ID, response); err != nil {
		s.logger.Printf("Error sending initialize response: %v\n", err)
	}
}

func (s *Server) Initialized(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	s.logger.Println("Server initialized")
}

func NewInitializeResponse() lsp.InitializeResult {
	return lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			TextDocumentSync:       lsp.TextDocumentSyncKindIncremental,
			HoverProvider:          true,
			DocumentSymbolProvider: true,
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
