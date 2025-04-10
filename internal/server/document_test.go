package server

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Norgate-AV/netlinx-language-server/internal/analysis"
	"github.com/Norgate-AV/netlinx-language-server/internal/logger"
	"github.com/Norgate-AV/netlinx-language-server/internal/lsp"

	"github.com/sourcegraph/jsonrpc2"
)

func TestTextDocumentDidOpen(t *testing.T) {
	// Setup
	log := logger.NewStdLogger()
	state := analysis.NewState()
	srv := NewServer(log, state)

	// Create test document parameters
	params := lsp.DidOpenTextDocumentParams{
		TextDocument: lsp.TextDocumentItem{
			URI:  "file:///test.axs",
			Text: "PROGRAM_NAME='Test'\nDEFINE_VARIABLE\nINTEGER x",
		},
	}

	// Create request with parameters
	paramsBytes, _ := json.Marshal(params)
	rawParams := json.RawMessage(paramsBytes)
	req := &jsonrpc2.Request{
		Method: "textDocument/didOpen",
		Params: &rawParams,
	}

	// Call handler
	srv.TextDocumentDidOpen(context.Background(), nil, req)

	// Verify document was added to state
	content, exists := state.GetDocument("file:///test.axs")
	if !exists {
		t.Fatal("Document was not added to state")
	}

	expectedContent := "PROGRAM_NAME='Test'\nDEFINE_VARIABLE\nINTEGER x"
	if content != expectedContent {
		t.Fatalf("Document content mismatch.\nExpected: %q\nGot: %q", expectedContent, content)
	}
}

func TestTextDocumentDidChange(t *testing.T) {
	// Setup
	log := logger.NewStdLogger()
	state := analysis.NewState()
	srv := NewServer(log, state)

	// First add a document to the state
	state.AddDocument("file:///test.axs", "PROGRAM_NAME='Test'\nDEFINE_VARIABLE\nINTEGER x")

	// Create change parameters
	params := lsp.DidChangeTextDocumentParams{
		TextDocument: lsp.VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: lsp.TextDocumentIdentifier{
				URI: "file:///test.axs",
			},
			Version: 2,
		},
		ContentChanges: []lsp.TextDocumentContentChangeEvent{
			{
				Text: "PROGRAM_NAME='Updated'\nDEFINE_VARIABLE\nINTEGER x, y",
			},
		},
	}

	// Create request with parameters
	paramsBytes, _ := json.Marshal(params)
	rawParams := json.RawMessage(paramsBytes)
	req := &jsonrpc2.Request{
		Method: "textDocument/didChange",
		Params: &rawParams,
	}

	// Call handler
	srv.TextDocumentDidChange(context.Background(), nil, req)

	// Verify document was updated in state
	content, exists := state.GetDocument("file:///test.axs")
	if !exists {
		t.Fatal("Document not found in state after update")
	}

	expectedContent := "PROGRAM_NAME='Updated'\nDEFINE_VARIABLE\nINTEGER x, y"
	if content != expectedContent {
		t.Fatalf("Document content mismatch after update.\nExpected: %q\nGot: %q", expectedContent, content)
	}
}

func TestTextDocumentDidClose(t *testing.T) {
	// Setup
	log := logger.NewStdLogger()
	state := analysis.NewState()
	srv := NewServer(log, state)

	// First add a document to the state
	state.AddDocument("file:///test.axs", "PROGRAM_NAME='Test'\nDEFINE_VARIABLE\nINTEGER x")

	// Create close parameters
	params := lsp.DidCloseTextDocumentParams{
		TextDocument: lsp.TextDocumentIdentifier{
			URI: "file:///test.axs",
		},
	}

	// Create request with parameters
	paramsBytes, _ := json.Marshal(params)
	rawParams := json.RawMessage(paramsBytes)
	req := &jsonrpc2.Request{
		Method: "textDocument/didClose",
		Params: &rawParams,
	}

	// Call handler
	srv.TextDocumentDidClose(context.Background(), nil, req)

	// Verify document was removed from state
	_, exists := state.GetDocument("file:///test.axs")
	if exists {
		t.Fatal("Document still exists in state after close")
	}
}

func TestInvalidParameters(t *testing.T) {
	// Setup
	log := logger.NewStdLogger()
	state := analysis.NewState()
	srv := NewServer(log, state)

	// Test cases with invalid JSON
	testCases := []struct {
		name   string
		method string
		params string
	}{
		{
			name:   "Invalid didOpen params",
			method: "textDocument/didOpen",
			params: `{"textDocument": {"uri": 123}}`, // uri should be a string
		},
		{
			name:   "Invalid didChange params",
			method: "textDocument/didChange",
			params: `{"textDocument": {"uri": "file:///test.axs"}}`, // missing contentChanges
		},
		{
			name:   "Invalid didClose params",
			method: "textDocument/didClose",
			params: `{"textDocument": {"uri": []}}`, // uri should be a string
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create invalid request
			rawParams := json.RawMessage(tc.params)
			req := &jsonrpc2.Request{
				Method: tc.method,
				Params: &rawParams,
			}

			// Depending on the method, call the appropriate handler
			// These should not panic even with invalid parameters
			switch tc.method {
			case "textDocument/didOpen":
				srv.TextDocumentDidOpen(context.Background(), nil, req)
			case "textDocument/didChange":
				srv.TextDocumentDidChange(context.Background(), nil, req)
			case "textDocument/didClose":
				srv.TextDocumentDidClose(context.Background(), nil, req)
			}

			// No assertion needed here - just making sure the handlers
			// handle invalid parameters gracefully without panicking
		})
	}
}
