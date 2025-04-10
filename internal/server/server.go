package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Norgate-AV/netlinx-language-server/internal/analysis"
	"github.com/Norgate-AV/netlinx-language-server/internal/logger"
	"github.com/Norgate-AV/netlinx-language-server/internal/lsp"

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

func (s *Server) Stop() {}

func (s *Server) registerHandlers() {
	s.handlers["initialize"] = s.handleInitialize
	s.handlers["initialized"] = s.handleInitialized
	s.handlers["shutdown"] = s.handleShutdown
	s.handlers["exit"] = s.handleExit
	s.handlers["textDocument/didOpen"] = s.handleTextDocumentDidOpen
	s.handlers["textDocument/didChange"] = s.handleTextDocumentDidChange
	s.handlers["textDocument/didClose"] = s.handleTextDocumentDidClose
	s.handlers["textDocument/hover"] = s.handleHover
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

func (s *Server) handleInitialize(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params lsp.InitializeRequestParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		s.logger.Printf("Error unmarshalling initialize params: %v\n", err)
		sendError(ctx, conn, req.ID, createError(jsonrpc2.CodeParseError, fmt.Sprintf("Invalid initialize params: %v", err)))
		return
	}

	s.logger.Printf("Connected to: %s %s", params.ClientInfo.Name, params.ClientInfo.Version)

	response := lsp.NewInitializeResponse(0)
	if err := conn.Reply(ctx, req.ID, response); err != nil {
		s.logger.Printf("Error sending initialize response: %v\n", err)
	}
}

// Add these implementations
func (s *Server) handleInitialized(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	s.logger.Println("Server initialized")
	// No response needed for a notification
}

func (s *Server) handleShutdown(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	s.logger.Println("Shutdown request received")
	if err := conn.Reply(ctx, req.ID, nil); err != nil {
		s.logger.Printf("Error sending shutdown response: %v\n", err)
	}
}

func (s *Server) handleExit(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	s.logger.Println("Exit notification received")
	// Close connection - typically this would terminate the process
}

func (s *Server) handleTextDocumentDidOpen(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params lsp.DidOpenTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		s.logger.Printf("Error unmarshalling didOpen params: %v\n", err)
		return
	}

	s.logger.Printf("Opened document: %s\n", params.TextDocument.URI)
	s.state.AddDocument(params.TextDocument.URI, params.TextDocument.Text)
}

// Add this simple change handler
func (s *Server) handleTextDocumentDidChange(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
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

// Add this simple close handler
func (s *Server) handleTextDocumentDidClose(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params lsp.DidCloseTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		s.logger.Printf("Error unmarshalling didClose params: %v\n", err)
		return
	}

	s.logger.Printf("Closed document: %s\n", params.TextDocument.URI)
	s.state.CloseDocument(params.TextDocument.URI)
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
