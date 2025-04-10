package lsp

type SemanticTokens struct {
	ResultID string   `json:"resultId,omitempty"`
	Data     []uint32 `json:"data"`
}

type SemanticTokensParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type SemanticTokensDelta struct {
	ResultID string               `json:"resultId,omitempty"`
	Edits    []SemanticTokensEdit `json:"edits"`
}

type SemanticTokensEdit struct {
	Start       uint32   `json:"start"`
	DeleteCount uint32   `json:"deleteCount"`
	Data        []uint32 `json:"data,omitempty"`
}

type SemanticTokensRangeParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Range        Range                  `json:"range"`
}

type SemanticTokensOptions struct {
	Legend SemanticTokensLegend `json:"legend"`
	Full   bool                 `json:"full,omitempty"`
	Range  bool                 `json:"range,omitempty"`
}

type SemanticTokensLegend struct {
	TokenTypes     []string `json:"tokenTypes"`
	TokenModifiers []string `json:"tokenModifiers"`
}

const (
	// Token types
	TokenTypeNamespace  = "namespace"
	TokenTypeType       = "type"
	TokenTypeClass      = "class"
	TokenTypeEnum       = "enum"
	TokenTypeInterface  = "interface"
	TokenTypeStruct     = "struct"
	TokenTypeFunction   = "function"
	TokenTypeVariable   = "variable"
	TokenTypeParameter  = "parameter"
	TokenTypeProperty   = "property"
	TokenTypeEnumMember = "enumMember"
	TokenTypeEvent      = "event"
	TokenTypeOperator   = "operator"
	TokenTypeKeyword    = "keyword"
	TokenTypeComment    = "comment"
	TokenTypeString     = "string"
	TokenTypeNumber     = "number"
	TokenTypeRegexp     = "regexp"
	TokenTypeDecorator  = "decorator"

	// Token modifiers
	TokenModifierDeclaration    = "declaration"
	TokenModifierDefinition     = "definition"
	TokenModifierReadonly       = "readonly"
	TokenModifierStatic         = "static"
	TokenModifierDeprecated     = "deprecated"
	TokenModifierAbstract       = "abstract"
	TokenModifierAsync          = "async"
	TokenModifierModification   = "modification"
	TokenModifierDocumentation  = "documentation"
	TokenModifierDefaultLibrary = "defaultLibrary"
)
