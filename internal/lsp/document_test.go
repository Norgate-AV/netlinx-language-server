package lsp

import (
	"encoding/json"
	"testing"
)

func TestDidOpenTextDocumentParamsSerialization(t *testing.T) {
	params := DidOpenTextDocumentParams{
		TextDocument: TextDocumentItem{
			URI:        "file:///test.axs",
			LanguageID: "netlinx",
			Version:    1,
			Text:       "PROGRAM_NAME='Test'",
		},
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Failed to marshal DidOpenTextDocumentParams: %v", err)
	}

	var unmarshaled DidOpenTextDocumentParams
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal DidOpenTextDocumentParams: %v", err)
	}

	if unmarshaled.TextDocument.URI != "file:///test.axs" {
		t.Errorf("Expected URI 'file:///test.axs', got '%s'", unmarshaled.TextDocument.URI)
	}
	if unmarshaled.TextDocument.Text != "PROGRAM_NAME='Test'" {
		t.Errorf("Expected text 'PROGRAM_NAME='Test'', got '%s'", unmarshaled.TextDocument.Text)
	}
}

func TestTextDocumentContentChangeEventSerialization(t *testing.T) {
	line := uint32(5)
	character := uint32(10)
	changeEvent := TextDocumentContentChangeEvent{
		Range: &Range{
			Start: Position{Line: line, Character: character},
			End:   Position{Line: line, Character: character + 5},
		},
		Text: "NewText",
	}

	data, err := json.Marshal(changeEvent)
	if err != nil {
		t.Fatalf("Failed to marshal TextDocumentContentChangeEvent: %v", err)
	}

	var unmarshaled TextDocumentContentChangeEvent
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal TextDocumentContentChangeEvent: %v", err)
	}

	if unmarshaled.Range.Start.Line != line {
		t.Errorf("Expected line %d, got %d", line, unmarshaled.Range.Start.Line)
	}
	if unmarshaled.Text != "NewText" {
		t.Errorf("Expected text 'NewText', got '%s'", unmarshaled.Text)
	}
}
