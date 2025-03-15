package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Norgate-AV/netlinx-language-server/internal/log"
	"github.com/Norgate-AV/netlinx-language-server/internal/lsp"
)

func main() {
	logger := log.NewLogger()
	defer logger.Sync()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cancel()
	}()

	server := lsp.NewServer(logger)
	if err := server.Run(ctx); err != nil {
		logger.Fatal("Server error", err)
	}
}
