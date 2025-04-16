package transport

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"

	"github.com/Norgate-AV/netlinx-language-server/internal/logger"

	"github.com/containerd/fifo"
	"github.com/sourcegraph/jsonrpc2"
)

type PipeTransport struct {
	name    string
	logger  logger.Logger
	inPipe  io.ReadCloser
	outPipe io.WriteCloser
}

func NewPipeTransport(name string, logger logger.Logger) (*PipeTransport, error) {
	return &PipeTransport{
		name:   name,
		logger: logger,
	}, nil
}

func (t *PipeTransport) Start(ctx context.Context, handler jsonrpc2.Handler) (<-chan struct{}, error) {
	tempDir := os.TempDir()
	inPipePath := filepath.Join(tempDir, fmt.Sprintf("%s-in", t.name))
	outPipePath := filepath.Join(tempDir, fmt.Sprintf("%s-out", t.name))

	// Clean up existing pipes
	os.Remove(inPipePath)
	os.Remove(outPipePath)

	// Create output FIFO (to client)
	outFifo, err := fifo.OpenFifo(ctx, outPipePath, syscall.O_WRONLY|syscall.O_CREAT, 0o666)
	if err != nil {
		return nil, fmt.Errorf("failed to create output pipe: %w", err)
	}

	// Print connection info
	fmt.Printf("Server listening on pipes:\nIN: %s\nOUT: %s\n", inPipePath, outPipePath)

	// Create input FIFO (from client)
	inFifo, err := fifo.OpenFifo(ctx, inPipePath, syscall.O_RDONLY|syscall.O_CREAT, 0o666)
	if err != nil {
		outFifo.Close()
		return nil, fmt.Errorf("failed to create input pipe: %w", err)
	}

	t.inPipe = inFifo
	t.outPipe = outFifo

	// Create a stream from the pipes
	stream := jsonrpc2.NewBufferedStream(&pipeStream{inFifo, outFifo}, jsonrpc2.VSCodeObjectCodec{})
	conn := jsonrpc2.NewConn(ctx, stream, handler)

	return conn.DisconnectNotify(), nil
}

func (t *PipeTransport) Close() error {
	var err1, err2 error

	if t.inPipe != nil {
		err1 = t.inPipe.Close()
	}

	if t.outPipe != nil {
		err2 = t.outPipe.Close()
	}

	if err1 != nil {
		return err1
	}

	return err2
}

type pipeStream struct {
	in  io.ReadCloser
	out io.WriteCloser
}

func (p *pipeStream) Read(data []byte) (int, error) {
	return p.in.Read(data)
}

func (p *pipeStream) Write(data []byte) (int, error) {
	return p.out.Write(data)
}

func (p *pipeStream) Close() error {
	err1 := p.in.Close()
	err2 := p.out.Close()

	if err1 != nil {
		return err1
	}

	return err2
}
