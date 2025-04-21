package captures

// Sections (just starting with the most common ones for now)
const (
	CaptureSectionDefineDevice   = "section.define_device"
	CaptureSectionDefineConstant = "section.define_constant"
	CaptureSectionDefineType     = "section.define_type"
	CaptureSectionDefineVariable = "section.define_variable"
	CaptureSectionDefineStart    = "section.define_start"
	CaptureSectionDefineEvent    = "section.define_event"
)

// Symbols
const (
	CaptureSymbolDeclaration = "symbol.declaration"
	CaptureSymbolStorage     = "symbol.storage"
	CaptureSymbolQualifier   = "symbol.qualifier"
	CaptureSymbolType        = "symbol.type"
	CaptureSymbolIdentifier  = "symbol.identifier"
	CaptureSymbolValue       = "symbol.value"
	CaptureSymbolSize        = "symbol.size"
)

// Functions
const (
	CaptureFunctionDefinition    = "function.definition"
	CaptureFunctionReturnType    = "function.return_type"
	CaptureFunctionName          = "function.name"
	CaptureFunctionParameter     = "function.parameter"
	CaptureFunctionParameterType = "function.parameter.type"
	CaptureFunctionParameterName = "function.parameter.name"
	CaptureFunctionParameterSize = "function.parameter.size"
)

// Types
const (
	CaptureTypeDefinition      = "type.definition"
	CaptureTypeIdentifier      = "type.identifier"
	CaptureTypeBody            = "type.body"
	CaptureTypeField           = "type.field"
	CaptureTypeFieldType       = "type.field.type"
	CaptureTypeFieldIdentifier = "type.field.identifier"
	CaptureTypeFieldValue      = "type.field.value"
	CaptureTypeFieldSize       = "type.field.size"
)
