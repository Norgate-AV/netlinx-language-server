package server_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Norgate-AV/netlinx-language-server/internal/logger"
	"github.com/Norgate-AV/netlinx-language-server/internal/lsp"
	"github.com/Norgate-AV/netlinx-language-server/internal/server"
	"github.com/Norgate-AV/netlinx-language-server/internal/workspace"
	"github.com/Norgate-AV/netlinx-language-server/parser"

	"github.com/sourcegraph/jsonrpc2"
)

const testDocumentURI = "file:///test.axs"

func TestTextDocumentDidOpen(t *testing.T) {
	// Setup
	log := logger.NewStdLogger()
	ts, err := parser.NewTreeSitter()
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}

	defer ts.Close()

	state := workspace.NewState(&workspace.NewStateOptions{
		TreeSitter: ts,
		Logger:     log,
	})

	srv := server.NewServer(log, state)

	// Create test document parameters
	params := lsp.DidOpenTextDocumentParams{
		TextDocument: lsp.TextDocumentItem{
			URI:  testDocumentURI,
			Text: "PROGRAM_NAME='Test'\nDEFINE_VARIABLE\nINTEGER x",
		},
	}

	// Create request with parameters
	paramsBytes, _ := json.Marshal(params)
	rawParams := json.RawMessage(paramsBytes)
	req := &jsonrpc2.Request{
		Method: lsp.MethodTextDocumentDidOpen,
		Params: &rawParams,
	}

	// Call handler
	srv.TextDocumentDidOpen(context.Background(), nil, req)

	// Verify document was added to state
	document, exists := state.GetDocument(testDocumentURI)
	if !exists {
		t.Fatal("Document was not added to state")
	}

	expectedContent := "PROGRAM_NAME='Test'\nDEFINE_VARIABLE\nINTEGER x"
	if document.Content != expectedContent {
		t.Fatalf("Document content mismatch.\nExpected: %q\nGot: %q", expectedContent, document.Content)
	}
}

func TestTextDocumentDidChange(t *testing.T) {
	// Setup
	log := logger.NewStdLogger()
	ts, err := parser.NewTreeSitter()
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}

	defer ts.Close()

	state := workspace.NewState(&workspace.NewStateOptions{
		TreeSitter: ts,
		Logger:     log,
	})

	srv := server.NewServer(log, state)

	// First add a document to the state
	state.AddDocument(testDocumentURI, "PROGRAM_NAME='Test'\nDEFINE_VARIABLE\nINTEGER x")

	// Create change parameters
	params := lsp.DidChangeTextDocumentParams{
		TextDocument: lsp.VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: lsp.TextDocumentIdentifier{
				URI: testDocumentURI,
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
		Method: lsp.MethodTextDocumentDidChange,
		Params: &rawParams,
	}

	// Call handler
	srv.TextDocumentDidChange(context.Background(), nil, req)

	// Verify document was updated in state
	document, exists := state.GetDocument(testDocumentURI)
	if !exists {
		t.Fatal("Document not found in state after update")
	}

	expectedContent := "PROGRAM_NAME='Updated'\nDEFINE_VARIABLE\nINTEGER x, y"
	if document.Content != expectedContent {
		t.Fatalf("Document content mismatch after update.\nExpected: %q\nGot: %q", expectedContent, document.Content)
	}
}

func TestTextDocumentDidClose(t *testing.T) {
	// Setup
	log := logger.NewStdLogger()
	ts, err := parser.NewTreeSitter()
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}

	defer ts.Close()

	state := workspace.NewState(&workspace.NewStateOptions{
		TreeSitter: ts,
		Logger:     log,
	})

	srv := server.NewServer(log, state)

	// First add a document to the state
	state.AddDocument(testDocumentURI, "PROGRAM_NAME='Test'\nDEFINE_VARIABLE\nINTEGER x")

	// Create close parameters
	params := lsp.DidCloseTextDocumentParams{
		TextDocument: lsp.TextDocumentIdentifier{
			URI: testDocumentURI,
		},
	}

	// Create request with parameters
	paramsBytes, _ := json.Marshal(params)
	rawParams := json.RawMessage(paramsBytes)
	req := &jsonrpc2.Request{
		Method: lsp.MethodTextDocumentDidClose,
		Params: &rawParams,
	}

	// Call handler
	srv.TextDocumentDidClose(context.Background(), nil, req)

	// Verify document was removed from state
	_, exists := state.GetDocument(testDocumentURI)
	if exists {
		t.Fatal("Document still exists in state after close")
	}
}

func TestInvalidParameters(t *testing.T) {
	// Setup
	log := logger.NewStdLogger()
	ts, err := parser.NewTreeSitter()
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}

	defer ts.Close()

	state := workspace.NewState(&workspace.NewStateOptions{
		TreeSitter: ts,
		Logger:     log,
	})

	srv := server.NewServer(log, state)

	// Test cases with invalid JSON
	testCases := []struct {
		name   string
		method string
		params string
	}{
		{
			name:   "Invalid didOpen params",
			method: lsp.MethodTextDocumentDidOpen,
			params: `{"textDocument": {"uri": 123}}`, // uri should be a string
		},
		{
			name:   "Invalid didChange params",
			method: lsp.MethodTextDocumentDidChange,
			params: `{"textDocument": {"uri": "file:///test.axs"}}`, // missing contentChanges
		},
		{
			name:   "Invalid didClose params",
			method: lsp.MethodTextDocumentDidClose,
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
			case lsp.MethodTextDocumentDidOpen:
				srv.TextDocumentDidOpen(context.Background(), nil, req)
			case lsp.MethodTextDocumentDidChange:
				srv.TextDocumentDidChange(context.Background(), nil, req)
			case lsp.MethodTextDocumentDidClose:
				srv.TextDocumentDidClose(context.Background(), nil, req)
			}

			// No assertion needed here - just making sure the handlers
			// handle invalid parameters gracefully without panicking
		})
	}
}
