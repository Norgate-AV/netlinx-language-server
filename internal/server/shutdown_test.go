package server_test

import (
	"context"
	"testing"

	"github.com/Norgate-AV/netlinx-language-server/internal/logger"
	"github.com/Norgate-AV/netlinx-language-server/internal/server"
	test "github.com/Norgate-AV/netlinx-language-server/internal/testing"
	"github.com/Norgate-AV/netlinx-language-server/internal/workspace"
	"github.com/Norgate-AV/netlinx-language-server/parser"

	"github.com/sourcegraph/jsonrpc2"
)

func Shutdown(s *server.Server, ctx context.Context, conn any, req *jsonrpc2.Request) {
	// Type assertion to check if conn implements the necessary method
	if replier, ok := conn.(interface {
		Reply(ctx context.Context, id jsonrpc2.ID, result any) error
	}); ok {
		s.Logger.LogServerEvent("Shutdown")

		if err := replier.Reply(ctx, req.ID, nil); err != nil {
			s.Logger.Error("Failed to send shutdown response", nil)
		}
	}
}

func TestShutdown(t *testing.T) {
	log := logger.NewStdLogger()
	ts, err := parser.NewTreeSitter()
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}

	defer ts.Close()

	state := workspace.NewState(&workspace.NewStateOptions{
		TreeSitter: ts,
		Logger:     log,
	})

	srv := server.NewServer(log, state)

	mockConn := &test.MockConn{}

	req := &jsonrpc2.Request{
		ID: jsonrpc2.ID{Num: 1},
	}

	Shutdown(srv, context.Background(), mockConn, req)

	// Verify response was sent with nil payload
	if !mockConn.ReplyCalled || mockConn.ReplyID != req.ID || mockConn.ReplyResult != nil {
		t.Errorf("Expected Reply(%v, nil), got Reply(%v, %v)", req.ID, mockConn.ReplyID, mockConn.ReplyResult)
	}
}
