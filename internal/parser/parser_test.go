package parser

import (
	"testing"
)

func TestBasicParsing(t *testing.T) {
	p := NewParser()

	code := `PROGRAM_NAME='Test'
DEFINE_DEVICE
dvTP = 10001:1:0
`

	tree, err := p.Parse([]byte(code))
	if err != nil {
		t.Fatalf("Failed to parse valid NetLinx code: %v", err)
	}

	// Check root node exists
	root := tree.RootNode()
	if root == nil {
		t.Fatal("Expected non-nil root node")
	}

	// Basic validation of the parse tree
	if root.ChildCount() < 1 {
		t.Errorf("Expected at least one child node, got %d", root.ChildCount())
	}
}

func TestInvalidSyntax(t *testing.T) {
	p := NewParser()

	code := `PROGRAM_NAME='Unterminated`

	_, err := p.Parse([]byte(code))
	if err == nil {
		t.Fatal("Expected error for invalid syntax, got nil")
	}
}
