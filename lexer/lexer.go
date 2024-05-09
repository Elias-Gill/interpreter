package lexer

import "github.com/sl2.0/tokens"

type Lexer struct {
	input string
	// both variables are initialized on 0 by default (go's default behavior).
	currentPosition int // position of the current character
	nextPosition    int // position of the next character
	ch              byte
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

	// first search for comments and ignore them, consuming every
	// character till the end of the line (or end of the file)
	if l.ch == '/' && l.pickChar() == '/' {
		for l.ch != '\n' && l.ch != 0 {
			l.readChar()
		}
		l.skipLineBreaks()
		l.burnWhiteSpaces()
	}

	// start generating tokens
	switch l.ch {
	// operators
	case '-':
		token = newSingleToken(tokens.MINUS, l.ch)
	case '+':
		token = newSingleToken(tokens.PLUS, l.ch)
	case '*':
		token = newSingleToken(tokens.ASTERISC, l.ch)
	case '/':
		token = newSingleToken(tokens.SLASH, l.ch)
	case '<':
		token = newSingleToken(tokens.LT, l.ch)
	case '>':
		token = newSingleToken(tokens.GT, l.ch)
	case '!':
		ch := l.pickChar()
		if ch == '=' {
			token = newMultiToken(tokens.NOTEQUAL, "!=")
			l.readChar()
		} else {
			token = newSingleToken(tokens.BANG, '!')
		}
	case '=':
		ch := l.pickChar()
		if ch == '=' {
			token = newMultiToken(tokens.EQUALS, "==")
			l.readChar()
		} else {
			token = newSingleToken(tokens.ASIGN, '=')
		}

		// especial chars
	case ',':
		token = newSingleToken(tokens.COMMA, l.ch)
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
	case '"':
		str := ""
		for l.pickChar() != '"' {
			l.readChar()
			str += string(l.ch)
		}
		// skeep the final '"'
		l.readChar()
		token = newMultiToken(tokens.STRING, str)
	case '\n':
		l.skipLineBreaks()
		// early return to avoid errors with some multiline characters
		return newMultiToken(tokens.LINEBREAK, "")
	case 0:
		token = newMultiToken(tokens.EOF, "")

		// keywords and identifiers (aka, multi-char tokens)
	default:
		if isLetter(l.ch) {
			ident := l.extractIdentifier()
			// early return to prevent reading (and skipping) the next char
			return newMultiToken(tokens.ResolveType(ident), ident)
		}

		if isNumber(l.ch) {
			return newMultiToken(tokens.NUMBER, l.extractNumber())
		}

		token = newSingleToken(tokens.ILLEGAL, l.ch)
	}

	l.readChar()

	return token
}
