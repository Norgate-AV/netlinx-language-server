package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

func EncodeMessage(msg any) string {
	content, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf(
		"Content-Length: %d\r\n\r\n%s",
		len(content),
		content,
	)
}

type BaseMessage struct {
	Method string `json:"method"`
}

func DecodeMessage(msg []byte) (string, int, error) {
	// Split the message into header and content
	header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return "", 0, errors.New("did not find the separator")
	}

	// Parse the header
	// Content-Length: <number>
	length, err := strconv.Atoi(string(header[len("Content-Length: "):]))
	if err != nil {
		return "", 0, err
	}

	// Check if the content length is correct
	if length != len(content) {
		return "", 0, errors.New("Content-Length does not match the content")
	}

	var message BaseMessage
	if err := json.Unmarshal(content[:length], &message); err != nil {
		return "", 0, err
	}

	return message.Method, length, nil
}
