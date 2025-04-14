package lsp

type MessageType = int

const (
	MessageTypeRequest MessageType = iota
	MessageTypeNotification
)

const (
	// General
	MethodInitialize  string = "initialize"
	MethodInitialized string = "initialized"
	MethodShutdown    string = "shutdown"
	MethodExit        string = "exit"

	// Document Management
	MethodTextDocumentDidOpen           string = "textDocument/didOpen"
	MethodTextDocumentDidChange         string = "textDocument/didChange"
	MethodTextDocumentDidClose          string = "textDocument/didClose"
	MethodTextDocumentDidSave           string = "textDocument/didSave"
	MethodTextDocumentWillSave          string = "textDocument/willSave"
	MethodTextDocumentWillSaveWaitUntil string = "textDocument/willSaveWaitUntil"

	// Navigation
	MethodTextDocumentDefinition     string = "textDocument/definition"
	MethodTextDocumentTypeDefinition string = "textDocument/typeDefinition"
	MethodTextDocumentImplementation string = "textDocument/implementation"
	MethodTextDocumentReferences     string = "textDocument/references"
	MethodTextDocumentDeclaration    string = "textDocument/declaration"

	// Information
	MethodTextDocumentHover             string = "textDocument/hover"
	MethodTextDocumentDocumentHighlight string = "textDocument/documentHighlight"
	MethodTextDocumentDocumentSymbol    string = "textDocument/documentSymbol"
	MethodTextDocumentSignatureHelp     string = "textDocument/signatureHelp"

	// Editing Support
	MethodTextDocumentCompletion         string = "textDocument/completion"
	MethodCompletionItemResolve          string = "completionItem/resolve"
	MethodTextDocumentFormatting         string = "textDocument/formatting"
	MethodTextDocumentRangeFormatting    string = "textDocument/rangeFormatting"
	MethodTextDocumentOnTypeFormatting   string = "textDocument/onTypeFormatting"
	MethodTextDocumentRename             string = "textDocument/rename"
	MethodTextDocumentPrepareRename      string = "textDocument/prepareRename"
	MethodTextDocumentCodeAction         string = "textDocument/codeAction"
	MethodTextDocumentCodeLens           string = "textDocument/codeLens"
	MethodCodeLensResolve                string = "codeLens/resolve"
	MethodTextDocumentLinkedEditingRange string = "textDocument/linkedEditingRange"

	// Folding/Structure
	MethodTextDocumentFoldingRange   string = "textDocument/foldingRange"
	MethodTextDocumentSelectionRange string = "textDocument/selectionRange"

	// Advanced Analysis
	MethodTextDocumentDiagnostic           string = "textDocument/diagnostic"
	MethodTextDocumentSemanticTokensFull   string = "textDocument/semanticTokens/full"
	MethodTextDocumentSemanticTokensDelta  string = "textDocument/semanticTokens/full/delta"
	MethodTextDocumentSemanticTokensRange  string = "textDocument/semanticTokens/range"
	MethodTextDocumentCallHierarchyPrepare string = "textDocument/prepareCallHierarchy"
	MethodCallHierarchyIncomingCalls       string = "callHierarchy/incomingCalls"
	MethodCallHierarchyOutgoingCalls       string = "callHierarchy/outgoingCalls"
	MethodTextDocumentMoniker              string = "textDocument/moniker"
	MethodTextDocumentInlineValue          string = "textDocument/inlineValue"

	// Workspace
	MethodWorkspaceSymbol           string = "workspace/symbol"
	MethodWorkspaceWorkspaceFolders string = "workspace/workspaceFolders"
	MethodWorkspaceExecuteCommand   string = "workspace/executeCommand"
	MethodWorkspaceApplyEdit        string = "workspace/applyEdit"

	// Workspace Notifications
	MethodWorkspaceDidChangeConfig           string = "workspace/didChangeConfiguration"
	MethodWorkspaceDidChangeWatchedFiles     string = "workspace/didChangeWatchedFiles"
	MethodWorkspaceDidChangeWorkspaceFolders string = "workspace/didChangeWorkspaceFolders"

	// Window
	MethodWindowShowMessage        string = "window/showMessage"
	MethodWindowShowMessageRequest string = "window/showMessageRequest"
	MethodWindowLogMessage         string = "window/logMessage"

	// Server-to-Client
	MethodTextDocumentPublishDiagnostics string = "textDocument/publishDiagnostics"

	// Special
	MethodCancelRequest              string = "$/cancelRequest"
	MethodProgress                   string = "$/progress"
	MethodClientRegisterCapability   string = "client/registerCapability"
	MethodClientUnregisterCapability string = "client/unregisterCapability"

	// Custom
	MethodNetLinxServerLogPath string = "netlinx/serverLogPath"
)

var MessageTypeMap = map[string]MessageType{
	// General
	MethodInitialize:  MessageTypeRequest,
	MethodInitialized: MessageTypeNotification,
	MethodShutdown:    MessageTypeRequest,
	MethodExit:        MessageTypeNotification,

	// Document Management (all notifications except willSaveWaitUntil)
	MethodTextDocumentDidOpen:           MessageTypeNotification,
	MethodTextDocumentDidChange:         MessageTypeNotification,
	MethodTextDocumentDidClose:          MessageTypeNotification,
	MethodTextDocumentDidSave:           MessageTypeNotification,
	MethodTextDocumentWillSave:          MessageTypeNotification,
	MethodTextDocumentWillSaveWaitUntil: MessageTypeRequest,

	// Navigation (all requests)
	MethodTextDocumentDefinition:     MessageTypeRequest,
	MethodTextDocumentTypeDefinition: MessageTypeRequest,
	MethodTextDocumentImplementation: MessageTypeRequest,
	MethodTextDocumentReferences:     MessageTypeRequest,
	MethodTextDocumentDeclaration:    MessageTypeRequest,

	// Information (all requests)
	MethodTextDocumentHover:             MessageTypeRequest,
	MethodTextDocumentDocumentHighlight: MessageTypeRequest,
	MethodTextDocumentDocumentSymbol:    MessageTypeRequest,
	MethodTextDocumentSignatureHelp:     MessageTypeRequest,

	// Editing Support (all requests)
	MethodTextDocumentCompletion:         MessageTypeRequest,
	MethodCompletionItemResolve:          MessageTypeRequest,
	MethodTextDocumentFormatting:         MessageTypeRequest,
	MethodTextDocumentRangeFormatting:    MessageTypeRequest,
	MethodTextDocumentOnTypeFormatting:   MessageTypeRequest,
	MethodTextDocumentRename:             MessageTypeRequest,
	MethodTextDocumentPrepareRename:      MessageTypeRequest,
	MethodTextDocumentCodeAction:         MessageTypeRequest,
	MethodTextDocumentCodeLens:           MessageTypeRequest,
	MethodCodeLensResolve:                MessageTypeRequest,
	MethodTextDocumentLinkedEditingRange: MessageTypeRequest,

	// Folding/Structure (all requests)
	MethodTextDocumentFoldingRange:   MessageTypeRequest,
	MethodTextDocumentSelectionRange: MessageTypeRequest,

	// Advanced Analysis (all requests)
	MethodTextDocumentDiagnostic:           MessageTypeRequest,
	MethodTextDocumentSemanticTokensFull:   MessageTypeRequest,
	MethodTextDocumentSemanticTokensDelta:  MessageTypeRequest,
	MethodTextDocumentSemanticTokensRange:  MessageTypeRequest,
	MethodTextDocumentCallHierarchyPrepare: MessageTypeRequest,
	MethodCallHierarchyIncomingCalls:       MessageTypeRequest,
	MethodCallHierarchyOutgoingCalls:       MessageTypeRequest,
	MethodTextDocumentMoniker:              MessageTypeRequest,
	MethodTextDocumentInlineValue:          MessageTypeRequest,

	// Workspace (mix of requests and notifications)
	MethodWorkspaceSymbol:           MessageTypeRequest,
	MethodWorkspaceWorkspaceFolders: MessageTypeRequest,
	MethodWorkspaceExecuteCommand:   MessageTypeRequest,
	MethodWorkspaceApplyEdit:        MessageTypeRequest,

	// Workspace Notifications (all notifications)
	MethodWorkspaceDidChangeConfig:           MessageTypeNotification,
	MethodWorkspaceDidChangeWatchedFiles:     MessageTypeNotification,
	MethodWorkspaceDidChangeWorkspaceFolders: MessageTypeNotification,

	// Window (mix)
	MethodWindowShowMessage:        MessageTypeNotification,
	MethodWindowShowMessageRequest: MessageTypeRequest,
	MethodWindowLogMessage:         MessageTypeNotification,

	// Server-to-Client
	MethodTextDocumentPublishDiagnostics: MessageTypeNotification,

	// Special
	MethodCancelRequest:              MessageTypeNotification,
	MethodProgress:                   MessageTypeNotification,
	MethodClientRegisterCapability:   MessageTypeRequest,
	MethodClientUnregisterCapability: MessageTypeRequest,

	// Custom
	MethodNetLinxServerLogPath: MessageTypeRequest,
}

func IsNotification(method string) bool {
	msgType, exists := MessageTypeMap[method]
	if !exists {
		return false
	}

	return msgType == MessageTypeNotification
}
