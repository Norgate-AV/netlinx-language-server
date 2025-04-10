package server

import (
	"context"

	"github.com/sourcegraph/jsonrpc2"
)

func (s *Server) Exit(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	s.logger.Println("Exit notification received")
}
