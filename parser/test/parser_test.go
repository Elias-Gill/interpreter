package test

import (
	"strings"
	"testing"

	"github.com/sl2.0/ast"
)

func TestFuncCall(t *testing.T) {
	testCases := []struct {
		input string
		args  []string
	}{
		{
			input: `new_function(x, y + 1)`,
			args: []string{
				`x`,
				`infix expression:
 left:
    y
 operator: +
 right:
    1`,
			},
		},
		{
			input: `new_function(x * (4 + 33), y + 1)`,
			args: []string{
				`infix expression:
 left:
    x
 operator: *
 right:
    infix expression:
     left:
        4
     operator: +
     right:
        33`,
				`infix expression:
 left:
    y
 operator: +
 right:
    1`,
			},
		},
	}

	for _, tc := range testCases {
		p := generateProgram(t, tc.input)

		if p == nil {
			t.Fatalf("ParseProgram() returned nil")
		}

		if len(p.Statements) != 1 {
			t.Fatalf("Number of statements found: %d", len(p.Statements))
		}

		stmt, ok := p.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Cannot convert statement to ast.ExpressionStatement")
		}

		exp, ok := stmt.Expression.(*ast.FunctionCall)
		if !ok {
			t.Fatalf("Cannot convert statement to ast.FunctionCall")
		}

		if exp.Identifier.ToString(0) != "new_function" {
			t.Fatalf("Expected identifier 'new_function'. Got %v", exp.Identifier.ToString(0))
		}

		if len(exp.Arguments) != 2 {
			t.Fatalf("Expected 2 arguments. Got %v", len(exp.Arguments))
		}

		for i, v := range exp.Arguments {
			expected := strings.TrimSpace(tc.args[i])
			actual := strings.TrimSpace(v.ToString(0))
			if actual != expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", expected, actual)
			}
		}
	}
}
