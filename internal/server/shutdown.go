package server

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/jsonrpc2"
)

func (s *Server) Shutdown(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	s.Logger.LogServerEvent("Shutdown")

	s.Logger.LogResponse(req.Method, req.ID)
	if err := conn.Reply(ctx, req.ID, nil); err != nil {
		s.Logger.Error("Failed to send shutdown response", logrus.Fields{
			"error": err.Error(),
		})
	}
}
