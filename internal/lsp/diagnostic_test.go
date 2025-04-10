package lsp

import (
	"encoding/json"
	"testing"
)

func TestDiagnosticSerialization(t *testing.T) {
	severity := DiagnosticSeverityError
	source := "netlinx-lsp"
	diagnostic := Diagnostic{
		Range: Range{
			Start: Position{Line: 10, Character: 5},
			End:   Position{Line: 10, Character: 10},
		},
		Severity: &severity,
		Source:   &source,
		Message:  "Undefined variable",
	}

	data, err := json.Marshal(diagnostic)
	if err != nil {
		t.Fatalf("Failed to marshal Diagnostic: %v", err)
	}

	var unmarshaled Diagnostic
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Diagnostic: %v", err)
	}

	if unmarshaled.Message != "Undefined variable" {
		t.Errorf("Expected message 'Undefined variable', got '%s'", unmarshaled.Message)
	}
	if *unmarshaled.Severity != DiagnosticSeverityError {
		t.Errorf("Expected severity %d, got %d", DiagnosticSeverityError, *unmarshaled.Severity)
	}
}

func TestPublishDiagnosticsParamsSerialization(t *testing.T) {
	severity := DiagnosticSeverityError
	source := "netlinx-lsp"
	params := PublishDiagnosticsParams{
		URI: "file:///test.axs",
		Diagnostics: []Diagnostic{
			{
				Range: Range{
					Start: Position{Line: 10, Character: 5},
					End:   Position{Line: 10, Character: 10},
				},
				Severity: &severity,
				Source:   &source,
				Message:  "Undefined variable",
			},
		},
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Failed to marshal PublishDiagnosticsParams: %v", err)
	}

	var unmarshaled PublishDiagnosticsParams
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal PublishDiagnosticsParams: %v", err)
	}

	if unmarshaled.URI != "file:///test.axs" {
		t.Errorf("Expected URI 'file:///test.axs', got '%s'", unmarshaled.URI)
	}
	if len(unmarshaled.Diagnostics) != 1 {
		t.Fatalf("Expected 1 diagnostic, got %d", len(unmarshaled.Diagnostics))
	}
}
