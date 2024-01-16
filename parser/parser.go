package parser

import (
	"github.com/sl2.0/ast"
	"github.com/sl2.0/lexer"
	"github.com/sl2.0/tokens"
)

type (
	prefixFn func() ast.Expression
	infixFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	lexer  *lexer.Lexer
	errors []string

	currentToken tokens.Token
	nextToken    tokens.Token

	infixParseFns  map[tokens.TokenType]infixFn
	prefixParseFns map[tokens.TokenType]prefixFn
}

const (
	LOWEST    = iota
	EQUALS    // ==
	GREATLESS // < >
	SUM       // + -
	MULT      // * /
	PREFIX    // -X  !X
	CALL      // foo(bar)
)

func NewParser(input string) *Parser {
	parser := &Parser{
		lexer:  lexer.NewLexer(input),
		errors: []string{},

		infixParseFns:  make(map[tokens.TokenType]infixFn),
		prefixParseFns: make(map[tokens.TokenType]prefixFn),
	}

	// to setup the parser in the correct initial state
	parser.advanceToken()
	parser.advanceToken()

	parser.registerPrefix(tokens.IDENT, parser.parseIdentifier)

	// register infix parsing functions

	return parser
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
	case tokens.LINEBREAK:
		// TODO: ver que onda
		return nil
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{
		Token:      p.currentToken,
		Expression: p.parseExpression(LOWEST),
	}

	// to support expression with optional semicolon
	if p.nextToken.Type == tokens.SEMICOLON {
		p.advanceToken()
	}

	return stmt
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

// ------------------------------------------
// ----- Parsing expression functions -------
// ------------------------------------------
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currentToken.Type]

	if prefix == nil {
		return nil
	}

	leftExpr := prefix()

	return leftExpr
}

func (p *Parser) parseGroupedExpression() *ast.Expression { return nil }

func (p *Parser) parseIdentifier() ast.Expression {
	return ast.NewIdentifier(p.currentToken)
}
