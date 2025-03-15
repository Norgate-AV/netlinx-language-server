package jsonrpc

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
)

// Request represents a JSON-RPC request
type Request struct {
	ID      *json.RawMessage `json:"id,omitempty"`
	Method  string           `json:"method"`
	Params  *json.RawMessage `json:"params,omitempty"`
	JSONRPC string           `json:"jsonrpc"`
}

// Response represents a JSON-RPC response
type Response struct {
	ID      *json.RawMessage `json:"id,omitempty"`
	Result  interface{}      `json:"result,omitempty"`
	Error   *Error           `json:"error,omitempty"`
	JSONRPC string           `json:"jsonrpc"`
}

// Error represents a JSON-RPC error
type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorCode constants
const (
	CodeParseError     = -32700
	CodeInvalidRequest = -32600
	CodeMethodNotFound = -32601
	CodeInvalidParams  = -32602
	CodeInternalError  = -32603
)

// Handler handles JSON-RPC requests
type Handler interface {
	Handle(ctx context.Context, conn *Conn, req *Request)
}

// Conn represents a JSON-RPC connection
type Conn struct {
	in     *bufio.Reader
	out    io.Writer
	mu     sync.Mutex
	done   chan struct{}
	handle Handler
}

// NewConn creates a new connection
func NewConn(ctx context.Context, reader io.Reader, writer io.Writer, handler Handler) *Conn {
	conn := &Conn{
		in:     bufio.NewReader(reader),
		out:    writer,
		done:   make(chan struct{}),
		handle: handler,
	}

	go conn.listen(ctx)
	return conn
}

// listen reads messages from the connection
func (c *Conn) listen(ctx context.Context) {
	defer close(c.done)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			header, err := c.in.ReadString('\n')
			if err != nil {
				return
			}

			parts := strings.SplitN(header, ":", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Content-Length") {
				continue
			}

			contentLength, err := strconv.Atoi(strings.TrimSpace(parts[1]))
			if err != nil {
				continue
			}

			// Read the empty line
			if _, err := c.in.ReadString('\n'); err != nil {
				return
			}

			content := make([]byte, contentLength)
			if _, err := io.ReadFull(c.in, content); err != nil {
				return
			}

			var req Request
			if err := json.Unmarshal(content, &req); err != nil {
				c.sendParseError(err)
				continue
			}

			if req.JSONRPC != "2.0" {
				c.sendInvalidRequest("Invalid JSON-RPC version")
				continue
			}

			go c.handle.Handle(ctx, c, &req)
		}
	}
}

// IsNotify returns true if the request is a notification
func (r *Request) IsNotify() bool {
	return r.ID == nil
}

// Reply sends a response with a result
func (c *Conn) Reply(ctx context.Context, id *json.RawMessage, result interface{}) error {
	resp := Response{
		ID:      id,
		Result:  result,
		JSONRPC: "2.0",
	}
	return c.send(resp)
}

// ReplyWithError sends an error response
func (c *Conn) ReplyWithError(ctx context.Context, id *json.RawMessage, err *Error) error {
	resp := Response{
		ID:      id,
		Error:   err,
		JSONRPC: "2.0",
	}
	return c.send(resp)
}

// Close closes the connection
func (c *Conn) Close() error {
	return nil
}

// DisconnectNotify returns a channel that is closed when the connection is closed
func (c *Conn) DisconnectNotify() <-chan struct{} {
	return c.done
}

// send sends a JSON-RPC message
func (c *Conn) send(v interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(data), data)
	_, err = io.WriteString(c.out, msg)
	return err
}

// Helper methods for sending error responses
func (c *Conn) sendParseError(err error) {
	id := json.RawMessage(`null`)
	c.ReplyWithError(context.Background(), &id, &Error{
		Code:    CodeParseError,
		Message: fmt.Sprintf("Parse error: %v", err),
	})
}

func (c *Conn) sendInvalidRequest(msg string) {
	id := json.RawMessage(`null`)
	c.ReplyWithError(context.Background(), &id, &Error{
		Code:    CodeInvalidRequest,
		Message: msg,
	})
}
