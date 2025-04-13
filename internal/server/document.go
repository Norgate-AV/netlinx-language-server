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
		s.Logger.Error("Failed to unmarshal document params", logrus.Fields{
			"method": lsp.MethodTextDocumentDidOpen,
			"error":  err.Error(),
		})
		return
	}

	s.Logger.LogDocumentEvent("open", params.TextDocument.URI)
	s.state.AddDocument(params.TextDocument.URI, params.TextDocument.Text)
}

func (s *Server) TextDocumentDidChange(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params lsp.DidChangeTextDocumentParams

	if err := json.Unmarshal(*req.Params, &params); err != nil {
		s.Logger.Error("Failed to unmarshal document params", logrus.Fields{
			"method": lsp.MethodTextDocumentDidChange,
			"error":  err.Error(),
		})

		return
	}

	s.Logger.LogDocumentEvent("change", params.TextDocument.URI)
	if len(params.ContentChanges) > 0 {
		s.state.UpdateDocument(params.TextDocument.URI, params.ContentChanges[0].Text)
	}
}

func (s *Server) TextDocumentDidClose(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params lsp.DidCloseTextDocumentParams

	if err := json.Unmarshal(*req.Params, &params); err != nil {
		s.Logger.Error("Failed to unmarshal document params", logrus.Fields{
			"method": lsp.MethodTextDocumentDidClose,
			"error":  err.Error(),
		})

		return
	}

	s.Logger.LogDocumentEvent("close", params.TextDocument.URI)
	s.state.CloseDocument(params.TextDocument.URI)
}

func (s *Server) TextDocumentDidSave(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params lsp.DidSaveTextDocumentParams

	if err := json.Unmarshal(*req.Params, &params); err != nil {
		s.Logger.Error("Failed to unmarshal document params", logrus.Fields{
			"method": "textDocument/didSave",
			"error":  err.Error(),
		})

		return
	}

	s.Logger.LogDocumentEvent("save", params.TextDocument.URI)
}

func (s *Server) TextDocumentDiagnostic(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params lsp.PublishDiagnosticsParams

	if err := json.Unmarshal(*req.Params, &params); err != nil {
		s.Logger.Error("Failed to unmarshal document params", logrus.Fields{
			"method": lsp.MethodTextDocumentDiagnostic,
			"error":  err.Error(),
		})

		return
	}

	s.Logger.LogDocumentEvent("diagnostic", params.URI)
}

func (s *Server) TextDocumentSymbol(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params lsp.DocumentSymbolParams

	if err := json.Unmarshal(*req.Params, &params); err != nil {
		s.Logger.Error("Failed to unmarshal document params", logrus.Fields{
			"method": lsp.MethodTextDocumentDocumentSymbol,
			"error":  err.Error(),
		})

		s.sendError(ctx, conn, req.ID, createError(jsonrpc2.CodeParseError, err.Error()))

		return
	}

	s.Logger.LogDocumentEvent("symbol", params.TextDocument.URI)

	symbols, err := s.state.ExtractSymbols(params.TextDocument.URI)
	if err != nil {
		s.Logger.Error("Failed to extract symbols", logrus.Fields{
			"uri":   params.TextDocument.URI,
			"error": err.Error(),
		})

		// Return empty result on error
		symbols = []lsp.DocumentSymbol{}
	}

	// Send response
	s.Logger.LogResponse(req.Method, req.ID)
	if err := conn.Reply(ctx, req.ID, symbols); err != nil {
		s.Logger.Error("Failed to send symbol response", logrus.Fields{
			"error": err.Error(),
		})
	}
}
