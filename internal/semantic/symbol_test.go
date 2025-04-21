package semantic_test

import (
	"fmt"
	"testing"

	"github.com/Norgate-AV/netlinx-language-server/internal/captures"
	"github.com/Norgate-AV/netlinx-language-server/internal/parser"
	"github.com/Norgate-AV/netlinx-language-server/internal/semantic"

	"github.com/stretchr/testify/assert"
)

const code = `PROGRAM_NAME='Test'

DEFINE_DEVICE
constant dev dvTP1 = 10001:1:0
dvTP1 = 10002:1:0
dev dvTP1 = 10003:1:0
volatile dvTP4 = 10004:1:0
dvTP5 = 10005:(first_local_port+100):0

DEFINE_CONSTANT
MAX_USERS = 10
USERS[][50] = {
	'Alice',
	'Bob',
	'Charlie'
}

MOREUSERS[][][50] = {
	{
		'Alice',
		'Bob',
		'Charlie'
	},
	{
		'Alice',
		'Bob',
		'Charlie'
	}
}

DEFINE_TYPE
struct User {
	char name[50]
	integer age
}

DEFINE_VARIABLE
count
count2 = 0
INTEGER userCount
volatile User users[MAX_USERS]

DEFINE_FUNCTION integer GetUserCount() {
	return userCount
}

define_function integer GetUserAge(User user, char foo) {
	return user.age
}
`

func TestCollectSymbols(t *testing.T) {
	// Arrange
	ts, err := parser.NewTreeSitter()
	assert.NoError(t, err)
	defer ts.Close()

	// Act
	tree := ts.Parser.Parse([]byte(code), nil)
	fmt.Println(parser.PrettyPrint(tree, parser.PrettyPrintOptions{ShowRanges: true}))

	symbolTable := semantic.GetSymbolTable(tree, []byte(code))
	semantic.PrintSymbolTable(symbolTable)

	// Assert
	assert.Contains(t, symbolTable.Symbols, "dvTP1")
	assert.Contains(t, symbolTable.Symbols, "dvTP4")
	assert.Contains(t, symbolTable.Symbols, "dvTP5")
	assert.Contains(t, symbolTable.Symbols, "MAX_USERS")
	assert.Contains(t, symbolTable.Symbols, "USERS")
	// assert.Contains(t, symbolsTable.Symbols, "MOREUSERS")
	assert.Contains(t, symbolTable.Symbols, "count")
	assert.Contains(t, symbolTable.Symbols, "count2")
	assert.Contains(t, symbolTable.Symbols, "userCount")
	assert.Contains(t, symbolTable.Symbols, "users")

	assert.Equal(t, 3, len(symbolTable.Symbols["dvTP1"]))
	assert.Equal(t, 1, len(symbolTable.Symbols["dvTP4"]))
	assert.Equal(t, 1, len(symbolTable.Symbols["dvTP5"]))
	assert.Equal(t, 1, len(symbolTable.Symbols["MAX_USERS"]))
	assert.Equal(t, 1, len(symbolTable.Symbols["USERS"]))
	// assert.Equal(t, 1, len(symbolsTable.Symbols["MOREUSERS"]))
	assert.Equal(t, 1, len(symbolTable.Symbols["count"]))
	assert.Equal(t, 1, len(symbolTable.Symbols["count2"]))
	assert.Equal(t, 1, len(symbolTable.Symbols["userCount"]))
	assert.Equal(t, 1, len(symbolTable.Symbols["users"]))

	for _, s := range symbolTable.Symbols["dvTP4"] {
		assert.Equal(t, "dvTP4", s.Name)
		assert.Equal(t, captures.CaptureSectionDefineDevice, s.Section)
		assert.Equal(t, semantic.SymbolKindDevice, s.Kind)
		assert.Equal(t, semantic.SymbolQualifierVolatile, s.Qualifier)
		assert.Equal(t, semantic.SymbolDataTypeDev, s.DataType)
		assert.Equal(t, "10004:1:0", s.Value)
		assert.Equal(t, semantic.SizeOfDataType(semantic.SymbolDataTypeDev), s.Size)
		assert.Equal(t, uint(0), s.Dimensions)
	}
	// if s, ok := symbolTable.Symbols["dvTP4"]; ok {
	// 	assert.Equal(t, "dvTP4", s[0].Name)
	// 	assert.Equal(t, semantic.CaptureSectionDefineDevice, s[0].Section)
	// 	assert.Equal(t, semantic.SymbolKindDevice, s[0].Kind)
	// 	assert.Equal(t, semantic.StorageTypeVolatile, s[0].StorageType)
	// 	assert.Equal(t, semantic.DataTypeDev, s[0].DataType)
	// 	assert.Equal(t, "10004:1:0", s[0].Value)
	// 	assert.Equal(t, semantic.SizeOfDataType(semantic.DataTypeDev), s[0].Size)
	// 	assert.Equal(t, uint(0), s[0].Dimensions)
	// }

	// if s, ok := symbolsTable.Symbols["MAX_USERS"]; ok {
	// 	assert.Equal(t, "MAX_USERS", s.Name)
	// 	// assert.Equal(t, semantic.SectionDefineConstant, s.Section)
	// 	assert.Equal(t, semantic.SymbolKindConstant, s.Kind)
	// 	assert.Equal(t, semantic.StorageTypeConstant, s.StorageType)
	// 	assert.Equal(t, semantic.DataTypeInteger, s.DataType)
	// 	assert.Equal(t, "10", s.Value)
	// 	assert.Equal(t, semantic.SizeOfDataType(semantic.DataTypeInteger), s.Size)
	// 	assert.Equal(t, uint(0), s.Dimensions)
	// }

	// if s, ok := symbolsTable.Symbols["USERS"]; ok {
	// 	assert.Equal(t, "USERS", s.Name)
	// 	// assert.Equal(t, semantic.SectionDefineConstant, s.Section)
	// 	assert.Equal(t, semantic.SymbolKindConstant, s.Kind)
	// 	assert.Equal(t, semantic.StorageTypeConstant, s.StorageType)
	// 	assert.Equal(t, semantic.DataTypeChar, s.DataType)
	// 	// assert.Equal(t, "{'Alice','Bob','Charlie'}", s.Value)
	// 	// assert.Equal(t, uint(1), s.Size)
	// 	// assert.Equal(t, uint(50), s.Dimensions)
	// 	// assert.Contains(t, s.Value, "Alice")
	// 	// assert.Contains(t, s.Value, "Bob")
	// 	// assert.Contains(t, s.Value, "Charlie")
	// }

	// if s, ok := symbolsTable.Symbols["MOREUSERS"]; ok {
	// 	assert.Equal(t, "MOREUSERS", s.Name)
	// 	// assert.Equal(t, semantic.SectionDefineConstant, s.Section)
	// 	assert.Equal(t, semantic.SymbolKindConstant, s.Kind)
	// 	assert.Equal(t, semantic.StorageTypeConstant, s.StorageType)
	// 	assert.Equal(t, semantic.DataTypeChar, s.DataType)
	// 	// assert.Equal(t, "{'Alice','Bob','Charlie'}", s.Value)
	// 	// assert.Equal(t, uint(1), s.Size)
	// 	// assert.Equal(t, uint(50), s.Dimensions)
	// 	// assert.Contains(t, s.Value, "Alice")
	// 	// assert.Contains(t, s.Value, "Bob")
	// 	// assert.Contains(t, s.Value, "Charlie")
	// }

	// assert.Contains(t, symbolsTable, "GetUserCount")
	// assert.Contains(t, symbolsTable, "GetUserAge")
}
