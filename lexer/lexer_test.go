package lexer

import (
	"testing"

	"github.com/sl2.0/tokens"
)

func TestNextToken(t *testing.T) {
	type testCase struct {
		input    string
		expected []tokens.Token
	}

	testCases := []testCase{
		{ // especial tokens (kee an eye on the line break)
			`+-:={}()==!*,;<>!=
            =

            `,
			[]tokens.Token{
				{Type: tokens.PLUS, Literal: "+"},
				{Type: tokens.MINUS, Literal: "-"},
				{Type: tokens.COLON, Literal: ":"},
				{Type: tokens.ASIGN, Literal: "="},
				{Type: tokens.LBRAC, Literal: "{"},
				{Type: tokens.RBRAC, Literal: "}"},
				{Type: tokens.LPAR, Literal: "("},
				{Type: tokens.RPAR, Literal: ")"},
				{Type: tokens.EQUALS, Literal: "=="},
				{Type: tokens.BANG, Literal: "!"},
				{Type: tokens.ASTERISC, Literal: "*"},
				{Type: tokens.COMMA, Literal: ","},
				{Type: tokens.SEMICOLON, Literal: ";"},
				{Type: tokens.LT, Literal: "<"},
				{Type: tokens.GT, Literal: ">"},
				{Type: tokens.NOTEQUAL, Literal: "!="},
                {Type: tokens.LINEBREAK, Literal: ""},
                {Type: tokens.ASIGN, Literal: "="},
                {Type: tokens.LINEBREAK, Literal: ""},
                {Type: tokens.EOF, Literal: ""},
			},
		},

		{
			`var nuevo: entero = 22;`,
			[]tokens.Token{
				{Type: tokens.VAR, Literal: "var"},
				{Type: tokens.IDENT, Literal: "nuevo"},
				{Type: tokens.COLON, Literal: ":"},
				{Type: tokens.DATATYPE, Literal: "entero"},
				{Type: tokens.ASIGN, Literal: "="},
				{Type: tokens.NUMBER, Literal: "22"},
				{Type: tokens.SEMICOLON, Literal: ";"},
				{Type: tokens.EOF, Literal: ""},
			},
		},
		{
			`func nuevo(nombre: cadena): entero {
                var auxiliar: entero
                si algo == true {
                    return false
                }

                return nombre
            }`,
			[]tokens.Token{
				{Type: tokens.FUNCTION, Literal: "func"},
				{Type: tokens.IDENT, Literal: "nuevo"},
				{Type: tokens.LPAR, Literal: "("},
				{Type: tokens.IDENT, Literal: "nombre"},
				{Type: tokens.COLON, Literal: ":"},
				{Type: tokens.DATATYPE, Literal: "cadena"},
				{Type: tokens.RPAR, Literal: ")"},
				{Type: tokens.COLON, Literal: ":"},
				{Type: tokens.DATATYPE, Literal: "entero"},

				{Type: tokens.LBRAC, Literal: "{"},
                {Type: tokens.LINEBREAK, Literal: ""},

				{Type: tokens.VAR, Literal: "var"},
				{Type: tokens.IDENT, Literal: "auxiliar"},
				{Type: tokens.COLON, Literal: ":"},
				{Type: tokens.DATATYPE, Literal: "entero"},

                {Type: tokens.LINEBREAK, Literal: ""},

				{Type: tokens.IF, Literal: "si"},
				{Type: tokens.IDENT, Literal: "algo"},
				{Type: tokens.EQUALS, Literal: "=="},
				{Type: tokens.BOOL, Literal: "true"},
				{Type: tokens.LBRAC, Literal: "{"},

                {Type: tokens.LINEBREAK, Literal: ""},

				{Type: tokens.RETURN, Literal: "return"},
				{Type: tokens.BOOL, Literal: "false"},
                {Type: tokens.LINEBREAK, Literal: ""},
				{Type: tokens.RBRAC, Literal: "}"},

                {Type: tokens.LINEBREAK, Literal: ""},

				{Type: tokens.RETURN, Literal: "return"},
				{Type: tokens.IDENT, Literal: "nombre"},

                {Type: tokens.LINEBREAK, Literal: ""},

				{Type: tokens.RBRAC, Literal: "}"},
				{Type: tokens.EOF, Literal: ""},
			},
		},
		{ // invalid tokens
			`~@#$^&`,
			[]tokens.Token{
				{Type: tokens.ILLEGAL, Literal: "~"},
				{Type: tokens.ILLEGAL, Literal: "@"},
				{Type: tokens.ILLEGAL, Literal: "#"},
				{Type: tokens.ILLEGAL, Literal: "$"},
				{Type: tokens.ILLEGAL, Literal: "^"},
				{Type: tokens.ILLEGAL, Literal: "&"},
			},
		},
	}

	for id, test := range testCases {
		lexer := NewLexer(test.input)

		for i := 0; i < len(test.expected); i++ {
			token := lexer.NexToken()
			if token.Type != test.expected[i].Type || token.Literal != test.expected[i].Literal {
				t.Fatalf("\nExpected token type: %s \n\tGot: %s \n\nExpected token value: %s \n\tGot: %s \nTest case: %d, token %d",
					test.expected[i].Type, token.Type,
					test.expected[i].Literal, token.Literal,
					id+1, i,
				)
			}
		}
	}
}
