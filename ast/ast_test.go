package ast

import (
	"strings"
	"testing"

	"github.com/sl2.0/tokens"
)

func TestString(t *testing.T) {
	program := Program{
		Statements: []Statement{
			&VarStatement{
				Token: tokens.Token{Type: tokens.VAR, Literal: "var"},
				Identifier: &Identifier{
					Token: tokens.Token{Type: tokens.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: tokens.Token{Type: tokens.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	expected := `var statement:
  identifier: myVar
  value: 
    anotherVar
`

	expected = strings.TrimSpace(expected)
	actual := strings.TrimSpace(program.ToString(0))
	if actual != expected {
		t.Errorf("program.ToString() wrong.\nExpected:\n%s\nGot:\n%s", expected, actual)
	}
}
