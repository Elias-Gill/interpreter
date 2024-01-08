package lexer

import (
	"testing"

	"github.com/sl2.0/tokens"
)

type testCase struct {
	input    string
	expected []tokens.Token
}

func TestNextToken(t *testing.T) {
	testCases := []testCase{
		{
			`+-:={}()`,
			[]tokens.Token{
				{Type: tokens.PLUS, Literal: "+"},
				{Type: tokens.MINUS, Literal: "-"},
				{Type: tokens.COLON, Literal: ":"},
				{Type: tokens.ASIGN, Literal: "="},
				{Type: tokens.LBRAC, Literal: "{"},
				{Type: tokens.RBRAC, Literal: "}"},
				{Type: tokens.LPAR, Literal: "("},
				{Type: tokens.RPAR, Literal: ")"},
			},
		},

		{
			`var nuevo: entero = 22;`,
			[]tokens.Token{
				{Type: tokens.VAR, Literal: "var"},
				{Type: tokens.IDENT, Literal: "nuevo"},
				{Type: tokens.COLON, Literal: ":"},
				{Type: tokens.INTEGER, Literal: "entero"},
				{Type: tokens.ASIGN, Literal: "="},
				{Type: tokens.INTEGER, Literal: "22"},
                {Type: tokens.SEMICOLON, Literal: ";"},
				{Type: tokens.EOF, Literal: ""},
			},
		},
	}

	for id, test := range testCases {
		lexer := NewLexer(test.input)

		for i := 0; i < len(test.expected); i++ {
			token := lexer.NexToken()
			if token.Type != test.expected[i].Type || token.Literal != test.expected[i].Literal {
				t.Fatalf("\nExpected token type to be: %s \n\tGot: %s \nExpected token value to be: %s \n\tGot: %s \nTest case: %d", test.expected[i].Type, token.Type,
					test.expected[i].Literal, token.Literal, id+1)
			}
		}
	}
}
