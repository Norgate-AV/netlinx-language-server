package server

import (
	"context"
	"log"

	"github.com/Norgate-AV/netlinx-language-server/internal/analysis"
	"github.com/Norgate-AV/netlinx-language-server/internal/logger"

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

	return handler.registerHandlers()
}

func (s *Server) Stop() {
	s.logger.Println("Stopping server...")
	s.logger.Println("Server stopped")
}

func (s *Server) registerHandlers() *Server {
	s.handlers["initialize"] = s.Initialize
	s.handlers["initialized"] = s.Initialized
	s.handlers["shutdown"] = s.Shutdown
	s.handlers["exit"] = s.Exit
	s.handlers["textDocument/didOpen"] = s.TextDocumentDidOpen
	s.handlers["textDocument/didChange"] = s.TextDocumentDidChange
	s.handlers["textDocument/didClose"] = s.TextDocumentDidClose
	s.handlers["textDocument/hover"] = s.Hover

	return s
}

func (s *Server) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	s.logger.Printf("Received method: %s\n", req.Method)

	if handler, ok := s.handlers[req.Method]; ok {
		handler(ctx, conn, req)
		return
	}

	s.logger.Printf("Method not implemented: %s\n", req.Method)

	if req.ID == (jsonrpc2.ID{}) {
		return
	}

	if err := conn.Reply(ctx, req.ID, nil); err != nil {
		s.logger.Printf("Error sending reply: %v\n", err)
	}
}

func createError(code int64, message string) *jsonrpc2.Error {
	return &jsonrpc2.Error{
		Code:    code,
		Message: message,
	}
}

func sendError(ctx context.Context, conn *jsonrpc2.Conn, id jsonrpc2.ID, err *jsonrpc2.Error) {
	if replyErr := conn.ReplyWithError(ctx, id, err); replyErr != nil {
		log.Printf("Error sending error response: %v\n", replyErr)
	}
}
