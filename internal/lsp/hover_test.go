package lsp_test

import (
	"encoding/json"
	"testing"

	"github.com/Norgate-AV/netlinx-language-server/internal/lsp"
)

func TestHoverSerialization(t *testing.T) {
	hover := lsp.Hover{
		Contents: lsp.MarkupContent{
			Kind:  lsp.MarkupKindMarkdown,
			Value: "**Function**: Test\n\nThis is a test function",
		},
		Range: &lsp.Range{
			Start: lsp.Position{Line: 5, Character: 10},
			End:   lsp.Position{Line: 5, Character: 15},
		},
	}

	data, err := json.Marshal(hover)
	if err != nil {
		t.Fatalf("Failed to marshal Hover: %v", err)
	}

	var unmarshaled lsp.Hover
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Hover: %v", err)
	}

	if unmarshaled.Contents.Kind != lsp.MarkupKindMarkdown {
		t.Errorf("Expected MarkupKindMarkdown, got %s", unmarshaled.Contents.Kind)
	}
	if unmarshaled.Contents.Value != "**Function**: Test\n\nThis is a test function" {
		t.Errorf("Expected hover content incorrect")
	}
}
