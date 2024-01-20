package parser

import (
	"github.com/sl2.0/ast"
	"github.com/sl2.0/tokens"
)

// --------------------------------
// ----- Parsing statements -------
// --------------------------------

func (p *Parser) parseStatement() ast.Statement {
    switch p.currentToken.Type {
    case tokens.VAR:
        return p.parseVarStatement()
    case tokens.RETURN:
        return p.parseReturnStatement()
    case tokens.FUNCTION:
        return p.parseFunctionStatement()
    case tokens.LINEBREAK:
        // Do nothing for now (TODO:)
        return nil
    default:
        return p.parseExpressionStatement()
    }
}

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

    stmt.Ident = ast.NewIdentifier(p.currentToken)

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
    tree := &ast.BlockStatement{}
    tree.Statements = []ast.Statement{}

    // jump '{'
    p.advanceToken()

    for !p.curTokenIs(tokens.EOF) && !p.curTokenIs(tokens.RBRAC) {
        stmt := p.parseStatement()

        if stmt != nil {
            tree.Statements = append(tree.Statements, stmt)
        }

        p.advanceToken()
    }

    return tree
}

func (p *Parser) parseFunctionStatement() *ast.ExpressionStatement {
    // TODO:
    return nil
}

