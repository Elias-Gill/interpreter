package parser

import (
	"fmt"

	"github.com/sl2.0/tokens"
)

func (p *Parser) advanceToken() {
	p.currentToken = p.nextToken
	p.nextToken = p.lexer.NexToken()
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

// Compares the current token type with the expected type.
func (p *Parser) curTokenIs(expTy tokens.TokenType) bool {
	return p.currentToken.Type == expTy
}

// Compares the next token type with the expected type.
func (p *Parser) nextTokenIs(expTy tokens.TokenType) bool {
	return p.nextToken.Type == expTy
}

// Compares the current token type with the expected type and generates a parsing error if
// false.
func (p *Parser) expectCurToken(expTy tokens.TokenType) bool {
	if p.currentToken.Type == expTy {
		return true
	}

	msg := fmt.Sprintf("Expected '%s'. Got %s", expTy, p.currentToken.Literal)
	p.errors = append(p.errors, msg)

	return false
}

func (p *Parser) expectNextToken(expTy tokens.TokenType) bool {
	if p.nextToken.Type == expTy {
		return true
	}

	msg := fmt.Sprintf("Expected '%s'. Got %s", expTy, p.nextToken.Type)
	p.errors = append(p.errors, msg)

	return false
}

// Advances to the next token if equals the given token type. Generates a parsing error if
// false.
func (p *Parser) advanceIfNextToken(expTy tokens.TokenType) bool {
	if p.nextToken.Type == expTy {
		p.advanceToken()
		return true
	}

	msg := fmt.Sprintf("Expected '%s'. Got %s", expTy, p.nextToken.Type)
	p.errors = append(p.errors, msg)

	return false
}

func (p *Parser) advanceIfCurToken(expTy tokens.TokenType) bool {
	if p.currentToken.Type == expTy {
		p.advanceToken()
		return true
	}

	msg := fmt.Sprintf("Expected '%s'. Got %s", expTy, p.currentToken.Literal)
	p.errors = append(p.errors, msg)

	return false
}
