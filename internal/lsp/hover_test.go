package lsp

import (
    "encoding/json"
    "testing"
)

func TestHoverSerialization(t *testing.T) {
    hover := Hover{
        Contents: MarkupContent{
            Kind:  MarkupKindMarkdown,
            Value: "**Function**: Test\n\nThis is a test function",
        },
        Range: &Range{
            Start: Position{Line: 5, Character: 10},
            End:   Position{Line: 5, Character: 15},
        },
    }

    data, err := json.Marshal(hover)
    if err != nil {
        t.Fatalf("Failed to marshal Hover: %v", err)
    }

    var unmarshaled Hover
    if err := json.Unmarshal(data, &unmarshaled); err != nil {
        t.Fatalf("Failed to unmarshal Hover: %v", err)
    }

    if unmarshaled.Contents.Kind != MarkupKindMarkdown {
        t.Errorf("Expected MarkupKindMarkdown, got %s", unmarshaled.Contents.Kind)
    }
    if unmarshaled.Contents.Value != "**Function**: Test\n\nThis is a test function" {
        t.Errorf("Expected hover content incorrect")
    }
}
