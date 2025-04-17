package transport

import (
	"context"
	"io"
	"os"

	"github.com/sourcegraph/jsonrpc2"
)

type StdioTransport struct {
	reader io.Reader
	writer io.Writer
}

func NewStdioTransport() (*StdioTransport, error) {
	return NewStdioTransportWithStreams(os.Stdin, os.Stdout)
}

func NewStdioTransportWithStreams(reader io.Reader, writer io.Writer) (*StdioTransport, error) {
	return &StdioTransport{
		reader: reader,
		writer: writer,
	}, nil
}

func (t *StdioTransport) Start(ctx context.Context, handler jsonrpc2.Handler) (<-chan struct{}, error) {
	stream := jsonrpc2.NewBufferedStream(&stdio{
		reader: t.reader,
		writer: t.writer,
	}, jsonrpc2.VSCodeObjectCodec{})

	conn := jsonrpc2.NewConn(ctx, stream, handler)

	return conn.DisconnectNotify(), nil
}

func (t *StdioTransport) Close() error {
	return nil
}

type stdio struct {
	reader io.Reader
	writer io.Writer
}

func (s *stdio) Read(p []byte) (n int, err error) {
	return s.reader.Read(p)
}

func (s *stdio) Write(p []byte) (n int, err error) {
	return s.writer.Write(p)
}

func (s *stdio) Close() error {
	if closer, ok := s.reader.(io.Closer); ok {
		if err := closer.Close(); err != nil {
			return err
		}
	}

	if closer, ok := s.writer.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}
