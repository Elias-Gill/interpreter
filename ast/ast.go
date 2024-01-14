package ast

import "github.com/sl2.0/tokens"

type Ast struct {
	Statements []Statement
}

type Node interface {
	TokenLiteral() string
}

type Expression interface {
	Node
	expressionNode()
}

type Statement interface {
	Node
	statementNode()
}

// --------
// EVERY "declaration" is an statement. Example: a variable declaration,
// and if statement, a function declration.

// the var declaration statement
type VarStatement struct {
	Ident  *Identifier
	Value Expression
	Token tokens.Token
}

func (v *VarStatement) statementNode() {}
func (v *VarStatement) TokenLiteral() string {
	return v.Token.Literal
}

// return statement
type ReturnStatement struct {
    ReturnValue Expression
    Token tokens.Token
}

func (v *ReturnStatement) statementNode() {}
func (v *ReturnStatement) TokenLiteral() string {
    return v.Token.Literal
}

// --------
// THE CONTENT of a statement is called expression. So a statement is a
// TREE OF EXPRESSIONS.

// the identifier is an expression
type Identifier struct {
	Value string
	Token tokens.Token
}

func NewIdentifier(t tokens.Token) *Identifier {
	return &Identifier{
		Value: t.Literal,
		Token: t,
	}
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
