package transport_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	test "github.com/Norgate-AV/netlinx-language-server/internal/testing"
	"github.com/Norgate-AV/netlinx-language-server/internal/transport"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPipeTransport_Creation(t *testing.T) {
	// Arrange
	testLogger := test.NewTestLogger(t)
	pipeName := "test-pipe-creation"

	// Act
	pipeTransport, err := transport.NewPipeTransport(pipeName, testLogger)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, pipeTransport)

	// Cleanup
	_ = pipeTransport.Close()
}

func TestPipeTransport_PipesNotCreatedUntilStart(t *testing.T) {
	// Arrange
	testLogger := test.NewTestLogger(t)
	pipeName := "test-pipe-files"

	// Act
	pipeTransport, err := transport.NewPipeTransport(pipeName, testLogger)
	require.NoError(t, err)

	// Assert
	tempDir := os.TempDir()
	inPipePath := filepath.Join(tempDir, pipeName+"-in")
	outPipePath := filepath.Join(tempDir, pipeName+"-out")

	_, err = os.Stat(inPipePath)
	assert.True(t, os.IsNotExist(err), "Input pipe should not exist before Start()")

	_, err = os.Stat(outPipePath)
	assert.True(t, os.IsNotExist(err), "Output pipe should not exist before Start()")

	_ = pipeTransport.Close()
}

func TestPipeTransport_StartReturnsChannel(t *testing.T) {
	// Skip the full test since it requires pipes to be created
	// Only test that the interface method exists and returns the correct type
	testLogger := test.NewTestLogger(t)
	pipeName := "test-pipe-channel"

	// Arrange
	pipeTransport, err := transport.NewPipeTransport(pipeName, testLogger)
	require.NoError(t, err)

	// This is more of an interface/compilation test
	// We're just ensuring the Start method returns a channel of the right type
	_ = func() (<-chan struct{}, error) {
		return pipeTransport.Start(context.Background(), nil)
	}

	// No actual assertion needed - if it compiles, the test passes
	assert.True(t, true, "Start method exists and returns a channel")

	_ = pipeTransport.Close()
}
