package server_test

import (
	"encoding/json"
	"testing"

	"github.com/Norgate-AV/netlinx-language-server/internal/analysis"
	"github.com/Norgate-AV/netlinx-language-server/internal/logger"
	"github.com/Norgate-AV/netlinx-language-server/internal/lsp"
	"github.com/Norgate-AV/netlinx-language-server/internal/server"

	"github.com/sourcegraph/jsonrpc2"
)

func TestLSPHandlerCreation(t *testing.T) {
	logger := logger.NewStdLogger()
	state := analysis.NewState()

	server := server.NewServer(logger, state)
	if server == nil {
		t.Fatal("Expected non-nil handler")
	}

	state.AddDocument("file:///test.axs", "test content")

	content, ok := state.GetDocument("file:///test.axs")
	if !ok {
		t.Fatal("Expected document to be added to state")
	}

	if content != "test content" {
		t.Errorf("Expected document content 'test content', got '%s'", content)
	}
}

func TestInitializeResultJSON(t *testing.T) {
	result := server.NewInitializeResponse()

	// Test that the result can be marshaled to JSON
	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal InitializeResult: %v", err)
	}

	// Unmarshal and verify the contents
	var unmarshaled lsp.InitializeResult
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

func TestCreateError(t *testing.T) {
	err := &jsonrpc2.Error{
		Code:    jsonrpc2.CodeInternalError,
		Message: "Test error message",
	}

	if err.Code != jsonrpc2.CodeInternalError {
		t.Errorf("Expected error code %d, got %d", jsonrpc2.CodeInternalError, err.Code)
	}

	if err.Message != "Test error message" {
		t.Errorf("Expected error message 'Test error message', got '%s'", err.Message)
	}
}
