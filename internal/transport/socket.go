package transport

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"sync"

	"github.com/Norgate-AV/netlinx-language-server/internal/logger"

	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/jsonrpc2"
)

type SocketTransport struct {
	port      string
	logger    logger.Logger
	listener  net.Listener
	conn      *jsonrpc2.Conn
	closeOnce sync.Once
}

func NewSocketTransport(port string, logger logger.Logger) (*SocketTransport, error) {
	_, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("invalid port number: %w", err)
	}

	return &SocketTransport{
		port:   port,
		logger: logger,
	}, nil
}

func (t *SocketTransport) Start(ctx context.Context, handler jsonrpc2.Handler) (<-chan struct{}, error) {
	port, err := strconv.Atoi(t.port)
	if err != nil {
		return nil, fmt.Errorf("invalid port number: %w", err)
	}

	addr := &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: port,
	}

	t.listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on port %d: %w", port, err)
	}

	t.logger.Info("Listening for TCP connections", logrus.Fields{
		"address": t.Address(),
		"port":    port,
	})

	// Print the address so clients know where to connect
	fmt.Printf("Server listening on: %s...\n", t.listener.Addr())

	// Accept a single connection
	doneCh := make(chan struct{})

	go func() {
		defer close(doneCh)

		conn, err := t.listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return // Exit silently if context canceled
			default:
				t.logger.Error("Failed to accept connection", logrus.Fields{
					"error": err.Error(),
				})
			}

			return
		}

		t.logger.Info("Accepted connection", logrus.Fields{
			"remote_addr": conn.RemoteAddr().String(),
		})

		// Create the JSON-RPC connection
		stream := jsonrpc2.NewBufferedStream(conn, jsonrpc2.VSCodeObjectCodec{})
		t.conn = jsonrpc2.NewConn(ctx, stream, handler)

		// Wait for the connection to close
		<-t.conn.DisconnectNotify()
	}()

	return doneCh, nil
}

func (t *SocketTransport) Close() error {
	var err error

	t.closeOnce.Do(func() {
		if t.listener != nil {
			err = t.listener.Close()
		}
	})

	return err
}

func (t *SocketTransport) Address() string {
	if t.listener == nil {
		return ""
	}

	return t.listener.Addr().String()
}
