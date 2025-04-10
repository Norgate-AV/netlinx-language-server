package lsp

import (
	"encoding/json"
	"testing"
)

func TestPositionSerialization(t *testing.T) {
	position := Position{
		Line:      10,
		Character: 20,
	}

	data, err := json.Marshal(position)
	if err != nil {
		t.Fatalf("Failed to marshal Position: %v", err)
	}

	var unmarshaled Position
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Position: %v", err)
	}

	if unmarshaled.Line != 10 {
		t.Errorf("Expected line 10, got %d", unmarshaled.Line)
	}
	if unmarshaled.Character != 20 {
		t.Errorf("Expected character 20, got %d", unmarshaled.Character)
	}
}

func TestRangeSerialization(t *testing.T) {
	rangeObj := Range{
		Start: Position{Line: 5, Character: 10},
		End:   Position{Line: 5, Character: 20},
	}

	data, err := json.Marshal(rangeObj)
	if err != nil {
		t.Fatalf("Failed to marshal Range: %v", err)
	}

	var unmarshaled Range
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Range: %v", err)
	}

	if unmarshaled.Start.Line != 5 {
		t.Errorf("Expected start line 5, got %d", unmarshaled.Start.Line)
	}
}
