package protocol

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ServerCapabilities struct {
	TextDocumentSync                 TextDocumentSyncKind             `json:"textDocumentSync,omitempty"`
	CompletionProvider               *CompletionOptions               `json:"completionProvider"`
	HoverProvider                    bool                             `json:"hoverProvider"`
	SignatureHelpProvider            *SignatureHelpOptions            `json:"signatureHelpProvider,omitempty"`
	DeclarationProvider              bool                             `json:"declarationProvider"`
	DefinitionProvider               bool                             `json:"definitionProvider"`
	TypeDefinitionProvider           bool                             `json:"typeDefinitionProvider"`
	ImplementationProvider           bool                             `json:"implementationProvider"`
	ReferencesProvider               bool                             `json:"referencesProvider"`
	DocumentHighlightProvider        bool                             `json:"documentHighlightProvider"`
	DocumentSymbolProvider           bool                             `json:"documentSymbolProvider"`
	CodeActionProvider               bool                             `json:"codeActionProvider"`
	CodeLensProvider                 *CodeLensOptions                 `json:"codeLensProvider,omitempty"`
	DocumentLinkProvider             *DocumentLinkOptions             `json:"documentLinkProvider,omitempty"`
	ColorProvider                    bool                             `json:"colorProvider"`
	DocumentFormattingProvider       bool                             `json:"documentFormattingProvider"`
	DocumentRangeFormattingProvider  bool                             `json:"documentRangeFormattingProvider"`
	DocumentOnTypeFormattingProvider *DocumentOnTypeFormattingOptions `json:"documentOnTypeFormattingProvider,omitempty"`
	RenameProvider                   bool                             `json:"renameProvider"`
	FoldingRangeProvider             bool                             `json:"foldingRangeProvider"`
	ExecuteCommandProvider           *ExecuteCommandOptions           `json:"executeCommandProvider,omitempty"`
	SelectionRangeProvider           bool                             `json:"selectionRangeProvider"`
	SemanticTokensProvider           *SemanticTokensOptions           `json:"semanticTokensProvider,omitempty"`
	InlayValueProvider               *InlayValueOptions               `json:"inlayValueProvider,omitempty"`
	InlayHintProvider                *InlayHintOptions                `json:"inlayHintProvider,omitempty"`
	DiagnosticProvider               *DiagnosticOptions               `json:"diagnosticProvider,omitempty"`
	WorkspaceSymbolProvider          bool                             `json:"workspaceSymbolProvider"`
}



type ExecuteCommandOptions struct {
	Commands []string `json:"commands"`
}
