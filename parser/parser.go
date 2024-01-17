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

	parser.registerPrefix(tokens.IDENT, parser.parseIdentifier)
	parser.registerPrefix(tokens.NUMBER, parser.parseNumber)
	parser.registerPrefix(tokens.LPAR, parser.parseGroupedExpression)

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

// --------------------------------
// ----- Parsing statements -------
// --------------------------------

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case tokens.VAR:
		return p.parseVarDeclaration()
	case tokens.RETURN:
		return p.parseReturn()
	case tokens.LINEBREAK:
		// Do nothing for now (TODO)
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

	p.advanceToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.nextToken.Type == tokens.SEMICOLON {
		p.advanceToken()
	}

	return stmt
}

func (p *Parser) parseVarDeclaration() *ast.VarStatement {
	stmt := &ast.VarStatement{
		Token: p.currentToken,
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

	p.advanceToken()

	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

// ---------------------------------
// ----- Parsing expressions -------
// ---------------------------------
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currentToken.Type]

	if prefix == nil {
		return nil
	}

	leftExpr := prefix()

	// aca hacer el for y los infijos

	return leftExpr
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
