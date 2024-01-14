package ast

import (
	"testing"

	"github.com/sl2.0/tokens"
)

func TestString(t *testing.T) {
	program := &Ast{
		Statements: []Statement{
			&VarStatement{
				Token: tokens.Token{Type: tokens.VAR, Literal: "let"},
				Ident: &Identifier{
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

	if program.ToString() != "let myVar = anotherVar;" {
		t.Errorf("program.String() wrong. got=%q", program.ToString())
	}
}
