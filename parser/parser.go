package parser

import (
	"github.com/sl2.0/ast"
	"github.com/sl2.0/lexer"
	"github.com/sl2.0/tokens"
)

type Parser struct {
	lexer  *lexer.Lexer
	errors []string

	currentToken tokens.Token
	nextToken    tokens.Token
}

func NewParser(input string) *Parser {
	parser := &Parser{
		lexer:  lexer.NewLexer(input),
		errors: []string{},
	}

	// advance tokens two times to setup the parser
	// in the correct initial state (like we did with the lexer)
	parser.advanceToken()
	parser.advanceToken()

	return parser
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) ParseProgram() *ast.Ast {
	tree := &ast.Ast{}
	tree.Statements = []ast.Statement{}

	for p.currentToken.Type != tokens.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			tree.Statements = append(tree.Statements, stmt)
		}
		p.advanceToken()
	}

	return tree
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case tokens.VAR:
		return p.parseVarDeclaration()
	case tokens.RETURN:
		return p.parseReturn()
	default:
		return nil
	}
}

func (p *Parser) parseReturn() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{
		Token: p.currentToken,
	}

	// TODO: parse expressions
	for p.currentToken.Type != tokens.SEMICOLON {
		p.advanceToken()
		if tokens.EOF == p.currentToken.Type {
			p.newParserError(`Expected ";" in return statement`)
			return nil
		}
	}

	return stmt
}

func (p *Parser) parseVarDeclaration() *ast.VarStatement {
	stmt := &ast.VarStatement{
		Token: p.currentToken,
		Value: nil,
	}

	// set the identifier
	p.advanceToken()
	if p.currentToken.Type != tokens.IDENT {
		return nil
	}
	stmt.Ident = ast.NewIdentifier(p.currentToken)

	// search for "="
	p.advanceToken()
	if p.currentToken.Type != tokens.ASIGN {
		p.newParserError(`Expected "=" sign`)
		return nil
	}

	// TODO: parse expressions
	for p.currentToken.Type != tokens.SEMICOLON {
		p.advanceToken()
		// p.parseExpression()
	}

	return stmt
}

func (p *Parser) parseExpression() *ast.Expression {
	var expression *ast.Expression

	p.advanceToken()
	switch p.currentToken.Type {
	case tokens.NUMBER:

	case tokens.IDENT:

	case tokens.SEMICOLON:
		return expression

	default:
		p.newParserError("no se, un error")
		return nil
	}

	return nil
}

func (p *Parser) parseGroupedExpression() *ast.Expression { return nil }
