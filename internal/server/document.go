package server

import (
	"context"
	"encoding/json"

	"github.com/Norgate-AV/netlinx-language-server/internal/lsp"
	"github.com/sirupsen/logrus"

	"github.com/sourcegraph/jsonrpc2"
)

func (s *Server) TextDocumentDidOpen(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params lsp.DidOpenTextDocumentParams

	if err := json.Unmarshal(*req.Params, &params); err != nil {
		s.logger.Error("Failed to unmarshal document params", logrus.Fields{
			"method": "textDocument/didOpen",
			"error":  err.Error(),
		})
		return
	}

	s.logger.LogDocumentEvent("open", params.TextDocument.URI)
	s.state.AddDocument(params.TextDocument.URI, params.TextDocument.Text)
}

func (s *Server) TextDocumentDidChange(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params lsp.DidChangeTextDocumentParams

	if err := json.Unmarshal(*req.Params, &params); err != nil {
		s.logger.Error("Failed to unmarshal document params", logrus.Fields{
			"method": "textDocument/didChange",
			"error":  err.Error(),
		})

		return
	}

	s.logger.LogDocumentEvent("change", params.TextDocument.URI)
	if len(params.ContentChanges) > 0 {
		s.state.UpdateDocument(params.TextDocument.URI, params.ContentChanges[0].Text)
	}
}

func (s *Server) TextDocumentDidClose(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params lsp.DidCloseTextDocumentParams

	if err := json.Unmarshal(*req.Params, &params); err != nil {
		s.logger.Error("Failed to unmarshal document params", logrus.Fields{
			"method": "textDocument/didClose",
			"error":  err.Error(),
		})

		return
	}

	s.logger.LogDocumentEvent("close", params.TextDocument.URI)
	s.state.CloseDocument(params.TextDocument.URI)
}
