package transport_test

import (
	"context"
	"os"
	"testing"

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
	defer func() { os.Stdin, os.Stdout = origStdin, origStdout }()

	r, w, _ := os.Pipe()
	defer r.Close()
	defer w.Close()

	os.Stdin, os.Stdout = r, w

	// Act
	transport, err := transport.NewStdioTransport()

	// Assert
	require.NoError(t, err)
	require.NotNil(t, transport)
}

func TestStdioTransport_StartReturnsChannel(t *testing.T) {
	// Arrange
	origStdin, origStdout := os.Stdin, os.Stdout
	defer func() { os.Stdin, os.Stdout = origStdin, origStdout }()

	r, w, _ := os.Pipe()
	defer r.Close()
	defer w.Close()

	os.Stdin, os.Stdout = r, w

	transport, _ := transport.NewStdioTransport()
	ctx := context.Background()

	// Simple handler
	handler := jsonrpc2.HandlerWithError(
		func(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
			return nil, nil
		})

	// Act
	disconnectCh, err := transport.Start(ctx, handler)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, disconnectCh)
}

func TestStdioTransport_CanBeStartedAndStopped(t *testing.T) {
	// Arrange
	origStdin, origStdout := os.Stdin, os.Stdout
	defer func() { os.Stdin, os.Stdout = origStdin, origStdout }()

	r, w, _ := os.Pipe()
	defer r.Close()
	defer w.Close()

	os.Stdin, os.Stdout = r, w

	transport, _ := transport.NewStdioTransport()
	ctx, cancel := context.WithCancel(context.Background())

	// Act
	_, err := transport.Start(ctx, jsonrpc2.HandlerWithError(
		func(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
			return nil, nil
		}))

	// Assert
	require.NoError(t, err)

	// Clean up
	cancel()
}
