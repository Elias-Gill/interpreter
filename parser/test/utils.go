package test

/*
Here are utility functions for the main tests. The real test suit is on
"parser_test.go"
*/

import (
	"fmt"
	"testing"

	"github.com/sl2.0/ast"
	"github.com/sl2.0/parser"
)

const colorMagenta = "\033[35m"
const colorGreen = "\033[33m"
const colorNone = "\033[0m"

/*
Generate a new parser with the given input and parse the program.
Fatal if parsing errors encountered, or the parsing program returns an empty AST.
Then return the given AST.
*/
func generateProgram(t *testing.T, input string) *ast.Program {
	parser := parser.NewParser(input)
	p := parser.ParseProgram()

	// check errors
	if len(parser.Errors()) != 0 {
		s := ""
		for _, value := range parser.Errors() {
			s += fmt.Sprintf("\n%s", value)
		}

		t.Fatalf("Found %s%d%s parsing errors: %s",
			colorMagenta, len(parser.Errors()), colorNone, s)
	}

	if p == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(p.Statements) == 0 {
		t.Fatalf("Empty AST")
	}

	return p
}

func testLiteralExpression(t *testing.T, expression ast.Expression, expected interface{}) bool {
	if expression == nil {
		t.Errorf("Wtf bro, you submited a nil expression")
	}

	switch expected := expected.(type) {
	case bool:
		return testBoolLiteral(t, expression) == expected
	case int64:
		return testIntLiteral(t, expression, expected)
	case int:
		return testIntLiteral(t, expression, int64(expected))
	case string:
		return testIdentifier(t, expression, expected)
	default:
		t.Errorf("Type of expected result not handled. Got %T", expected)
		return false
	}
}

func testBoolLiteral(t *testing.T, exp ast.Expression) bool {
	value, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("Cannot convert statement to ast.Integer. \n\tGot: %v", exp.ToString(0))
		return false
	}

	if value.Token.Literal == "true" {
		return true
	}

	return false
}

func testIntLiteral(t *testing.T, exp ast.Expression, expected int64) bool {
	number, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("Cannot convert statement to ast.Integer. \n\tGot: %v", exp.ToString(0))
		return false
	}

	if number.Value != expected {
		t.Errorf("Expected value %v. Got: %v", expected, number.Value)
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, expected string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("Cannot convert statement to ast.Integer. \n\tGot: %v", exp.ToString(0))
		return false
	}

	if ident.Value != expected {
		t.Errorf("Expected value %v. Got: %v", expected, ident.Value)
		return false
	}

	return true
}

func testInfix(t *testing.T, exp ast.Expression, expected string) {
	if expected != exp.ToString(0) {
		t.Errorf("Expected: %s. Got: %s", expected, exp.ToString(0))
		return
	}
}

func testVar(t *testing.T, exp ast.Statement, identifier string, value interface{}) {

	if exp.TokenLiteral() != "var" {
		t.Errorf("Cannot convert statement to ast.ReturnStatement")
		return
	}

	v, ok := exp.(*ast.VarStatement)
	if !ok {
		t.Errorf("Expected 'var'. Got: %s", exp.ToString(0))
		return
	}

	if v.Identifier.Value != identifier {
		t.Errorf("Expected value %s. Got: %s", identifier, v.Identifier.Value)
		return
	}

	testLiteralExpression(t, v.Value, value)
}

func testReturnFunc(t *testing.T, stmt ast.Statement, expectedValue interface{}) {
	// try to convert to type ReturnStatement
	if stmt.TokenLiteral() != "retorna" {
		t.Errorf("Expected return statement\n\tGot: %v",
			stmt.TokenLiteral())
	}

	ret, ok := stmt.(*ast.ReturnStatement)
	if !ok {
		t.Errorf("Cannot convert statement to ast.ReturnStatement")
	}

	testLiteralExpression(t, ret.ReturnValue, expectedValue)
}
