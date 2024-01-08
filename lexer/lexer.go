package lexer

import "github.com/sl2.0/tokens"

type Lexer struct {
	input        string
	position     int // position of the current character
	readPosition int // position of the next character
	ch           byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input: input,
	}

	// initialize the lexer in a full working state
	l.readChar()

	return l
}

func (l *Lexer) NexToken() tokens.Token {
	var token tokens.Token

	l.burnWhiteSpaces()

	// start generating tokens
	switch l.ch {
	case '-':
		token = newSingleToken(tokens.MINUS, l.ch)
	case '+':
		token = newSingleToken(tokens.PLUS, l.ch)
	case '=':
		l.readChar()
		if l.ch == '=' {
			token = newMultiToken(tokens.COMPARE, "==")
		} else {
			return newSingleToken(tokens.ASIGN, '=')
		}
	case ';':
		token = newSingleToken(tokens.SEMICOLON, l.ch)
	case ':':
		token = newSingleToken(tokens.COLON, l.ch)
	case '{':
		token = newSingleToken(tokens.LBRAC, l.ch)
	case '}':
		token = newSingleToken(tokens.RBRAC, l.ch)
	case ')':
		token = newSingleToken(tokens.RPAR, l.ch)
	case '(':
		token = newSingleToken(tokens.LPAR, l.ch)
	case 0:
		token.Literal = ""
		token.Type = tokens.EOF

		// the default case are keywords and identifiers
	default:
		if isLetter(l.ch) {
			ident := l.extractIdentifier()
			// early return to prevent reading (and skipping) the next char
			return newMultiToken(tokens.TokenizeIdent(ident), ident)
		}

		if isNumber(l.ch) {
			return newMultiToken(tokens.NUMBER, l.extractNumber())
		}

		return newSingleToken(tokens.ILLEGAL, l.ch)
	}

	l.readChar()

	return token
}
