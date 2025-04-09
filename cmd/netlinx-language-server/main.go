package main

import (
	"context"
	"os"

	"github.com/Norgate-AV/netlinx-language-server/internal/analysis"
	"github.com/Norgate-AV/netlinx-language-server/internal/logger"
	"github.com/Norgate-AV/netlinx-language-server/internal/server"
	"github.com/sourcegraph/jsonrpc2"
)

func main() {
	// logger := getLogger("netlinx-language-server.log")
	log, err := logger.NewFileLogger("netlinx-language-server.log")
	if err != nil {
		log = logger.NewStdLogger()
		log.Printf("Failed to initialize file logger: %v, falling back to stderr", err)
	}

	log.Println("Started Netlinx Language Server...")

	state := analysis.NewState()

	// Create JSONRPC handler
	server := server.NewServer(log, state)

	// Create and run JSONRPC connection using standard input/output
	<-jsonrpc2.NewConn(
		context.Background(),
		jsonrpc2.NewBufferedStream(&stdinStdout{}, jsonrpc2.VSCodeObjectCodec{}),
		server,
	).DisconnectNotify()
}

type stdinStdout struct{}

func (stdinStdout) Read(p []byte) (n int, err error) {
	return os.Stdin.Read(p)
}

func (stdinStdout) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

func (stdinStdout) Close() error {
	if err := os.Stdin.Close(); err != nil {
		return err
	}
	return os.Stdout.Close()
}
