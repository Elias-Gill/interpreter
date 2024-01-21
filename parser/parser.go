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
	parser.registerPrefixFn(tokens.IF, parser.parseIfExpression)
    parser.registerPrefixFn(tokens.FUNCTION, parser.parseFunctionExpression)

	parser.registerInfixFn(tokens.MINUS, parser.parseInfixExpression)
	parser.registerInfixFn(tokens.PLUS, parser.parseInfixExpression)
	parser.registerInfixFn(tokens.SLASH, parser.parseInfixExpression)
	parser.registerInfixFn(tokens.ASTERISC, parser.parseInfixExpression)
	parser.registerInfixFn(tokens.GT, parser.parseInfixExpression)
	parser.registerInfixFn(tokens.LT, parser.parseInfixExpression)
	parser.registerInfixFn(tokens.EQUALS, parser.parseInfixExpression)
	parser.registerInfixFn(tokens.NOTEQUAL, parser.parseInfixExpression)
	parser.registerInfixFn(tokens.LPAR, parser.parseCall)

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

