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

	// initializes the lexer in a full working state
	l.readChar()

	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func isLetter(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}

func (l *Lexer) extractIdentifier() string {
	auxPos := l.position

	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[auxPos:l.position]
}

func isNumber(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func (l *Lexer) extractNumber() string {
	auxPos := l.position

	for isNumber(l.ch) {
		l.readChar()
	}

	return l.input[auxPos:l.position]
}

func (l *Lexer) burnWhiteSpaces() {
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) NexToken() tokens.Token {
	var token tokens.Token

	l.burnWhiteSpaces()

	// start generating tokens
	switch l.ch {
	case '-':
		token = newToken(tokens.MINUS, l.ch)
	case '+':
		token = newToken(tokens.PLUS, l.ch)
	case '=':
		token = newToken(tokens.ASIGN, l.ch)
	case ';':
		token = newToken(tokens.SEMICOLON, l.ch)
	case ':':
		token = newToken(tokens.COLON, l.ch)
	case '{':
		token = newToken(tokens.LBRAC, l.ch)
	case '}':
		token = newToken(tokens.RBRAC, l.ch)
	case ')':
		token = newToken(tokens.RPAR, l.ch)
	case '(':
		token = newToken(tokens.LPAR, l.ch)
	case 0:
		token.Literal = ""
		token.Type = tokens.EOF
	default: // the default case are keywords and identifiers
		if isLetter(l.ch) {
			ident := l.extractIdentifier()
			token.Type = tokens.TokenizeIdent(ident)
			token.Literal = ident

			// early return to prevent the next readChar or it will jump
			// one char more than expected
			return token
		} else if isNumber(l.ch) {
            number := l.extractNumber()
            token.Literal = number
            token.Type = tokens.INTEGER

            return token
		} else {
			return newToken(tokens.ILLEGAL, l.ch)
		}
	}

	l.readChar()

	return token
}

func newToken(ty tokens.TokenType, ch byte) tokens.Token {
	return tokens.Token{
		Type:    ty,
		Literal: string(ch),
	}
}
