package workspace_test

import (
	"testing"

	"github.com/Norgate-AV/netlinx-language-server/internal/logger"
	"github.com/Norgate-AV/netlinx-language-server/internal/workspace"
	"github.com/Norgate-AV/netlinx-language-server/parser"
)

const testDocumentURI = "file:///test.axs"

func TestDocumentManagement(t *testing.T) {
	log := logger.NewStdLogger()
	ts, err := parser.NewTreeSitter()
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}

	defer ts.Close()

	state := workspace.NewWorkspace(&workspace.Options{
		TreeSitter: ts,
		Logger:     log,
	})

	// Test adding a document
	state.AddDocument(testDocumentURI, "content")

	// Test retrieving a document
	document, ok := state.GetDocument(testDocumentURI)
	if !ok {
		t.Fatal("Expected document to exist")
	}
	if document.Content != "content" {
		t.Errorf("Expected content 'content', got '%s'", document.Content)
	}

	// Test updating a document
	state.UpdateDocument(testDocumentURI, "updated")
	document, _ = state.GetDocument(testDocumentURI)
	if document.Content != "updated" {
		t.Errorf("Expected content 'updated', got '%s'", document.Content)
	}

	// Test closing a document
	state.CloseDocument(testDocumentURI)
	_, ok = state.GetDocument(testDocumentURI)
	if ok {
		t.Error("Expected document to be removed")
	}
}
