package parser

import (
	"fmt"

	"github.com/sl2.0/ast"
	"github.com/sl2.0/tokens"
)

// ---------------------------------
// ----- Parsing expressions -------
// ---------------------------------

// First parse the prefix side of the expression (identifiers, numbers and unary operators),
// then parse the infix part of the expression if exists
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currentToken.Type]

	if prefix == nil {
		p.errors = append(p.errors, "Not prefixFn found for: "+p.currentToken.Literal)
		return nil
	}

	exp := prefix()

	for !p.nextTokenIs(tokens.SEMICOLON) && precedence < p.nextPrecendence() {
		infix := p.infixParseFns[p.nextToken.Type]

		if infix == nil {
			return exp
		}

		p.advanceToken()

		// parse and create an infix expression adding the current prefix expression to it
		exp = infix(exp)
	}

	return exp
}

// Parses prefix expressions (like -X or !X)
func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}

	p.advanceToken()

	exp.Right = p.parseExpression(PREFIX)

	return exp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return ast.NewIdentifier(p.currentToken)
}

func (p *Parser) parseNumber() ast.Expression {
	exp := ast.NewInteger(p.currentToken)

	if exp == nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.currentToken.Literal)
		p.errors = append(p.errors, msg)
	}

	return exp
}

func (p *Parser) parseInfixExpression(e ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{
		Left:     e,
		Operator: p.currentToken.Literal,
		Token:    p.currentToken,
	}

	precedence := p.curPrecendence()

	p.advanceToken()

	exp.Right = p.parseExpression(precedence)

	return exp
}

func (p *Parser) parseBoolExpression() ast.Expression {
	exp := ast.NewBoolean(p.currentToken)

	if exp == nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.currentToken.Literal)
		p.errors = append(p.errors, msg)
	}

	return exp
}

func (p *Parser) parseGroupedExpression() ast.Expression { return nil }

func (p *Parser) parseIfExpression() ast.Expression {
	exp := ast.NewIfExpression(p.currentToken)

	if !p.advanceIfNextToken(tokens.LPAR) {
		p.errors = append(p.errors, "Missing '(' after if expression")
		return nil
	}

	p.advanceToken()

	condition := p.parseExpression(LOWEST)

	if condition == nil {
		return nil
	}

	exp.Condition = condition

	if !p.advanceIfNextToken(tokens.RPAR) {
		p.errors = append(p.errors, "Missing ')' on if expression")
		return nil
	}

	if !p.advanceIfNextToken(tokens.LBRAC) {
		p.errors = append(p.errors, "Missing '{' on if expression")
		return nil
	}

	exp.Consequence = p.parseBlockStatement()

	// if not "else" block, return
	if !p.curTokenIs(tokens.ELSE) {
		return exp
	}

	if !p.advanceIfNextToken(tokens.LBRAC) {
		return nil
	}

	exp.Alternative = p.parseBlockStatement()

	return exp
}

func (p *Parser) parseFunctionExpression() ast.Expression {
	f := ast.NewFunctionLiteral(p.currentToken)

	if !p.advanceIfNextToken(tokens.LPAR) {
		return nil
	}

	params := p.parseFuncParameters()
	if params == nil {
		return nil
	}

	f.Paramenters = params

	body := p.parseBlockStatement()
	if body == nil {
		return nil
	}

	f.Body = body

	return f
}
