package analysis

import (
	"testing"
)

func TestDocumentManagement(t *testing.T) {
	state := NewState()

	// Test adding a document
	state.AddDocument("file:///test.axs", "content")

	// Test retrieving a document
	content, ok := state.GetDocument("file:///test.axs")
	if !ok {
		t.Fatal("Expected document to exist")
	}
	if content != "content" {
		t.Errorf("Expected content 'content', got '%s'", content)
	}

	// Test updating a document
	state.UpdateDocument("file:///test.axs", "updated")
	content, _ = state.GetDocument("file:///test.axs")
	if content != "updated" {
		t.Errorf("Expected content 'updated', got '%s'", content)
	}

	// Test closing a document
	state.CloseDocument("file:///test.axs")
	_, ok = state.GetDocument("file:///test.axs")
	if ok {
		t.Error("Expected document to be removed")
	}
}
