package server_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/Norgate-AV/netlinx-language-server/internal/server"
)

func TestInitializeResponseFormat(t *testing.T) {
	// Get the initialize response
	response := server.NewInitializeResponse()

	// Marshal to JSON
	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal InitializeResult: %v", err)
	}

	// Check raw JSON format for capability structure
	var rawJSON map[string]any
	if err := json.Unmarshal(data, &rawJSON); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Validate capabilities field exists and is an object
	capabilities, ok := rawJSON["capabilities"].(map[string]any)
	if !ok {
		t.Fatal("Expected capabilities to be an object")
	}

	// Check if hoverProvider is true
	hoverProvider, ok := capabilities["hoverProvider"].(bool)
	if !ok || !hoverProvider {
		t.Fatal("Expected hoverProvider to be true")
	}

	// Verify complete JSON structure against LSP spec
	expectedJSON := `{
        "capabilities": {
            "textDocumentSync": 2,
            "hoverProvider": true,
            "documentSymbolProvider": true,
            "diagnosticProvider": {
                "interFileDependencies": false,
                "workspaceDiagnostics": false
            }
        },
        "serverInfo": {
            "name": "netlinx-language-server",
            "version": "0.1.0"
        }
    }`

	// Compare JSON structure (ignoring whitespace)
	var expected, actual any
	if err := json.Unmarshal([]byte(expectedJSON), &expected); err != nil {
		t.Fatalf("Invalid expected JSON: %v", err)
	}
	if err := json.Unmarshal(data, &actual); err != nil {
		t.Fatalf("Invalid actual JSON: %v", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Initialize response format doesn't match expected format\nExpected: %v\nGot: %v",
			expected, actual)
	}
}
