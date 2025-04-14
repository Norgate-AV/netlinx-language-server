package server

import (
	"context"

	"github.com/Norgate-AV/netlinx-language-server/internal/analysis"
	"github.com/Norgate-AV/netlinx-language-server/internal/logger"
	"github.com/Norgate-AV/netlinx-language-server/internal/lsp"

	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/jsonrpc2"
)

type handler func(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request)

type Server struct {
	Logger   logger.Logger
	state    *analysis.State
	handlers map[string]handler
}

func NewServer(logger logger.Logger, state *analysis.State) *Server {
	handler := &Server{
		Logger:   logger,
		state:    state,
		handlers: make(map[string]handler),
	}

	return handler.registerHandlers()
}

func (s *Server) Stop() {
	s.Logger.LogServerEvent("Stopping")
	s.Logger.LogServerEvent("Stopped")
}

func (s *Server) registerHandlers() *Server {
	s.handlers[lsp.MethodInitialize] = s.Initialize
	s.handlers[lsp.MethodInitialized] = s.Initialized
	s.handlers[lsp.MethodShutdown] = s.Shutdown
	s.handlers[lsp.MethodExit] = s.Exit

	s.handlers[lsp.MethodTextDocumentDidOpen] = s.TextDocumentDidOpen
	s.handlers[lsp.MethodTextDocumentDidChange] = s.TextDocumentDidChange
	s.handlers[lsp.MethodTextDocumentDidClose] = s.TextDocumentDidClose
	s.handlers[lsp.MethodTextDocumentDidSave] = s.TextDocumentDidSave

	s.handlers[lsp.MethodTextDocumentHover] = s.Hover
	s.handlers[lsp.MethodTextDocumentDocumentSymbol] = s.TextDocumentSymbol
	s.handlers[lsp.MethodTextDocumentDiagnostic] = s.TextDocumentDiagnostic
	// s.handlers[lsp.MethodWorkspaceDidChangeWatchedFiles] = s.WortkspaceDidChangeWatchedFiles

	s.handlers[lsp.MethodNetLinxServerLogPath] = s.NetLinxServerLogPath

	return s
}

func (s *Server) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	s.Logger.LogMessage(req.Method, req.ID)

	if handler, ok := s.handlers[req.Method]; ok {
		handler(ctx, conn, req)
		return
	}

	s.Logger.Warn("Method not implemented", logrus.Fields{
		"method": req.Method,
	})

	if req.ID == (jsonrpc2.ID{}) {
		return
	}

	s.Logger.LogResponse(req.Method, req.ID)

	if err := conn.Reply(ctx, req.ID, nil); err != nil {
		s.Logger.Error("Failed to send response", logrus.Fields{
			"error": err.Error(),
		})
	}
}

func createError(code int64, message string) *jsonrpc2.Error {
	return &jsonrpc2.Error{
		Code:    code,
		Message: message,
	}
}

func (s *Server) sendError(ctx context.Context, conn *jsonrpc2.Conn, id jsonrpc2.ID, err *jsonrpc2.Error) {
	if replyErr := conn.ReplyWithError(ctx, id, err); replyErr != nil {
		s.Logger.Error("Failed to send error response", logrus.Fields{
			"error": replyErr.Error(),
		})
	}
}

func (s *Server) NetLinxServerLogPath(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	s.Logger.LogResponse(req.Method, req.ID)

	if err := conn.Reply(ctx, req.ID, s.Logger.GetFilePath()); err != nil {
		s.Logger.Error("Failed to send log path response", logrus.Fields{
			"error": err.Error(),
		})
	}
}
