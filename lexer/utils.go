package lexer

import "github.com/sl2.0/tokens"

// generates a new "single char token" from the given token type and char
func newSingleToken(ty tokens.TokenType, ch byte) tokens.Token {
	return tokens.Token{
		Type:    ty,
		Literal: string(ch),
	}
}

// generates a new "multi-char token" from the given token type and string literal
func newMultiToken(ty tokens.TokenType, lit string) tokens.Token {
	return tokens.Token{
		Type:    ty,
		Literal: string(lit),
	}
}

// reads a new character and advances the lexer state
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
