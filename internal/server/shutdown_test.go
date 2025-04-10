package server

import (
	"context"
	"testing"

	"github.com/Norgate-AV/netlinx-language-server/internal/analysis"
	"github.com/Norgate-AV/netlinx-language-server/internal/logger"
	test "github.com/Norgate-AV/netlinx-language-server/internal/testing"

	"github.com/sourcegraph/jsonrpc2"
)

func (s *Server) TestShutdown(ctx context.Context, conn interface{}, req *jsonrpc2.Request) {
	// Type assertion to check if conn implements the necessary method
	if replier, ok := conn.(interface {
		Reply(ctx context.Context, id jsonrpc2.ID, result interface{}) error
	}); ok {
		s.logger.LogServerEvent("Shutdown")

		if err := replier.Reply(ctx, req.ID, nil); err != nil {
			s.logger.Error("Failed to send shutdown response", nil)
		}
	}
}

func TestShutdown(t *testing.T) {
	log := logger.NewStdLogger()
	state := analysis.NewState()
	srv := NewServer(log, state)

	mockConn := &test.MockConn{}

	req := &jsonrpc2.Request{
		ID: jsonrpc2.ID{Num: 1},
	}

	srv.TestShutdown(context.Background(), mockConn, req)

	// Verify response was sent with nil payload
	if !mockConn.ReplyCalled || mockConn.ReplyID != req.ID || mockConn.ReplyResult != nil {
		t.Errorf("Expected Reply(%v, nil), got Reply(%v, %v)", req.ID, mockConn.ReplyID, mockConn.ReplyResult)
	}
}
