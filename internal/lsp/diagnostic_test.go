package lsp_test

import (
	"encoding/json"
	"testing"

	"github.com/Norgate-AV/netlinx-language-server/internal/lsp"
)

func TestDiagnosticSerialization(t *testing.T) {
	severity := lsp.DiagnosticSeverityError
	source := "netlinx-lsp"
	diagnostic := lsp.Diagnostic{
		Range: lsp.Range{
			Start: lsp.Position{Line: 10, Character: 5},
			End:   lsp.Position{Line: 10, Character: 10},
		},
		Severity: &severity,
		Source:   &source,
		Message:  "Undefined variable",
	}

	data, err := json.Marshal(diagnostic)
	if err != nil {
		t.Fatalf("Failed to marshal Diagnostic: %v", err)
	}

	var unmarshaled lsp.Diagnostic
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Diagnostic: %v", err)
	}

	if unmarshaled.Message != "Undefined variable" {
		t.Errorf("Expected message 'Undefined variable', got '%s'", unmarshaled.Message)
	}
	if *unmarshaled.Severity != lsp.DiagnosticSeverityError {
		t.Errorf("Expected severity %d, got %d", lsp.DiagnosticSeverityError, *unmarshaled.Severity)
	}
}

func TestPublishDiagnosticsParamsSerialization(t *testing.T) {
	severity := lsp.DiagnosticSeverityError
	source := "netlinx-lsp"
	params := lsp.PublishDiagnosticsParams{
		URI: "file:///test.axs",
		Diagnostics: []lsp.Diagnostic{
			{
				Range: lsp.Range{
					Start: lsp.Position{Line: 10, Character: 5},
					End:   lsp.Position{Line: 10, Character: 10},
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

	var unmarshaled lsp.PublishDiagnosticsParams
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
