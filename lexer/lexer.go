package lexer

import "github.com/sl2.0/tokens"

type Lexer struct {
	input           string
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
	case '\n':
		l.skipLineBreaks()
		token = newMultiToken(tokens.LINEBREAK, "")
	case 0:
		token = newMultiToken(tokens.EOF, "")

		// keywords and identifiers (aka, multi-char tokens)
	default:
		if isLetter(l.ch) {
			ident := l.extractIdentifier()
			// early return to prevent reading (and skipping) the next char
			return newMultiToken(tokens.ResolveIdent(ident), ident)
		}

		if isNumber(l.ch) {
			return newMultiToken(tokens.NUMBER, l.extractNumber())
		}

		token = newSingleToken(tokens.ILLEGAL, l.ch)
	}

	l.readChar()

	return token
}
