package server

import (
	"context"

	"github.com/sourcegraph/jsonrpc2"
)

func (s *Server) Shutdown(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	s.logger.Println("Shutdown request received")

	if err := conn.Reply(ctx, req.ID, nil); err != nil {
		s.logger.Printf("Error sending shutdown response: %v\n", err)
	}
}
