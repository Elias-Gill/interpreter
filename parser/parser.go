package parser

import (
	"fmt"

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
	PROD      // * /
	PREFIX    // -X  !X
	CALL      // foo(bar)
)

var precedences = map[string]int{
	tokens.EQUALS:   EQUALS,
	tokens.NOTEQUAL: EQUALS,
	tokens.LT:       GREATLESS,
	tokens.GT:       GREATLESS,
	tokens.PLUS:     SUM,
	tokens.MINUS:    SUM,
	tokens.ASTERISC: PROD,
	tokens.SLASH:    PROD,
	tokens.FUNCTION: CALL,
	tokens.LPAR:     CALL,
}

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

	parser.registerPrefixFn(tokens.IDENT, parser.parseIdentifier)
	parser.registerPrefixFn(tokens.NUMBER, parser.parseNumber)
	parser.registerPrefixFn(tokens.LPAR, parser.parseGroupedExpression)
	parser.registerPrefixFn(tokens.BANG, parser.parsePrefixExpression)
	parser.registerPrefixFn(tokens.MINUS, parser.parsePrefixExpression)
	parser.registerPrefixFn(tokens.TRUE, parser.parseBoolExpression)
	parser.registerPrefixFn(tokens.FALSE, parser.parseBoolExpression)

    parser.registerInfixFn(tokens.MINUS, parser.parseInfixExpression)
    parser.registerInfixFn(tokens.PLUS, parser.parseInfixExpression)
    parser.registerInfixFn(tokens.SLASH, parser.parseInfixExpression)
    parser.registerInfixFn(tokens.ASTERISC, parser.parseInfixExpression)
    parser.registerInfixFn(tokens.GT, parser.parseInfixExpression)
    parser.registerInfixFn(tokens.LT, parser.parseInfixExpression)
    parser.registerInfixFn(tokens.EQUALS, parser.parseInfixExpression)
    parser.registerInfixFn(tokens.NOTEQUAL, parser.parseInfixExpression)
	// parser.registerInfix(tokens.LPAR, parser.parseCall)

	return parser
}

func (p *Parser) ParseProgram() *ast.Ast {
	tree := &ast.Ast{}
	tree.Statements = []ast.Statement{}

	for !p.curTokenIs(tokens.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			tree.Statements = append(tree.Statements, stmt)
		}

		p.advanceToken()
	}

	return tree
}

// --------------------------------
// ----- Parsing statements -------
// --------------------------------

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case tokens.VAR:
		return p.parseVarStatement()
	case tokens.RETURN:
		return p.parseReturnStatement()
	case tokens.LINEBREAK:
		// Do nothing for now (TODO:)
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

// ---------------------------------
// ----- Parsing expressions -------
// ---------------------------------

// First parse the prefix side of the expression (identifiers, numbers and unary operators),
// then parse the infix part of the expression if exists
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currentToken.Type]

	if prefix == nil {
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

func (p *Parser) parseGroupedExpression() ast.Expression { return nil }

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
