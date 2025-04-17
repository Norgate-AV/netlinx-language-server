package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Norgate-AV/netlinx-language-server/internal/analysis/semantic"
	"github.com/Norgate-AV/netlinx-language-server/internal/lsp"

	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/jsonrpc2"
)

func (s *Server) Hover(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params lsp.HoverParams

	if err := json.Unmarshal(*req.Params, &params); err != nil {
		s.Logger.Error("Failed to unmarshal hover params", logrus.Fields{
			"error": err.Error(),
		})

		s.sendError(ctx, conn, req.ID, createError(jsonrpc2.CodeParseError, fmt.Sprintf("Invalid hover params: %v", err)))

		return
	}

	_, ok := s.state.GetDocument(params.TextDocument.URI)
	if !ok {
		s.Logger.Warn("Document not found for hover", logrus.Fields{
			"uri": params.TextDocument.URI,
		})

		s.Logger.LogResponse(req.Method, req.ID)

		if err := conn.Reply(ctx, req.ID, nil); err != nil {
			s.Logger.Error(fmt.Sprintf("Failed to send empty hover response: %v", err), nil)
		}

		return
	}

	response, err := s.GetHoverInfo(params.TextDocument.URI, params.Position)
	if err != nil {
		s.sendError(ctx, conn, req.ID, createError(jsonrpc2.CodeInternalError, fmt.Sprintf("Hover error: %v", err)))
		return
	}

	if response == nil {
		s.Logger.Warn("No hover information found", logrus.Fields{
			"uri":      params.TextDocument.URI,
			"position": params.Position,
		})

		s.Logger.LogResponse(req.Method, req.ID)

		if err := conn.Reply(ctx, req.ID, nil); err != nil {
			s.Logger.Error(fmt.Sprintf("Failed to send empty hover response: %v", err), nil)
		}

		return
	}

	s.Logger.LogResponse(req.Method, req.ID)
	if err := conn.Reply(ctx, req.ID, response); err != nil {
		s.Logger.Error("Failed to send hover response", logrus.Fields{
			"error": err.Error(),
		})
	}
}

func (s *Server) GetHoverInfo(uri string, position lsp.Position) (*lsp.Hover, error) {
	document, ok := s.state.GetDocument(uri)
	if !ok {
		return nil, nil // Document not found
	}

	// If document has semantic model
	if document.SemanticModel != nil {
		// Find symbol at position
		for _, symbol := range document.SemanticModel.Symbols {
			if containsPosition(symbol.Range, position) {
				// Create hover content based on symbol type
				content := createHoverForSymbol(symbol)
				return &lsp.Hover{
					Contents: lsp.MarkupContent{
						Kind:  lsp.MarkupKindMarkdown,
						Value: content,
					},
					Range: &symbol.Range,
				}, nil
			}
		}
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

func containsPosition(r lsp.Range, p lsp.Position) bool {
	// Check if position is within range
	return (p.Line > r.Start.Line || (p.Line == r.Start.Line && p.Character >= r.Start.Character)) &&
		(p.Line < r.End.Line || (p.Line == r.End.Line && p.Character <= r.End.Character))
}

func createHoverForSymbol(symbol *semantic.Symbol) string {
	switch symbol.Type {
	case semantic.DeviceSymbol:
		return fmt.Sprintf("**Device:** %s\n\n%s", symbol.Name, symbol.Value)
	case semantic.ConstantSymbol:
		return fmt.Sprintf("**Constant %s:** %s\n\n%s", symbol.DataType, symbol.Name, symbol.Value)
	case semantic.VariableSymbol:
		varType := "Variable"
		switch symbol.VariableKind {
		case semantic.VolatileVar:
			varType = "Volatile Variable"
		case semantic.NonVolatileVar:
			varType = "Non-volatile Variable"
		}
		return fmt.Sprintf("**%s %s:** %s", varType, symbol.DataType, symbol.Name)
	// Add other cases
	default:
		return fmt.Sprintf("**Symbol:** %s", symbol.Name)
	}
}
