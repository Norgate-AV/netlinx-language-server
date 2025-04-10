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

func (s *Server) TextDocumentDidSave(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params lsp.DidSaveTextDocumentParams

	if err := json.Unmarshal(*req.Params, &params); err != nil {
		s.logger.Error("Failed to unmarshal document params", logrus.Fields{
			"method": "textDocument/didSave",
			"error":  err.Error(),
		})

		return
	}

	s.logger.LogDocumentEvent("save", params.TextDocument.URI)
}

func (s *Server) TextDocumentDiagnostic(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params lsp.PublishDiagnosticsParams

	if err := json.Unmarshal(*req.Params, &params); err != nil {
		s.logger.Error("Failed to unmarshal document params", logrus.Fields{
			"method": "textDocument/diagnostic",
			"error":  err.Error(),
		})

		return
	}

	s.logger.LogDocumentEvent("diagnostic", params.URI)
}

func (s *Server) TextDocumentSymbol(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params lsp.DocumentSymbolParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		s.logger.Error("Failed to unmarshal document params", logrus.Fields{
			"method": "textDocument/documentSymbol",
			"error":  err.Error(),
		})
		s.sendError(ctx, conn, req.ID, createError(jsonrpc2.CodeParseError, err.Error()))
		return
	}

	s.logger.LogDocumentEvent("symbol", params.TextDocument.URI)

	// symbols, err := s.state.ExtractSymbols(params.TextDocument.URI)
	// if err != nil {
	// 	s.logger.Error("Failed to extract symbols", logrus.Fields{
	// 		"uri":   params.TextDocument.URI,
	// 		"error": err.Error(),
	// 	})

	// 	// Return empty result on error
	// 	symbols = []lsp.DocumentSymbol{}
	// }

	symbols := createFakeSymbols()

	// Send response
	if err := conn.Reply(ctx, req.ID, symbols); err != nil {
		s.logger.Error("Failed to send symbol response", logrus.Fields{
			"error": err.Error(),
		})
	}
}

func createFakeSymbols() []lsp.DocumentSymbol {
	// Create a program symbol
	programSymbol := lsp.DocumentSymbol{
		Name: "PROGRAM: Test Program",
		Kind: lsp.SymbolKindFile,
		Range: lsp.Range{
			Start: lsp.Position{Line: 0, Character: 0},
			End:   lsp.Position{Line: 0, Character: 25},
		},
		SelectionRange: lsp.Range{
			Start: lsp.Position{Line: 0, Character: 0},
			End:   lsp.Position{Line: 0, Character: 25},
		},
	}

	// Create a DEFINE_DEVICE section
	deviceSection := lsp.DocumentSymbol{
		Name: "DEFINE_DEVICE",
		Kind: lsp.SymbolKindNamespace,
		Range: lsp.Range{
			Start: lsp.Position{Line: 2, Character: 0},
			End:   lsp.Position{Line: 4, Character: 0},
		},
		SelectionRange: lsp.Range{
			Start: lsp.Position{Line: 2, Character: 0},
			End:   lsp.Position{Line: 2, Character: 13},
		},
		Children: []lsp.DocumentSymbol{
			{
				Name: "dvTP",
				Kind: lsp.SymbolKindVariable,
				Range: lsp.Range{
					Start: lsp.Position{Line: 3, Character: 0},
					End:   lsp.Position{Line: 3, Character: 15},
				},
				SelectionRange: lsp.Range{
					Start: lsp.Position{Line: 3, Character: 0},
					End:   lsp.Position{Line: 3, Character: 4},
				},
			},
		},
	}

	// Create a function
	functionSymbol := lsp.DocumentSymbol{
		Name:   "DoSomething",
		Detail: "INTEGER",
		Kind:   lsp.SymbolKindFunction,
		Range: lsp.Range{
			Start: lsp.Position{Line: 10, Character: 0},
			End:   lsp.Position{Line: 13, Character: 1},
		},
		SelectionRange: lsp.Range{
			Start: lsp.Position{Line: 10, Character: 25},
			End:   lsp.Position{Line: 10, Character: 36},
		},
	}

	// Return the collection of symbols
	return []lsp.DocumentSymbol{
		programSymbol,
		deviceSection,
		functionSymbol,
	}
}
