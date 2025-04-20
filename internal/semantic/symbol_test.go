package semantic_test

import (
	"fmt"
	"testing"

	"github.com/Norgate-AV/netlinx-language-server/internal/semantic"
	"github.com/Norgate-AV/netlinx-language-server/parser"

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

// DEFINE_TYPE
// struct User {
// 	char name[50]
// 	integer age
// }

DEFINE_VARIABLE
count
count2 = 0
INTEGER userCount
volatile User users[MAX_USERS]

DEFINE_FUNCTION integer GetUserCount() {
	return userCount
}

// define_function integer GetUserAge(User user) {
// 	return user.age
// }

// DEFINE_START
// users[1].name = 'Alice'
// users[1].age = 30
// userCount++

// users[2].name = 'Bob'
// users[2].age = 25
// userCount++

// users[3].name = 'Charlie'
// users[3].age = 35
// userCount++
`

func TestCollectSymbols(t *testing.T) {
	// Arrange
	ts, err := parser.NewTreeSitter()
	assert.NoError(t, err)
	defer ts.Close()

	// Act
	tree := ts.Parser.Parse([]byte(code), nil)

	fmt.Println(parser.PrettyPrint(tree, parser.PrettyPrintOptions{ShowRanges: true}))

	symbolsTable := semantic.GetSymbolTable(tree, []byte(code))

	semantic.PrintSymbolTable(symbolsTable)

	// Assert
	// assert.Contains(t, symbolsTable.Symbols, "dvTP")
	// assert.Contains(t, symbolsTable.Symbols, "MAX_USERS")
	// assert.Contains(t, symbolsTable.Symbols, "USERS")
	// assert.Contains(t, symbolsTable.Symbols, "MOREUSERS")

	// assert.Equal(t, 1, len(symbolsTable.Symbols["dvTP"]))
	// assert.Equal(t, 1, len(symbolsTable.Symbols["MAX_USERS"]))
	// assert.Equal(t, 1, len(symbolsTable.Symbols["USERS"]))
	// assert.Equal(t, 1, len(symbolsTable.Symbols["MOREUSERS"]))

	// if s, ok := symbolsTable.Symbols["dvTP"]; ok {
	// 	assert.Equal(t, "dvTP", s.Name)
	// 	// assert.Equal(t, semantic.SectionDefineDevice, s.Section)
	// 	assert.Equal(t, semantic.SymbolKindDevice, s.Kind)
	// 	assert.Equal(t, semantic.StorageTypeConstant, s.StorageType)
	// 	assert.Equal(t, semantic.DataTypeDev, s.DataType)
	// 	assert.Equal(t, "10001:1:0", s.Value)
	// 	assert.Equal(t, semantic.SizeOfDataType(semantic.DataTypeDev), s.Size)
	// 	assert.Equal(t, uint(0), s.Dimensions)
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

	// assert.Contains(t, symbolsTable.Symbols, "userCount")
	// assert.Contains(t, symbolsTable.Symbols, "users")
	// assert.Contains(t, symbolsTable, "GetUserCount")
	// assert.Contains(t, symbolsTable, "GetUserAge")
}
