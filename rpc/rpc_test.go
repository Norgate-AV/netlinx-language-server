package rpc_test

import (
	"testing"

	"github.com/Norgate-AV/netlinx-language-server/rpc"
)

type EncodingExample struct {
	Testing bool `json:"testing"`
}

func TestEncode(t *testing.T) {
	expected := "Content-Length: 16\r\n\r\n{\"testing\":true}"
	actual := rpc.EncodeMessage(EncodingExample{Testing: true})

	if expected != actual {
		t.Fatalf("Expected %s, Actual %s", expected, actual)
	}
}

func TestDecode(t *testing.T) {
	message := "Content-Length: 15\r\n\r\n{\"method\":\"hi\"}"
	method, length, err := rpc.DecodeMessage([]byte(message))
	if err != nil {
		t.Fatal(err)
	}

	if length != 15 {
		t.Fatalf("Expected length 15, Actual %d", length)
	}

	if method != "hi" {
		t.Fatalf("Expected method hi, Actual %s", method)
	}
}
