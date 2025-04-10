package server

import (
	"context"
	"encoding/json"

	"github.com/Norgate-AV/netlinx-language-server/internal/lsp"

	"github.com/sourcegraph/jsonrpc2"
)

func (s *Server) TextDocumentDidOpen(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params lsp.DidOpenTextDocumentParams

	if err := json.Unmarshal(*req.Params, &params); err != nil {
		s.logger.Printf("Error unmarshalling didOpen params: %v\n", err)
		return
	}

	s.logger.Printf("Opened document: %s\n", params.TextDocument.URI)
	s.state.AddDocument(params.TextDocument.URI, params.TextDocument.Text)
}

func (s *Server) TextDocumentDidChange(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params lsp.DidChangeTextDocumentParams

	if err := json.Unmarshal(*req.Params, &params); err != nil {
		s.logger.Printf("Error unmarshalling didChange params: %v\n", err)
		return
	}

	s.logger.Printf("Changed document: %s\n", params.TextDocument.URI)
	if len(params.ContentChanges) > 0 {
		s.state.UpdateDocument(params.TextDocument.URI, params.ContentChanges[0].Text)
	}
}

func (s *Server) TextDocumentDidClose(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params lsp.DidCloseTextDocumentParams

	if err := json.Unmarshal(*req.Params, &params); err != nil {
		s.logger.Printf("Error unmarshalling didClose params: %v\n", err)
		return
	}

	s.logger.Printf("Closed document: %s\n", params.TextDocument.URI)
	s.state.CloseDocument(params.TextDocument.URI)
}
