package server

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/jsonrpc2"
)

func (s *Server) Shutdown(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	s.logger.LogServerEvent("Shutdown")

	if err := conn.Reply(ctx, req.ID, nil); err != nil {
		s.logger.Error("Failed to send shutdown response", logrus.Fields{
			"error": err.Error(),
		})
	}
}
