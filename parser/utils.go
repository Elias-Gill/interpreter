package parser

import (
	"fmt"

	"github.com/sl2.0/tokens"
)

func (p *Parser) advanceToken() {
	p.currentToken = p.nextToken
	p.nextToken = p.lexer.NexToken()
}

// Generates a new parser error with the current token
func (p *Parser) newParserError(msg string) {
	err := fmt.Sprintf("Parse error: %s \n Token: %s", msg, p.currentToken.Type)
	p.errors = append(p.errors, err)
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) registerInfixFn(t tokens.TokenType, f infixFn) {
	p.infixParseFns[t] = f
}

func (p *Parser) registerPrefixFn(t tokens.TokenType, f prefixFn) {
	p.prefixParseFns[t] = f
}

// returns the precedence lvl of the current token
func (p *Parser) curPrecendence() int {
	value, ok := precedences[string(p.currentToken.Type)]
	if !ok {
		return LOWEST
	}

	return value
}

// returns the precedence lvl of the next token
func (p *Parser) nextPrecendence() int {
	value, ok := precedences[string(p.nextToken.Type)]
	if !ok {
		return LOWEST
	}

	return value
}
