package parser_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Norgate-AV/netlinx-language-server/parser"
)

func TestBasicParsing(t *testing.T) {
	ts, err := parser.NewTreeSitter()
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}

	defer ts.Close()

	code := `
		PROGRAM_NAME='Test'
		DEFINE_DEVICE
		dvTP = 10001:1:0
	`

	tree := ts.Parser.Parse([]byte(code), nil)
	defer tree.Close()
	if tree == nil {
		t.Fatal("Expected non-nil parse tree")
	}

	// Check root node exists
	root := tree.RootNode()
	if root == nil {
		t.Fatal("Expected non-nil root node")
	}

	// Basic validation of the parse tree
	if root.ChildCount() < 1 {
		t.Errorf("Expected at least one child node, got %d", root.ChildCount())
	}
}

func TestInvalidSyntax(t *testing.T) {
	ts, err := parser.NewTreeSitter()
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}

	defer ts.Close()

	code := `PROGRAM_NAME='Unterminated`

	tree := ts.Parser.Parse([]byte(code), nil)
	defer tree.Close()
	if tree == nil {
		t.Fatal("Expected non-nil parse tree")
	}

	if !tree.RootNode().HasError() {
		t.Fatal("Expected parse tree to have error node")
	}
}

func TestQueryParsing(t *testing.T) {
	ts, err := parser.NewTreeSitter()
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}

	defer ts.Close()

	sourceCode := []byte(`
		PROGRAM_NAME='Test'

		DEFINE_DEVICE
		dvTP = 10001:1:0

		DEFINE_CONSTANT
		constant integer FOO = 1

		DEFINE_VARIABLE
		volatile integer bar = 2
	`)

	tree := ts.Parser.Parse(sourceCode, nil)
	defer tree.Close()

	query, err := parser.CreateQuery(
		`
		(section) @section
		`,
	)
	if err != nil {
		t.Fatalf("Failed to create query: %v", err)
	}

	defer query.Close()

	qc := parser.CreateQueryCursor()
	defer qc.Close()

	captures := qc.Captures(query, tree.RootNode(), sourceCode)

	for match, index := captures.Next(); match != nil; match, index = captures.Next() {
		fmt.Printf(
			"Capture %d: %s\n",
			index,
			match.Captures[index].Node.Utf8Text(sourceCode),
		)
	}
}

func TestPrettyPrint(t *testing.T) {
	code := `PROGRAM_NAME='Test'

DEFINE_DEVICE
dvTP = 10001:1:0

DEFINE_CONSTANT
MAX_USERS = 10
USERS[][50] = {
	'Alice',
	'Bob',
	'Charlie'
}

DEFINE_TYPE
struct User {
	char name[50]
	integer age
}

DEFINE_VARIABLE
INTEGER userCount
volatile User users[MAX_USERS]

// DEFINE_FUNCTION integer GetUserCount() {
// 	return userCount
// }

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
// userCount++`

	ts, err := parser.NewTreeSitter()
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}

	defer ts.Close()
	tree := ts.Parser.Parse([]byte(code), nil)
	defer tree.Close()
	sexp := parser.PrettyPrint(tree, parser.PrettyPrintOptions{})

	expected := `(source_file
  (program_name
    (program_name_keyword)
    (string_literal
      (string_content)))
  (section
    (define_device_section
      (define_device_keyword)))
  (expression_statement
    (assignment_expression
      left: (identifier)
      right: (device_literal
        device: (number_literal)
        port: (number_literal)
        system: (number_literal))))
  (section
    (define_constant_section
      (define_constant_keyword)))
  (expression_statement
    (assignment_expression
      left: (identifier)
      right: (number_literal)))
  (expression_statement
    (comma_expression
      left: (assignment_expression
        left: (subscript_expression
          argument: (subscript_expression
            argument: (identifier)
            index: (identifier))
          index: (number_literal))
        (ERROR)
        right: (string_literal
          (string_content)))
      right: (comma_expression
        left: (string_literal
          (string_content))
        right: (string_literal
          (string_content)))))
  (ERROR)
  (section
    (define_type_section
      (define_type_keyword)))
  (type_definition
    (struct_specifier
      (struct_keyword)
      name: (type_identifier)
      body: (field_declaration_list
        (field_declaration
          type: (intrinsic_type
            (primitive_type
              (char_keyword)))
          declarator: (array_declarator
            declarator: (field_identifier)
            size: (number_literal)))
        (field_declaration
          type: (intrinsic_type
            (primitive_type
              (integer_keyword)))
          declarator: (field_identifier)))))
  (section
    (define_variable_section
      (define_variable_keyword)))
  (declaration
    type: (intrinsic_type
      (primitive_type
        (integer_keyword)))
    declarator: (identifier))
  (declaration
    (type_qualifier
      (volatile_keyword))
    (type_identifier)
    declarator: (array_declarator
      declarator: (identifier)
      size: (identifier)))
  (comment)
  (comment)
  (comment)
  (comment)
  (comment)
  (comment)
  (comment)
  (comment)
  (comment)
  (comment)
  (comment)
  (comment)
  (comment)
  (comment)
  (comment)
  (comment))`

	if strings.TrimSpace(sexp) != strings.TrimSpace(expected) {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, sexp)
	}
}

func TestPrettyPrintWithRanges(t *testing.T) {
	code := `PROGRAM_NAME='Test'

DEFINE_DEVICE
dvTP = 10001:1:0

DEFINE_CONSTANT
MAX_USERS = 10
USERS[][50] = {
	'Alice',
	'Bob',
	'Charlie'
}

DEFINE_TYPE
struct User {
	char name[50]
	integer age
}

DEFINE_VARIABLE
INTEGER userCount
volatile User users[MAX_USERS]

// DEFINE_FUNCTION integer GetUserCount() {
// 	return userCount
// }

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

	ts, err := parser.NewTreeSitter()
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}

	defer ts.Close()
	tree := ts.Parser.Parse([]byte(code), nil)
	defer tree.Close()
	sexp := parser.PrettyPrint(tree, parser.PrettyPrintOptions{
		ShowRanges: true,
	})

	expected := `(source_file [0, 0] - [43, 0]
  (program_name [0, 0] - [0, 19]
    (program_name_keyword [0, 0] - [0, 12])
    (string_literal [0, 13] - [0, 19]
      (string_content [0, 14] - [0, 18])))
  (section [2, 0] - [2, 13]
    (define_device_section [2, 0] - [2, 13]
      (define_device_keyword [2, 0] - [2, 13])))
  (expression_statement [3, 0] - [3, 16]
    (assignment_expression [3, 0] - [3, 16]
      left: (identifier [3, 0] - [3, 4])
      right: (device_literal [3, 7] - [3, 16]
        device: (number_literal [3, 7] - [3, 12])
        port: (number_literal [3, 13] - [3, 14])
        system: (number_literal [3, 15] - [3, 16]))))
  (section [5, 0] - [5, 15]
    (define_constant_section [5, 0] - [5, 15]
      (define_constant_keyword [5, 0] - [5, 15])))
  (expression_statement [6, 0] - [6, 14]
    (assignment_expression [6, 0] - [6, 14]
      left: (identifier [6, 0] - [6, 9])
      right: (number_literal [6, 12] - [6, 14])))
  (expression_statement [7, 0] - [10, 10]
    (comma_expression [7, 0] - [10, 10]
      left: (assignment_expression [7, 0] - [8, 8]
        left: (subscript_expression [7, 0] - [7, 11]
          argument: (subscript_expression [7, 0] - [7, 7]
            argument: (identifier [7, 0] - [7, 5])
            index: (identifier [7, 6] - [7, 6]))
          index: (number_literal [7, 8] - [7, 10]))
        (ERROR [7, 14] - [7, 15])
        right: (string_literal [8, 1] - [8, 8]
          (string_content [8, 2] - [8, 7])))
      right: (comma_expression [9, 1] - [10, 10]
        left: (string_literal [9, 1] - [9, 6]
          (string_content [9, 2] - [9, 5]))
        right: (string_literal [10, 1] - [10, 10]
          (string_content [10, 2] - [10, 9])))))
  (ERROR [11, 0] - [11, 1])
  (section [13, 0] - [13, 11]
    (define_type_section [13, 0] - [13, 11]
      (define_type_keyword [13, 0] - [13, 11])))
  (type_definition [14, 0] - [17, 1]
    (struct_specifier [14, 0] - [17, 1]
      (struct_keyword [14, 0] - [14, 6])
      name: (type_identifier [14, 7] - [14, 11])
      body: (field_declaration_list [14, 12] - [17, 1]
        (field_declaration [15, 1] - [15, 14]
          type: (intrinsic_type [15, 1] - [15, 5]
            (primitive_type [15, 1] - [15, 5]
              (char_keyword [15, 1] - [15, 5])))
          declarator: (array_declarator [15, 6] - [15, 14]
            declarator: (field_identifier [15, 6] - [15, 10])
            size: (number_literal [15, 11] - [15, 13])))
        (field_declaration [16, 1] - [16, 12]
          type: (intrinsic_type [16, 1] - [16, 8]
            (primitive_type [16, 1] - [16, 8]
              (integer_keyword [16, 1] - [16, 8])))
          declarator: (field_identifier [16, 9] - [16, 12])))))
  (section [19, 0] - [19, 15]
    (define_variable_section [19, 0] - [19, 15]
      (define_variable_keyword [19, 0] - [19, 15])))
  (declaration [20, 0] - [20, 17]
    type: (intrinsic_type [20, 0] - [20, 7]
      (primitive_type [20, 0] - [20, 7]
        (integer_keyword [20, 0] - [20, 7])))
    declarator: (identifier [20, 8] - [20, 17]))
  (declaration [21, 0] - [21, 30]
    (type_qualifier [21, 0] - [21, 8]
      (volatile_keyword [21, 0] - [21, 8]))
    (type_identifier [21, 9] - [21, 13])
    declarator: (array_declarator [21, 14] - [21, 30]
      declarator: (identifier [21, 14] - [21, 19])
      size: (identifier [21, 20] - [21, 29])))
  (comment [23, 0] - [23, 43])
  (comment [24, 0] - [24, 20])
  (comment [25, 0] - [25, 4])
  (comment [27, 0] - [27, 50])
  (comment [28, 0] - [28, 19])
  (comment [29, 0] - [29, 4])
  (comment [31, 0] - [31, 15])
  (comment [32, 0] - [32, 26])
  (comment [33, 0] - [33, 20])
  (comment [34, 0] - [34, 14])
  (comment [36, 0] - [36, 24])
  (comment [37, 0] - [37, 20])
  (comment [38, 0] - [38, 14])
  (comment [40, 0] - [40, 28])
  (comment [41, 0] - [41, 20])
  (comment [42, 0] - [42, 14]))
`

	if strings.TrimSpace(sexp) != strings.TrimSpace(expected) {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, sexp)
	}
}
