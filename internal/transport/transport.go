package transport

import (
	"context"

	"github.com/Norgate-AV/netlinx-language-server/internal/config"
	"github.com/Norgate-AV/netlinx-language-server/internal/logger"

	"github.com/sourcegraph/jsonrpc2"
)

type Transport interface {
	Start(ctx context.Context, handler jsonrpc2.Handler) (<-chan struct{}, error)
	Close() error
}

type Options struct {
	PipeName   string
	SocketPort string
	Logger     logger.Logger
}

func CreateTransport(transport string, options *Options) (Transport, error) {
	switch transport {
	case config.TransportTypeStdio:
		return NewStdioTransport()
	case config.TransportTypePipe:
		return NewPipeTransport(options.PipeName, options.Logger)
	case config.TransportTypeSocket:
		return NewSocketTransport(options.SocketPort, options.Logger)
	default:
		options.Logger.Info("Unknown transport requested, defaulting to stdio", nil)
		return NewStdioTransport()
	}
}
