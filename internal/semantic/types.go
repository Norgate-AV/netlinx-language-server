package semantic

type SymbolStorage = string

const (
	SymbolStorageStackVar SymbolStorage = "stack_var"
	SymbolStorageLocalVar SymbolStorage = "local_var"
)

type SymbolQualifier = string

const (
	SymbolQualifierConstant    SymbolQualifier = "constant"
	SymbolQualifierVolatile    SymbolQualifier = "volatile"
	SymbolQualifierNonVolatile SymbolQualifier = "non_volatile"
	SymbolQualifierPersistent  SymbolQualifier = "persistent"
)

type SymbolKind = string

const (
	SymbolKindDevice   SymbolKind = "device"
	SymbolKindConstant SymbolKind = "constant"
	SymbolKindVariable SymbolKind = "variable"
	SymbolKindFunction SymbolKind = "function"
	SymbolKindStruct   SymbolKind = "struct"
)

type SymbolDataType = string

const (
	SymbolDataTypeChar     SymbolDataType = "char"
	SymbolDataTypeWideChar SymbolDataType = "widechar"
	SymbolDataTypeInteger  SymbolDataType = "integer"
	SymbolDataTypeSinteger SymbolDataType = "sinteger"
	SymbolDataTypeLong     SymbolDataType = "long"
	SymbolDataTypeSlong    SymbolDataType = "slong"
	SymbolDataTypeFloat    SymbolDataType = "float"
	SymbolDataTypeDouble   SymbolDataType = "double"
	SymbolDataTypeDev      SymbolDataType = "dev"
	SymbolDataTypeDevChan  SymbolDataType = "devchan"
	SymbolDataTypeDevLev   SymbolDataType = "devlev"
)

type SymbolScope = string

const (
	SymbolScopeGlobal SymbolScope = "global"
	SymbolScopeLocal  SymbolScope = "local"
)
