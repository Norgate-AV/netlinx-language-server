package testing

import (
	"context"

	"github.com/sourcegraph/jsonrpc2"
)

// MockConn implements the minimal jsonrpc2.Conn interface needed for testing
type MockConn struct {
	// Tracking for Reply calls
	ReplyCalled bool
	ReplyID     jsonrpc2.ID
	ReplyResult interface{}
	ReplyError  error
}

// Reply implements the jsonrpc2.Conn Reply method
func (c *MockConn) Reply(ctx context.Context, id jsonrpc2.ID, result interface{}) error {
	c.ReplyCalled = true
	c.ReplyID = id
	c.ReplyResult = result
	return c.ReplyError
}
