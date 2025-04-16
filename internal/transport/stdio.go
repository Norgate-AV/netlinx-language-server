package transport

import (
	"context"
	"fmt"
	"os"

	"github.com/sourcegraph/jsonrpc2"
)

type StdioTransport struct{}

func NewStdioTransport() (*StdioTransport, error) {
	return &StdioTransport{}, nil
}

func (t *StdioTransport) Start(ctx context.Context, handler jsonrpc2.Handler) (<-chan struct{}, error) {
	stream := jsonrpc2.NewBufferedStream(&stdio{}, jsonrpc2.VSCodeObjectCodec{})
	conn := jsonrpc2.NewConn(ctx, stream, handler)

	fmt.Printf("Server listening on: %s...\n", os.Stdin.Name())

	return conn.DisconnectNotify(), nil
}

func (t *StdioTransport) Close() error {
	return nil
}

type stdio struct{}

func (stdio) Read(p []byte) (n int, err error) {
	return os.Stdin.Read(p)
}

func (stdio) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

func (stdio) Close() error {
	if err := os.Stdin.Close(); err != nil {
		return err
	}

	return os.Stdout.Close()
}
