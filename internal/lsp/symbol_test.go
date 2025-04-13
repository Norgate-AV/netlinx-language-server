package lsp_test

import (
	"encoding/json"
	"testing"

	"github.com/Norgate-AV/netlinx-language-server/internal/lsp"
)

func TestSymbolInformationSerialization(t *testing.T) {
	symbol := lsp.SymbolInformation{
		Name: "MyFunction",
		Kind: 12, // Function
		Location: lsp.Location{
			URI: "file:///test.axs",
			Range: lsp.Range{
				Start: lsp.Position{Line: 10, Character: 0},
				End:   lsp.Position{Line: 15, Character: 1},
			},
		},
		ContainerName: "ModuleName",
	}

	data, err := json.Marshal(symbol)
	if err != nil {
		t.Fatalf("Failed to marshal SymbolInformation: %v", err)
	}

	var unmarshaled lsp.SymbolInformation
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
	symbol := lsp.DocumentSymbol{
		Name:           "DEFINE_DEVICE",
		Detail:         "Device section",
		Kind:           2, // Module
		Range:          lsp.Range{Start: lsp.Position{Line: 5, Character: 0}, End: lsp.Position{Line: 10, Character: 0}},
		SelectionRange: lsp.Range{Start: lsp.Position{Line: 5, Character: 0}, End: lsp.Position{Line: 5, Character: 14}},
		Children: []lsp.DocumentSymbol{
			{
				Name:           "dvTP",
				Detail:         "10001:1:0",
				Kind:           13, // Variable
				Range:          lsp.Range{Start: lsp.Position{Line: 6, Character: 0}, End: lsp.Position{Line: 6, Character: 15}},
				SelectionRange: lsp.Range{Start: lsp.Position{Line: 6, Character: 0}, End: lsp.Position{Line: 6, Character: 4}},
			},
		},
	}

	data, err := json.Marshal(symbol)
	if err != nil {
		t.Fatalf("Failed to marshal DocumentSymbol: %v", err)
	}

	var unmarshaled lsp.DocumentSymbol
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
	params := lsp.DocumentSymbolParams{
		TextDocument: lsp.TextDocumentIdentifier{
			URI: "file:///test.axs",
		},
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Failed to marshal DocumentSymbolParams: %v", err)
	}

	var unmarshaled lsp.DocumentSymbolParams
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal DocumentSymbolParams: %v", err)
	}

	if unmarshaled.TextDocument.URI != "file:///test.axs" {
		t.Errorf("Expected URI 'file:///test.axs', got '%s'", unmarshaled.TextDocument.URI)
	}
}
