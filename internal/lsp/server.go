package lsp

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ServerCapabilities struct {
	TextDocumentSync                 TextDocumentSyncKind             `json:"textDocumentSync,omitempty"`
	CompletionProvider               *CompletionOptions               `json:"completionProvider,omitempty"`
	HoverProvider                    *bool                            `json:"hoverProvider,omitempty"`
	SignatureHelpProvider            *SignatureHelpOptions            `json:"signatureHelpProvider,omitempty"`
	DeclarationProvider              *bool                            `json:"declarationProvider,omitempty"`
	DefinitionProvider               *bool                            `json:"definitionProvider,omitempty"`
	TypeDefinitionProvider           *bool                            `json:"typeDefinitionProvider,omitempty"`
	ImplementationProvider           *bool                            `json:"implementationProvider,omitempty"`
	ReferencesProvider               *bool                            `json:"referencesProvider,omitempty"`
	DocumentHighlightProvider        *bool                            `json:"documentHighlightProvider,omitempty"`
	DocumentSymbolProvider           *bool                            `json:"documentSymbolProvider,omitempty"`
	CodeActionProvider               *bool                            `json:"codeActionProvider,omitempty"`
	CodeLensProvider                 *CodeLensOptions                 `json:"codeLensProvider,omitempty"`
	DocumentLinkProvider             *DocumentLinkOptions             `json:"documentLinkProvider,omitempty"`
	ColorProvider                    *bool                            `json:"colorProvider,omitempty"`
	DocumentFormattingProvider       *bool                            `json:"documentFormattingProvider,omitempty"`
	DocumentRangeFormattingProvider  *bool                            `json:"documentRangeFormattingProvider,omitempty"`
	DocumentOnTypeFormattingProvider *DocumentOnTypeFormattingOptions `json:"documentOnTypeFormattingProvider,omitempty"`
	RenameProvider                   *bool                            `json:"renameProvider,omitempty"`
	FoldingRangeProvider             *bool                            `json:"foldingRangeProvider,omitempty"`
	ExecuteCommandProvider           *ExecuteCommandOptions           `json:"executeCommandProvider,omitempty"`
	SelectionRangeProvider           *bool                            `json:"selectionRangeProvider,omitempty"`
	SemanticTokensProvider           *SemanticTokensOptions           `json:"semanticTokensProvider,omitempty"`
	InlayValueProvider               *InlayValueOptions               `json:"inlayValueProvider,omitempty"`
	InlayHintProvider                *InlayHintOptions                `json:"inlayHintProvider,omitempty"`
	DiagnosticProvider               *DiagnosticOptions               `json:"diagnosticProvider,omitempty"`
	WorkspaceSymbolProvider          *bool                            `json:"workspaceSymbolProvider,omitempty"`
}

type ExecuteCommandOptions struct {
	Commands []string `json:"commands"`
}
