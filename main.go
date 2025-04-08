package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Norgate-AV/netlinx-language-server/analysis"
	"github.com/Norgate-AV/netlinx-language-server/rpc"
	"github.com/sourcegraph/jsonrpc2"

	tree_sitter_netlinx "github.com/norgate-av/tree-sitter-netlinx/bindings/go"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

func main() {
	logger := getLogger("netlinx-language-server.log")
	logger.Println("Started Netlinx Language Server...")

	state := analysis.NewState()

	// Example tree-sitter parsing
	code := []byte("dvDisplay = 5001:1:0")
	parser := tree_sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_netlinx.Language()))

	tree := parser.Parse(code, nil)
	defer tree.Close()

	root := tree.RootNode()
	fmt.Println(root.ToSexp())

	// Create JSONRPC handler
	handler := rpc.NewLSPHandler(logger, &state)

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
