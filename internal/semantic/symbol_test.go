package semantic_test

import (
	"testing"

	"github.com/Norgate-AV/netlinx-language-server/internal/semantic"
	"github.com/Norgate-AV/netlinx-language-server/parser"

	"github.com/stretchr/testify/assert"
)

const code = `
PROGRAM_NAME='Test'

DEFINE_DEVICE
dvTP = 10001:1:0

DEFINE_CONSTANT
MAX_USERS = 10

DEFINE_TYPE
struct User {
	char name[50]
	integer age
}

DEFINE_VARIABLE
INTEGER userCount
volatile User users[MAX_USERS]

DEFINE_FUNCTION integer GetUserCount() {
	return userCount
}

define_function integer GetUserAge(User user) {
	return user.age
}

DEFINE_START
users[1].name = "Alice"
users[1].age = 30
userCount++

users[2].name = "Bob"
users[2].age = 25
userCount++

users[3].name = "Charlie"
users[3].age = 35
userCount++
`

func TestCollectSymbols(t *testing.T) {
	// Arrange
	ts, err := parser.NewTreeSitter()
	assert.NoError(t, err)
	defer ts.Close()

	// Act
	tree := ts.Parser.Parse([]byte(code), nil)

	symbolsTable := semantic.GetSymbolTable(tree.RootNode(), []byte(code))

	for name, symbol := range symbolsTable.Symbols {
		t.Logf("Symbol: %s = %v", name, symbol)
	}

	// Assert
	assert.Contains(t, symbolsTable.Symbols, "dvTP")
	assert.Contains(t, symbolsTable.Symbols, "MAX_USERS")
	// assert.Contains(t, symbolsTable, "userCount")
	// assert.Contains(t, symbolsTable, "users")
	// assert.Contains(t, symbolsTable, "GetUserCount")
	// assert.Contains(t, symbolsTable, "GetUserAge")
}
