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
	if l.nextPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.nextPosition]
	}

	l.currentPosition = l.nextPosition
	l.nextPosition++
}

// reads a new character WITHOUT changing the lexer state an returns the caracter
func (l Lexer) pickChar() byte {
	if l.nextPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.nextPosition]
	}
}

func isLetter(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}

func (l *Lexer) extractIdentifier() string {
	auxPos := l.currentPosition

	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[auxPos:l.currentPosition]
}

func isNumber(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func (l *Lexer) extractNumber() string {
	auxPos := l.currentPosition

	for isNumber(l.ch) {
		l.readChar()
	}

	return l.input[auxPos:l.currentPosition]
}

func (l *Lexer) burnWhiteSpaces() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipLineBreaks() {
	for l.ch == '\n' {
		l.readChar()
	}
}
