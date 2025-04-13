package parser_test

import (
	"fmt"
	"testing"

	"github.com/Norgate-AV/netlinx-language-server/parser"
)

func TestBasicParsing(t *testing.T) {
	ts, err := parser.NewTreeSitter()
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}

	defer ts.Close()

	code := `
		PROGRAM_NAME='Test'
		DEFINE_DEVICE
		dvTP = 10001:1:0
	`

	tree := ts.Parser.Parse([]byte(code), nil)
	defer tree.Close()
	if tree == nil {
		t.Fatal("Expected non-nil parse tree")
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
	ts, err := parser.NewTreeSitter()
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}

	defer ts.Close()

	code := `PROGRAM_NAME='Unterminated`

	tree := ts.Parser.Parse([]byte(code), nil)
	defer tree.Close()
	if tree == nil {
		t.Fatal("Expected non-nil parse tree")
	}

	if !tree.RootNode().HasError() {
		t.Fatal("Expected parse tree to have error node")
	}
}

func TestQueryParsing(t *testing.T) {
	ts, err := parser.NewTreeSitter()
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}

	defer ts.Close()

	sourceCode := []byte(`
		PROGRAM_NAME='Test'

		DEFINE_DEVICE
		dvTP = 10001:1:0

		DEFINE_CONSTANT
		constant integer FOO = 1

		DEFINE_VARIABLE
		volatile integer bar = 2
	`)

	tree := ts.Parser.Parse(sourceCode, nil)
	defer tree.Close()

	query, err := parser.CreateQuery(
		`
		(section) @section
		`,
	)
	if err != nil {
		t.Fatalf("Failed to create query: %v", err)
	}

	defer query.Close()

	qc := parser.CreateQueryCursor()
	defer qc.Close()

	captures := qc.Captures(query, tree.RootNode(), sourceCode)

	for match, index := captures.Next(); match != nil; match, index = captures.Next() {
		fmt.Printf(
			"Capture %d: %s\n",
			index,
			match.Captures[index].Node.Utf8Text(sourceCode),
		)
	}
}
