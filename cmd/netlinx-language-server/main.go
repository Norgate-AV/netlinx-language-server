package main

import (
	"context"
	"log"
	"os"

	"github.com/Norgate-AV/netlinx-language-server/internal/analysis"
	"github.com/Norgate-AV/netlinx-language-server/internal/server"
	"github.com/sourcegraph/jsonrpc2"
)

func main() {
	logger := getLogger("netlinx-language-server.log")
	logger.Println("Started Netlinx Language Server...")

	state := analysis.NewState()

	// Create JSONRPC handler
	handler := server.NewLSPHandler(logger, state)

	// Create and run JSONRPC connection using standard input/output
	<-jsonrpc2.NewConn(
		context.Background(),
		jsonrpc2.NewBufferedStream(&stdinStdout{}, jsonrpc2.VSCodeObjectCodec{}),
		handler,
	).DisconnectNotify()
}

// stdinStdout implements io.ReadWriteCloser for stdin/stdout communication
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

func getLogger(fileName string) *log.Logger {
	logFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o666)
	if err != nil {
		panic(err)
	}

	return log.New(logFile, "[netlinx-language-server]", log.Ldate|log.Ltime|log.Lshortfile)
}
