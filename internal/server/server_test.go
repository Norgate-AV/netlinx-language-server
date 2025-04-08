package server_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/Norgate-AV/netlinx-language-server/internal/analysis"
	"github.com/Norgate-AV/netlinx-language-server/internal/protocol"
	"github.com/Norgate-AV/netlinx-language-server/internal/server"
	"github.com/sourcegraph/jsonrpc2"
)

// TestLSPHandlerCreation tests that the LSP handler can be created and basic methods
func TestLSPHandlerCreation(t *testing.T) {
	// Setup a logger and state
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	state := &analysis.State{
		Documents: make(map[string]protocol.TextDocumentUri),
	}

	// Create a new handler
	handler := server.NewLSPHandler(logger, state)
	if handler == nil {
		t.Fatal("Expected non-nil handler")
	}

	// Test that the state's OpenDocument method works
	state.AddDocument("file:///test.axs", "test content")
	content, ok := state.GetDocument("file:///test.axs")
	if !ok {
		t.Fatal("Expected document to be added to state")
	}
	if content != "test content" {
		t.Errorf("Expected document content 'test content', got '%s'", content)
	}
}

// TestInitializeResultJSON tests that our InitializeResult struct marshals to JSON correctly
func TestInitializeResultJSON(t *testing.T) {
	result := protocol.NewInitializeResponse(0)

	// Test that the result can be marshaled to JSON
	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal InitializeResult: %v", err)
	}

	// Unmarshal and verify the contents
	var unmarshaled protocol.InitializeResult
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal InitializeResult: %v", err)
	}

	// Verify server info
	if unmarshaled.ServerInfo.Name != "netlinx-language-server" {
		t.Errorf("Expected server name 'netlinx-language-server', got '%s'", unmarshaled.ServerInfo.Name)
	}

	// Check that the server version exists
	if unmarshaled.ServerInfo.Version == "" {
		t.Error("Expected non-empty server version")
	}
}

// TestCreateError tests the creation of a jsonrpc2 error
func TestCreateError(t *testing.T) {
	// Create a test error
	err := &jsonrpc2.Error{
		Code:    jsonrpc2.CodeInternalError,
		Message: "Test error message",
	}

	// Test that the error has the expected code and message
	if err.Code != jsonrpc2.CodeInternalError {
		t.Errorf("Expected error code %d, got %d", jsonrpc2.CodeInternalError, err.Code)
	}
	if err.Message != "Test error message" {
		t.Errorf("Expected error message 'Test error message', got '%s'", err.Message)
	}
}
