package parser

import (
	"fmt"
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
