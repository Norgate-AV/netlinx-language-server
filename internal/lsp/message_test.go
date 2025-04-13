package lsp_test

import (
	"testing"

	"github.com/Norgate-AV/netlinx-language-server/internal/lsp"
)

func TestMethodTypes(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		isNotification bool
	}{
		// General
		{"Initialize", lsp.MethodInitialize, false},
		{"Initialized", lsp.MethodInitialized, true},
		{"Shutdown", lsp.MethodShutdown, false},
		{"Exit", lsp.MethodExit, true},

		// Document Management
		{"DidOpen", lsp.MethodTextDocumentDidOpen, true},
		{"DidChange", lsp.MethodTextDocumentDidChange, true},
		{"DidClose", lsp.MethodTextDocumentDidClose, true},
		{"DidSave", lsp.MethodTextDocumentDidSave, true},
		{"WillSave", lsp.MethodTextDocumentWillSave, true},
		{"WillSaveWaitUntil", lsp.MethodTextDocumentWillSaveWaitUntil, false},

		// Navigation
		{"Definition", lsp.MethodTextDocumentDefinition, false},
		{"References", lsp.MethodTextDocumentReferences, false},

		// Information
		{"Hover", lsp.MethodTextDocumentHover, false},
		{"DocumentSymbol", lsp.MethodTextDocumentDocumentSymbol, false},

		// Workspace
		{"WorkspaceDidChangeConfig", lsp.MethodWorkspaceDidChangeConfig, true},
		{"WorkspaceSymbol", lsp.MethodWorkspaceSymbol, false},

		// Special
		{"CancelRequest", lsp.MethodCancelRequest, true},
		{"RegisterCapability", lsp.MethodClientRegisterCapability, false},

		// Server-to-Client
		{"PublishDiagnostics", lsp.MethodTextDocumentPublishDiagnostics, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := lsp.IsNotification(tc.method)
			if result != tc.isNotification {
				t.Errorf("IsNotification(%q) = %v; want %v", tc.method, result, tc.isNotification)
			}
		})
	}
}

func TestUnknownMethods(t *testing.T) {
	// Unknown methods should default to request (false)
	unknownMethods := []string{
		"unknown/method",
		"",
		"textDocument/unknownOperation",
		"$/customExtension",
	}

	for _, method := range unknownMethods {
		if lsp.IsNotification(method) {
			t.Errorf("Expected unknown method %q to default to request type", method)
		}
	}
}

func TestAllMethodConstants(t *testing.T) {
	// This test ensures all our method constants can be processed
	// by IsNotification without panicking
	constants := []string{
		// Just testing a representative sample
		lsp.MethodInitialize,
		lsp.MethodTextDocumentDidOpen,
		lsp.MethodTextDocumentHover,
		lsp.MethodWorkspaceSymbol,
		lsp.MethodCancelRequest,
	}

	for _, method := range constants {
		// Just verify it doesn't panic
		_ = lsp.IsNotification(method)
	}
}

func TestMethodPatterns(t *testing.T) {
	// Test patterns in method naming

	// All "did" methods should be notifications
	didMethods := []string{
		lsp.MethodTextDocumentDidOpen,
		lsp.MethodTextDocumentDidChange,
		lsp.MethodTextDocumentDidClose,
		lsp.MethodTextDocumentDidSave,
		lsp.MethodWorkspaceDidChangeConfig,
		lsp.MethodWorkspaceDidChangeWatchedFiles,
		lsp.MethodWorkspaceDidChangeWorkspaceFolders,
	}

	for _, method := range didMethods {
		if !lsp.IsNotification(method) {
			t.Errorf("Expected 'did' method %q to be a notification", method)
		}
	}
}
