package parser

import (
	"fmt"
	"testing"

	"github.com/sl2.0/ast"
)

const colorMagenta = "\033[35m"
const colorGreen = "\033[33m"
const colorNone = "\033[0m"

type testCase struct {
	input     string
	numStatms int
	expected  []string // expected identifiers
}

// --- TESTS ---

func TestVarDeclaration(t *testing.T) {
	testCases := []testCase{
		{
			input: `
        var persona = 21;
        var x = 21;
        var y = 21;
        `,
			numStatms: 3,
			expected:  []string{"persona", "x", "y"},
		},
		{
			input:     ` var persona = 21; `,
			numStatms: 1,
			expected:  []string{"persona"},
		},
	}

	for i, tc := range testCases {
		t.Logf("\n%sRunning test case %d%s", colorMagenta, i, colorNone)

		parser := NewParser(tc.input)
		p := parser.ParseProgram()

		// if the parsing stage has errors, the test will FailNow()
		checkErrors(t, parser)

		if p == nil {
			t.Fatalf("ParseProgram() returned nil")
		}

		if len(p.Statements) != tc.numStatms {
			t.Fatalf("Number of statements must be: %d\n Found: %d",
				tc.numStatms, len(p.Statements),
			)
		}

		// chech if the parsed statemens match the expected test case
		testIndVar(t, p, tc)
	}
}

func TestReturn(t *testing.T) {
	tc := testCase{
		input: `
            return persona;
            return x;
            return y + 21;
            `,
		numStatms: 3,
	}

	parser := NewParser(tc.input)
	p := parser.ParseProgram()

	// if the parsing stage has errors, the test will FailNow()
	checkErrors(t, parser)

	if p == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(p.Statements) != tc.numStatms {
		t.Fatalf("Number of statements must be: %d\n Found: %d",
			tc.numStatms, len(p.Statements),
		)
	}

	for _, value := range p.Statements {
		if value.TokenLiteral() != "return" {
			t.Errorf("Parser error\n \tExpected variable declaration\n\tGot: %v",
				value.TokenLiteral())
			continue
		}

		// try to convert to type VarStatement
		_, ok := value.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("Parser error\n \tCannot convert statement to ast.VarStatement")
		}
	}

}

// --- UTILS for tests ---

func testIndVar(t *testing.T, p *ast.Ast, tc testCase) {
	hasErros := false
	for k, value := range p.Statements {
		// check if the generated statement is a variable declaration
		if value.TokenLiteral() != "var" {
			t.Errorf("Parser error\n \tExpected variable declaration\n\tGot: %v",
				value.TokenLiteral())
			hasErros = true
			continue
		}

		// try to convert to type VarStatement
		varStmt, ok := value.(*ast.VarStatement)
		if !ok {
			hasErros = true
			t.Errorf("Parser error\n \tCannot convert statement to ast.VarStatement")
		}

		// check the variable value
		if varStmt.Ident.Value != tc.expected[k] {
			hasErros = true
			t.Errorf("Parser error\n \tExpected: %v\n\tGot: %v",
				tc.expected[k], value.TokenLiteral())
		}

		// check the variable token
		if varStmt.Ident.TokenLiteral() != tc.expected[k] {
			hasErros = true
			t.Errorf("Parser error\n \tExpected: %v\n\tGot: %v",
				tc.expected[k], value.TokenLiteral())
		}
	}
	if !hasErros {
		t.Logf("\t%s✓ No errors found%s", colorGreen, colorNone)
	}
}

func checkErrors(t *testing.T, program *Parser) {
	if len(program.Errors()) != 0 {
		s := ""
		for _, value := range program.errors {
			s += fmt.Sprintf("\n%s", value)
		}

		t.Fatalf("Found %s%d%s parsing errors: %s",
			colorMagenta, len(program.Errors()), colorNone, s)
	}
}
