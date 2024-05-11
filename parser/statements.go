package parser

import (
	"github.com/sl2.0/ast"
	"github.com/sl2.0/tokens"
)

// --------------------------------
// ----- Parsing statements -------
// --------------------------------

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{
		Token: p.currentToken,
	}

	exp := p.parseExpression(LOWEST)
	if exp == nil {
		return nil
	}

	stmt.Expression = exp

	// to support expression with optional semicolon
	if p.nextTokenIs(tokens.SEMICOLON) {
		p.advanceToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{
		Token: p.currentToken,
	}

	// step over "retorna"
	p.advanceToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.nextToken.Type == tokens.SEMICOLON {
		p.advanceToken()
	}

	return stmt
}

func (p *Parser) parseVarStatement() *ast.VarStatement {
	stmt := &ast.VarStatement{
		Token: p.currentToken,
	}

	if !p.advanceIfNextToken(tokens.IDENT) {
		return nil
	}

	stmt.Identifier = ast.NewIdentifier(p.currentToken)

	if !p.advanceIfNextToken(tokens.ASIGN) {
		return nil
	}

	// step over "="
	p.advanceToken()

	stmt.Value = p.parseExpression(LOWEST)

	p.advanceToken()

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{}
	block.Statements = []ast.Statement{}

	if !p.advanceIfCurToken(tokens.LBRAC) {
		return nil
	}

	for !p.curTokenIs(tokens.RBRAC) {
		stmt := p.parseStatement()

		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		if !p.curTokenIs(tokens.RBRAC) {
			p.advanceToken()
		}
	}

	if !p.advanceIfCurToken(tokens.RBRAC) {
		return nil
	}

	return block
}

func (p *Parser) parseFunctionStatement() *ast.FunctionStatement {
	f := ast.NewFunctionStatement(p.currentToken)

	if !p.advanceIfNextToken(tokens.IDENT) {
		return nil
	}

	f.Identifier = ast.NewIdentifier(p.currentToken)

	f.Paramenters = p.parseFuncParameters()

	body := p.parseBlockStatement()
	if body == nil {
		return nil
	}

	f.Body = body

	return f
}

func (p *Parser) parseFuncParameters() []*ast.Identifier {
	var params []*ast.Identifier

	// jump over the "("
	if !p.advanceIfNextToken(tokens.LPAR) {
		return nil
	}
	p.advanceToken()

	for !p.curTokenIs(tokens.RPAR) {
		ident := ast.NewIdentifier(p.currentToken)
		params = append(params, ident)

		p.advanceToken()

		if p.curTokenIs(tokens.COMMA) {
			p.advanceToken()
		}
	}

	if !p.advanceIfCurToken(tokens.RPAR) {
		return nil
	}

	return params
}
