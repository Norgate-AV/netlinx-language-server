package lsp

type Diagnostic struct {
	Range              Range                           `json:"range"`
	Severity           *DiagnosticSeverity             `json:"severity,omitempty"`
	Code               int                             `json:"code"`
	CodeDescription    *CodeDescription                `json:"codeDescription,omitempty"`
	Source             *string                         `json:"source,omitempty"`
	Message            string                          `json:"message"`
	Tags               *[]DiagnosticTag                `json:"tags,omitempty"`
	RelatedInformation *[]DiagnosticRelatedInformation `json:"relatedInformation,omitempty"`
	Data               *LSPAny                         `json:"data,omitempty"`
}

type DiagnosticSeverity = int

const (
	DiagnosticSeverityError = iota + 1
	DiagnosticSeverityWarning
	DiagnosticSeverityInformation
	DiagnosticSeverityHint
)

type DiagnosticRelatedInformation struct {
	Location Location `json:"location"`
	Message  string   `json:"message"`
}

type DiagnosticTag = int

const (
	DiagnosticTagUnnecessary = iota + 1
	DiagnosticTagDeprecated
)

type CodeDescription struct {
	Href URI `json:"href"`
}

type DiagnosticOptions struct {
	InterFileDependencies bool `json:"interFileDependencies"`
	WorkspaceDiagnostics  bool `json:"workspaceDiagnostics"`
}
