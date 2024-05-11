package parser

import (
	"fmt"

	"github.com/sl2.0/ast"
	"github.com/sl2.0/tokens"
)

// ------------------------------
// -- Prefix parsing functions --
// ------------------------------

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

func (p *Parser) parseString() ast.Expression {
	return ast.NewString(p.currentToken)
}

func (p *Parser) parseBoolExpression() ast.Expression {
	exp := ast.NewBoolean(p.currentToken)

	if exp == nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.currentToken.Literal)
		p.errors = append(p.errors, msg)
	}

	return exp
}

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

func (p *Parser) parseAnonnymousFunction() ast.Expression {
	f := ast.NewAnonymousFunction(p.currentToken)

	params := p.parseFuncParameters()

	f.Paramenters = params

	body := p.parseBlockStatement()
	if body == nil {
		return nil
	}

	f.Body = body

	return f
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.advanceToken()

	exp := p.parseExpression(LOWEST)

	if !p.advanceIfNextToken(tokens.RPAR) {
		return nil
	}

	return exp
}

// -----------------------------
// -- Infix parsing functions --
// -----------------------------

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

func (p *Parser) parseCall(e ast.Expression) ast.Expression {
	f := ast.NewFunctionCall(p.currentToken, e)
	f.Arguments = p.parseCallArguments()
	return f
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	// empty arguments
	if p.nextTokenIs(tokens.RPAR) {
		p.advanceToken()
		return args
	}

	p.advanceToken()

	args = append(args, p.parseExpression(LOWEST))

	for p.nextTokenIs(tokens.COMMA) {
		// jump comma and place on next ident
		p.advanceToken()
		p.advanceToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.advanceIfNextToken(tokens.RPAR) {
		return nil
	}

	return args
}

func (p *Parser) parseForLoop() ast.Expression {
	exp := ast.NewForLoop(p.currentToken)

	if !p.advanceIfNextToken(tokens.NUMBER) {
		p.errors = append(p.errors, "Missing 'iterations' on for loop")
		return nil
	}

	exp.Iterations = *ast.NewInteger(p.currentToken)

	if !p.advanceIfNextToken(tokens.LBRAC) {
		p.errors = append(p.errors, "Missing opening '{' on for loop body")
		return nil
	}

	exp.Body = p.parseBlockStatement()

	return exp
}
