package transport

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Norgate-AV/netlinx-language-server/internal/logger"

	"github.com/sourcegraph/jsonrpc2"
)

type PipeTransport struct {
	name     string
	logger   logger.Logger
	pipe     io.ReadWriteCloser
	listener net.Listener
}

func NewPipeTransport(name string, logger logger.Logger) (*PipeTransport, error) {
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	return &PipeTransport{
		name:   name,
		logger: logger,
	}, nil
}

func (t *PipeTransport) Start(ctx context.Context, handler jsonrpc2.Handler) (<-chan struct{}, error) {
	// Platform-agnostic implementation that meets the LSP spec requirements
	var listener net.Listener
	var err error

	if runtime.GOOS == "windows" {
		// Windows: Emulate named pipe through TCP with pipe name file
		listener, err = net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return nil, fmt.Errorf("failed to create TCP listener: %w", err)
		}

		// Write pipe name and port to a special file for discovery
		pipeInfo := fmt.Sprintf("%d", listener.Addr().(*net.TCPAddr).Port)
		pipePath := filepath.Join(os.TempDir(), fmt.Sprintf("%s-pipe", t.name))

		if err := os.WriteFile(pipePath, []byte(pipeInfo), 0o666); err != nil {
			listener.Close()
			return nil, fmt.Errorf("failed to create pipe info file: %w", err)
		}

		t.logger.Info(fmt.Sprintf("Windows pipe transport listening on port %s (pipe info at %s)",
			pipeInfo, pipePath), nil)
	} else {
		// Unix: Use Unix domain socket (as per LSP spec)
		socketPath := filepath.Join(os.TempDir(), t.name)

		// Remove existing socket if present
		if err := os.Remove(socketPath); err != nil && !os.IsNotExist(err) {
			t.logger.Error(fmt.Sprintf("Failed to remove existing socket: %v", err), nil)
		}

		listener, err = net.Listen("unix", socketPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create unix domain socket: %w", err)
		}

		t.logger.Info(fmt.Sprintf("Unix socket transport listening on %s", socketPath), nil)
	}

	t.listener = listener

	// Create a disconnect channel
	disconnectCh := make(chan struct{})

	// Accept connections in a goroutine
	go func() {
		defer close(disconnectCh)

		conn, err := listener.Accept()
		if err != nil {
			t.logger.Error(fmt.Sprintf("Failed to accept connection: %v", err), nil)
			return
		}

		t.pipe = conn

		// Create a stream from the connection
		stream := jsonrpc2.NewBufferedStream(conn, jsonrpc2.VSCodeObjectCodec{})
		jsonConn := jsonrpc2.NewConn(ctx, stream, handler)

		select {
		case <-ctx.Done():
			jsonConn.Close()
		case <-jsonConn.DisconnectNotify():
			// Client disconnected
		}
	}()

	return disconnectCh, nil
}

func (t *PipeTransport) Close() error {
	var err error

	if t.pipe != nil {
		err = t.pipe.Close()
	}

	if t.listener != nil {
		// Only capture subsequent errors if we haven't found one yet
		if closeErr := t.listener.Close(); err == nil {
			err = closeErr
		}
	}

	// Clean up pipe info file on Windows
	if runtime.GOOS == "windows" {
		pipePath := filepath.Join(os.TempDir(), fmt.Sprintf("%s-pipe", t.name))
		// Only capture this error if we haven't found one yet
		if removeErr := os.Remove(pipePath); removeErr != nil && !os.IsNotExist(removeErr) && err == nil {
			err = removeErr
		}
	}

	return err
}
