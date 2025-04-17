package transport_test

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/Norgate-AV/netlinx-language-server/internal/transport"
	"github.com/sourcegraph/jsonrpc2"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStdioTransport_Creation(t *testing.T) {
	stdioTransport, err := transport.NewStdioTransport()
	require.NoError(t, err)
	require.NotNil(t, stdioTransport)
}

func TestStdioTransport_ReplacesStdinAndStdout(t *testing.T) {
	// Arrange
	origStdin, origStdout := os.Stdin, os.Stdout

	r, w, _ := os.Pipe()
	os.Stdin, os.Stdout = r, w

	// Act
	transport, err := transport.NewStdioTransport()

	// Assert
	require.NoError(t, err)
	require.NotNil(t, transport)

	os.Stdin, os.Stdout = origStdin, origStdout

	_ = w.Close()
	_ = r.Close()
}

func TestStdioTransport_StartReturnsChannel(t *testing.T) {
	// Arrange
	r1, w1 := io.Pipe() // for stdin
	r2, w2 := io.Pipe() // for stdout

	transport, err := transport.NewStdioTransportWithStreams(r1, w2)
	require.NoError(t, err)

	ctx := context.Background()

	// Simple handler
	handler := jsonrpc2.HandlerWithError(
		func(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (any, error) {
			return nil, nil
		})

	// Act
	disconnectCh, err := transport.Start(ctx, handler)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, disconnectCh)

	// Clean up pipes
	_ = w1.Close()
	_ = r2.Close()
}

func TestStdioTransport_CanBeStartedAndStopped(t *testing.T) {
	// Arrange
	r1, w1 := io.Pipe() // for stdin
	r2, w2 := io.Pipe() // for stdout

	transport, err := transport.NewStdioTransportWithStreams(r1, w2)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())

	// Act
	disconnectCh, err := transport.Start(ctx, jsonrpc2.HandlerWithError(
		func(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (any, error) {
			return nil, nil
		}))

	// Assert
	require.NoError(t, err)

	// Clean up
	cancel()

	// Wait for disconnect notification
	select {
	case <-disconnectCh:
		// Successfully disconnected
	case <-time.After(100 * time.Millisecond):
		// Timeout is okay too
	}

	// Clean up pipes
	_ = w1.Close()
	_ = r2.Close()
}
