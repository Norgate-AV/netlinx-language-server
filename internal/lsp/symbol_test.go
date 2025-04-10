package lsp

import (
	"encoding/json"
	"testing"
)

func TestSymbolInformationSerialization(t *testing.T) {
	symbol := SymbolInformation{
		Name: "MyFunction",
		Kind: 12, // Function
		Location: Location{
			URI: "file:///test.axs",
			Range: Range{
				Start: Position{Line: 10, Character: 0},
				End:   Position{Line: 15, Character: 1},
			},
		},
		ContainerName: "ModuleName",
	}

	data, err := json.Marshal(symbol)
	if err != nil {
		t.Fatalf("Failed to marshal SymbolInformation: %v", err)
	}

	var unmarshaled SymbolInformation
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal SymbolInformation: %v", err)
	}

	if unmarshaled.Name != "MyFunction" {
		t.Errorf("Expected name 'MyFunction', got '%s'", unmarshaled.Name)
	}

	if unmarshaled.Kind != 12 {
		t.Errorf("Expected kind 12, got %d", unmarshaled.Kind)
	}
}

func TestDocumentSymbolSerialization(t *testing.T) {
	symbol := DocumentSymbol{
		Name:           "DEFINE_DEVICE",
		Detail:         "Device section",
		Kind:           2, // Module
		Range:          Range{Start: Position{Line: 5, Character: 0}, End: Position{Line: 10, Character: 0}},
		SelectionRange: Range{Start: Position{Line: 5, Character: 0}, End: Position{Line: 5, Character: 14}},
		Children: []DocumentSymbol{
			{
				Name:           "dvTP",
				Detail:         "10001:1:0",
				Kind:           13, // Variable
				Range:          Range{Start: Position{Line: 6, Character: 0}, End: Position{Line: 6, Character: 15}},
				SelectionRange: Range{Start: Position{Line: 6, Character: 0}, End: Position{Line: 6, Character: 4}},
			},
		},
	}

	data, err := json.Marshal(symbol)
	if err != nil {
		t.Fatalf("Failed to marshal DocumentSymbol: %v", err)
	}

	var unmarshaled DocumentSymbol
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal DocumentSymbol: %v", err)
	}

	if unmarshaled.Name != "DEFINE_DEVICE" {
		t.Errorf("Expected name 'DEFINE_DEVICE', got '%s'", unmarshaled.Name)
	}

	if len(unmarshaled.Children) != 1 {
		t.Fatalf("Expected 1 child, got %d", len(unmarshaled.Children))
	}

	if unmarshaled.Children[0].Name != "dvTP" {
		t.Errorf("Expected child name 'dvTP', got '%s'", unmarshaled.Children[0].Name)
	}
}

func TestDocumentSymbolParamsSerialization(t *testing.T) {
	params := DocumentSymbolParams{
		TextDocument: TextDocumentIdentifier{
			URI: "file:///test.axs",
		},
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Failed to marshal DocumentSymbolParams: %v", err)
	}

	var unmarshaled DocumentSymbolParams
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal DocumentSymbolParams: %v", err)
	}

	if unmarshaled.TextDocument.URI != "file:///test.axs" {
		t.Errorf("Expected URI 'file:///test.axs', got '%s'", unmarshaled.TextDocument.URI)
	}
}
