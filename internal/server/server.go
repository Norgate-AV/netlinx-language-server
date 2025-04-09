package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Norgate-AV/netlinx-language-server/internal/analysis"
	"github.com/Norgate-AV/netlinx-language-server/internal/logger"
	"github.com/Norgate-AV/netlinx-language-server/internal/protocol"
	"github.com/sourcegraph/jsonrpc2"
)

type handler func(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request)

type Server struct {
	logger   logger.Logger
	state    *analysis.State
	handlers map[string]handler
}

func NewServer(logger logger.Logger, state *analysis.State) *Server {
	handler := &Server{
		logger:   logger,
		state:    state,
		handlers: make(map[string]handler),
	}

	handler.registerHandlers()
	return handler
}

func (s *Server) registerHandlers() {
	s.handlers["initialize"] = s.handleInitialize
	s.handlers["textDocument/didOpen"] = s.handleTextDocumentDidOpen
}

func (s *Server) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	s.logger.Printf("Received method: %s\n", req.Method)

	if handler, ok := s.handlers[req.Method]; ok {
		handler(ctx, conn, req)
	} else {
		s.logger.Printf("Method not implemented: %s\n", req.Method)
		if req.ID != (jsonrpc2.ID{}) {
			if err := conn.Reply(ctx, req.ID, nil); err != nil {
				s.logger.Printf("Error sending reply: %v\n", err)
			}
		}
	}
}

// handleInitialize handles the initialize request.
func (s *Server) handleInitialize(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params protocol.InitializeRequestParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		s.logger.Printf("Error unmarshalling initialize params: %v\n", err)
		sendError(ctx, conn, req.ID, createError(jsonrpc2.CodeParseError, fmt.Sprintf("Invalid initialize params: %v", err)))
		return
	}

	s.logger.Printf("Connected to: %s %s", params.ClientInfo.Name, params.ClientInfo.Version)

	response := protocol.NewInitializeResponse(0) // We'll ignore our own ID and use the one from the request
	if err := conn.Reply(ctx, req.ID, response); err != nil {
		s.logger.Printf("Error sending initialize response: %v\n", err)
	}
}

// handleTextDocumentDidOpen handles the textDocument/didOpen notification.
func (s *Server) handleTextDocumentDidOpen(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params protocol.DidOpenTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		s.logger.Printf("Error unmarshalling didOpen params: %v\n", err)
		return // No need to send a response for notifications
	}

	s.logger.Printf("Opened document: %s\n", params.TextDocument.URI)
	s.state.AddDocument(params.TextDocument.URI, params.TextDocument.Text)
}

// createError creates a jsonrpc2.Error from a code and message
func createError(code int64, message string) *jsonrpc2.Error {
	return &jsonrpc2.Error{
		Code:    code,
		Message: message,
	}
}

// sendError sends an error response.
func sendError(ctx context.Context, conn *jsonrpc2.Conn, id jsonrpc2.ID, err *jsonrpc2.Error) {
	if replyErr := conn.ReplyWithError(ctx, id, err); replyErr != nil {
		log.Printf("Error sending error response: %v\n", replyErr)
	}
}
