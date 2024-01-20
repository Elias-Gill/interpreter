package test

import (
	"testing"

	"github.com/sl2.0/ast"
)

// --- GENERAL TESTS ---

func TestVarStatement(t *testing.T) {
	testCases := []struct {
		expectedValue interface{}
		identifier    string
		input         string
	}{
		{
			input:         "var persona = 21;",
			expectedValue: 21,
			identifier:    "persona",
		},
		{
			input:         " var est_certo = true; ",
			expectedValue: true,
			identifier:    "est_certo",
		},
		{
			input:         " var es_falso = false; ",
			expectedValue: false,
			identifier:    "es_falso",
		},
	}

	for i, tc := range testCases {
		t.Logf("\n%sRunning test case %d%s", colorMagenta, i+1, colorNone)

		p := generateProgram(t, tc.input)

		if len(p.Statements) != 1 {
			t.Errorf("Number of statements found: %d", len(p.Statements))
			continue
		}

		if p.Statements[0].TokenLiteral() != "var" {
			t.Errorf("Parser error\n \tCannot convert statement to ast.ReturnStatement")
			continue
		}

		v, ok := p.Statements[0].(*ast.VarStatement)
		if !ok {
			t.Errorf("Parser error\n \tCannot convert statement to ast.ReturnStatement")
			continue
		}

		if v.Ident.Value != tc.identifier {
			t.Errorf("Expected value %s. Got: %s", tc.input, v.Ident.Value)
			continue
		}

		testLiteralExpression(t, v.Value, tc.expectedValue)
	}
}

func TestReturn(t *testing.T) {
	tesCases := []struct {
		expectedValue interface{}
		input         string
	}{
		{
			input:         "retorna persona;",
			expectedValue: "persona",
		},
	}

	for _, tc := range tesCases {
		p := generateProgram(t, tc.input)

		if p == nil {
			t.Fatalf("ParseProgram() returned nil")
		}

		if len(p.Statements) != 1 {
			t.Fatalf("Number of statements found: %d", len(p.Statements))
		}

		// try to convert to type ReturnStatement
		stmt := p.Statements[0]

		if stmt.TokenLiteral() != "retorna" {
			t.Errorf("Parser error\n \tExpected return statement\n\tGot: %v",
				stmt.TokenLiteral())
			continue
		}

		ret, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("Parser error\n \tCannot convert statement to ast.ReturnStatement")
		}

		testLiteralExpression(t, ret.ReturnValue, tc.expectedValue)
	}
}

func TestIntegerExpression(t *testing.T) {
	tesCases := []struct {
		expectedValue int
		input         string
	}{
		{
			input:         `4;`,
			expectedValue: 4,
		},
	}

	for _, tc := range tesCases {

		p := generateProgram(t, tc.input)

		if p == nil {
			t.Fatalf("ParseProgram() returned nil")
		}

		if len(p.Statements) != 1 {
			t.Errorf("Number of statements found: %d", len(p.Statements))
			continue
		}

		exp, ok := p.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("Parser error\n \tCannot convert statement to ast.ExpressionStatement")
			continue
		}

		testLiteralExpression(t, exp.Expression, tc.expectedValue)
	}
}

func TestIdentifierExpression(t *testing.T) {
	tesCases := []struct {
		expectedValue string
		input         string
	}{
		{
			input:         ` persona; `,
			expectedValue: "persona",
		},
	}

	for _, tc := range tesCases {

		p := generateProgram(t, tc.input)

		if p == nil {
			t.Fatalf("ParseProgram() returned nil")
		}

		if len(p.Statements) != 1 {
			t.Errorf("Number of statements found: %d", len(p.Statements))
		}

		// try to convert to type Identifier
		exp, ok := p.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("Parser error\n \tCannot convert statement to ast.ExpressionStatement")
			continue
		}

		testLiteralExpression(t, exp.Expression, tc.expectedValue)
	}
}

func TestPrefixExpression(t *testing.T) {
	tesCases := []struct {
		expectedValue string
		input         string
	}{
		{
			input:         ` -3; `,
			expectedValue: "(-3)",
		},
		{
			input:         ` -noviembre; `,
			expectedValue: "(-noviembre)",
		},
		{
			input:         ` !noviembre; `,
			expectedValue: "(!noviembre)",
		},
	}

	for _, tc := range tesCases {

		p := generateProgram(t, tc.input)

		if p == nil {
			t.Fatalf("ParseProgram() returned nil")
		}

		if len(p.Statements) != 1 {
			t.Errorf("Number of statements found: %d", len(p.Statements))
		}

		// try to convert to type expression statement
		exp, ok := p.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("Parser error\n \tCannot convert statement to ast.ExpressionStatement")
			continue
		}

		if tc.expectedValue != exp.ToString() {
			t.Errorf("Expected: %s. Got: %s", tc.expectedValue, exp.ToString())
			continue
		}
	}
}
func TestInfixExpression(t *testing.T) {
	tesCases := []struct {
		expectedValue string
		input         string
	}{
		{
			input:         ` 2+3; `,
			expectedValue: "(2+3)",
		},
		{
			input:         ` 21231 * nada; `,
			expectedValue: "(21231*nada)",
		},
		{
			input:         ` 21231 / nada; `,
			expectedValue: "(21231/nada)",
		},
		{
			input:         ` 21231 - nada; `,
			expectedValue: "(21231-nada)",
		},
	}

	for _, tc := range tesCases {

		p := generateProgram(t, tc.input)

		if p == nil {
			t.Fatalf("ParseProgram() returned nil")
		}

		if len(p.Statements) != 1 {
			t.Errorf("Number of statements found: %d", len(p.Statements))
		}

		// try to convert to type Identifier
		exp, ok := p.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("Parser error\n \tCannot convert statement to ast.ExpressionStatement")
			continue
		}

		if tc.expectedValue != exp.ToString() {
			t.Errorf("Expected: %s. Got: %s", tc.expectedValue, exp.ToString())
			continue
		}
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tesCases := []struct {
		expectedValue string
		input         string
	}{
		{
			input:         ` 2+-3; `,
			expectedValue: "(2+(-3))",
		},
        {
            input:         ` -2+-3; `,
            expectedValue: "((-2)+(-3))",
        },
		{
			input:         ` -2 > 5 + 4*nada == 33; `,
			expectedValue: "(((-2)>(5+(4*nada)))==33)",
		},
	}

	for _, tc := range tesCases {

		p := generateProgram(t, tc.input)

		if p == nil {
			t.Fatalf("ParseProgram() returned nil")
		}

		if len(p.Statements) != 1 {
			t.Errorf("Number of statements found: %d", len(p.Statements))
		}

		// try to convert to type Identifier
		exp, ok := p.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("Parser error\n \tCannot convert statement to ast.ExpressionStatement")
			continue
		}

		if tc.expectedValue != exp.ToString() {
			t.Errorf("Expected: %s. Got: %s", tc.expectedValue, exp.ToString())
			continue
		}
	}
}
	// TODO:
}
