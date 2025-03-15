package lsp

import (
	"context"
	"os"

	"github.com/Norgate-AV/netlinx-language-server/internal/jsonrpc"
	"github.com/Norgate-AV/netlinx-language-server/internal/log"
)

// Server represents the NetLinx language server instance
type Server struct {
	conn   *jsonrpc.Conn
	logger *log.Logger
	docs   *DocumentManager
}

// NewServer creates a new language server instance
func NewServer(logger *log.Logger) *Server {
	return &Server{
		logger: logger,
		docs:   NewDocumentManager(),
	}
}

// Run starts the language server
func (s *Server) Run(ctx context.Context) error {
	s.logger.Info("Starting NetLinx Language Server")

	handler := &Handler{
		server: s,
		logger: s.logger,
	}

	s.conn = jsonrpc.NewConn(ctx, os.Stdin, os.Stdout, handler)

	<-s.conn.DisconnectNotify()
	s.logger.Info("Language server stopped")
	return nil
}
