package transport_test

import (
	"context"
	"testing"

	test "github.com/Norgate-AV/netlinx-language-server/internal/testing"
	"github.com/Norgate-AV/netlinx-language-server/internal/transport"
	"github.com/sourcegraph/jsonrpc2"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSocketTransport_Creation(t *testing.T) {
	testLogger := test.NewTestLogger(t)
	port := "0" // random port

	transport, err := transport.NewSocketTransport(port, testLogger)
	require.NoError(t, err)
	require.NotNil(t, transport)

	// Check that no listener is active yet
	assert.Empty(t, transport.Address(), "Address should be empty before Start()")

	_ = transport.Close()
}

func TestSocketTransport_SpecificPort(t *testing.T) {
	// Arrange
	testLogger := test.NewTestLogger(t)
	port := "9000" // specific port

	// Act
	transport, err := transport.NewSocketTransport(port, testLogger)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, transport)

	// Cleanup
	_ = transport.Close()
}

func TestSocketTransport_StartReturnsChannel(t *testing.T) {
	// Arrange
	testLogger := test.NewTestLogger(t)
	transport, _ := transport.NewSocketTransport("0", testLogger)

	ctx := context.Background()

	handler := jsonrpc2.HandlerWithError(
		func(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (any, error) {
			return nil, nil
		})

	// Act
	disconnectCh, err := transport.Start(ctx, handler)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, disconnectCh)

	// Cleanup - cancel context to stop listener
	_, cancel := context.WithCancel(context.Background())
	cancel()

	_ = transport.Close()
}

func TestSocketTransport_AddressAfterStart(t *testing.T) {
	// Arrange
	testLogger := test.NewTestLogger(t)
	transport, _ := transport.NewSocketTransport("0", testLogger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handler := jsonrpc2.HandlerWithError(
		func(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (any, error) {
			return nil, nil
		})

	// Act
	_, err := transport.Start(ctx, handler)
	require.NoError(t, err)

	// Assert
	address := transport.Address()
	assert.NotEmpty(t, address, "Address should not be empty after Start()")

	_ = transport.Close()
}

func TestSocketTransport_Close(t *testing.T) {
	// Arrange
	testLogger := test.NewTestLogger(t)
	transport, _ := transport.NewSocketTransport("0", testLogger)

	// Act
	err := transport.Close()

	// Assert
	assert.NoError(t, err, "Close should not return an error")
}

func TestSocketTransport_InvalidPort(t *testing.T) {
	// Arrange
	testLogger := test.NewTestLogger(t)

	// Act
	transport, err := transport.NewSocketTransport("invalid-port", testLogger)

	// Assert
	if assert.Error(t, err, "Should return error with invalid port") {
		assert.Contains(t, err.Error(), "port")
	}

	assert.Nil(t, transport, "Transport should be nil with invalid port")
}
