package lsp

import (
	"encoding/json"
	"testing"
)

func TestInitializeResultSerialization(t *testing.T) {
	b := func(v bool) *bool { return &v }

	result := InitializeResult{
		Capabilities: ServerCapabilities{
			TextDocumentSync: TextDocumentSyncKindIncremental,
			HoverProvider:    b(true),
		},
		ServerInfo: ServerInfo{
			Name:    "TestServer",
			Version: "1.0.0",
		},
	}

	// Serialize the result to JSON
	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to serialize InitializeResult: %v", err)
	}

	// Verify the JSON output
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		t.Fatalf("Failed to deserialize JSON: %v", err)
	}

	// Check if the capabilities are present
	capabilities, ok := jsonMap["capabilities"].(map[string]interface{})
	if !ok {
		t.Fatal("Missing or invalid capabilities field")
	}

	// Verify hover support is enabled
	if hoverProvider, ok := capabilities["hoverProvider"].(bool); !ok || !hoverProvider {
		t.Errorf("Expected hoverProvider to be true, got %v", capabilities["hoverProvider"])
	}

	// Verify text document sync is set correctly
	if textDocSync, ok := capabilities["textDocumentSync"].(float64); !ok || int(textDocSync) != int(TextDocumentSyncKindIncremental) {
		t.Errorf("Expected textDocumentSync to be %d, got %v", TextDocumentSyncKindIncremental, capabilities["textDocumentSync"])
	}

	// Check server info
	serverInfo, ok := jsonMap["serverInfo"].(map[string]interface{})
	if !ok {
		t.Fatal("Missing or invalid serverInfo field")
	}

	if name, ok := serverInfo["name"].(string); !ok || name != "TestServer" {
		t.Errorf("Expected server name to be 'TestServer', got %v", serverInfo["name"])
	}
}

func TestInitializeRequestParamsSerialization(t *testing.T) {
	jsonData := `{
		"clientInfo": {
			"name": "TestClient",
			"version": "1.0.0"
		}
	}`

	var params InitializeRequestParams
	if err := json.Unmarshal([]byte(jsonData), &params); err != nil {
		t.Fatalf("Failed to unmarshal InitializeRequestParams: %v", err)
	}

	// Verify data
	if params.ClientInfo == nil {
		t.Fatal("ClientInfo is nil")
	}

	if params.ClientInfo.Name != "TestClient" {
		t.Errorf("Expected client name 'TestClient', got '%s'", params.ClientInfo.Name)
	}

	if params.ClientInfo.Version != "1.0.0" {
		t.Errorf("Expected client version '1.0.0', got '%s'", params.ClientInfo.Version)
	}

	// Serialize back to JSON
	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Failed to marshal InitializeRequestParams: %v", err)
	}

	// Verify the JSON output
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		t.Fatalf("Failed to unmarshal JSON to map: %v", err)
	}

	clientInfo, ok := jsonMap["clientInfo"].(map[string]interface{})
	if !ok {
		t.Fatal("Missing or invalid clientInfo field")
	}

	if name, ok := clientInfo["name"].(string); !ok || name != "TestClient" {
		t.Errorf("Expected client name 'TestClient', got '%s'", clientInfo["name"])
	}

	if version, ok := clientInfo["version"].(string); !ok || version != "1.0.0" {
		t.Errorf("Expected client version '1.0.0', got '%s'", clientInfo["version"])
	}
}
