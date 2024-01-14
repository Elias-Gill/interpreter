package ast

import (
	"bytes"

	"github.com/sl2.0/tokens"
)

type Node interface {
    // returns the token literal of the current token
	TokenLiteral() string

    // returns a string representation of the statements in the ast
	ToString() string
}

type Expression interface {
	Node
	expressionNode()
}

type Statement interface {
	Node
	statementNode()
}

type Ast struct {
	Statements []Statement
}

// writte the string representation of every statemen in a buffer and then
// return it
func (a *Ast) ToString() string {
	var buffer bytes.Buffer

	for _, stmt := range a.Statements {
		buffer.WriteString(stmt.ToString())
	}

	return buffer.String()
}

// --------
// EVERY "declaration" is an statement. Example: a variable declaration,
// and if statement, a function declration.

// --- variable declaration statement ---
type VarStatement struct {
	Ident *Identifier
	Value Expression
	Token tokens.Token
}

func (v *VarStatement) statementNode() {}
func (v *VarStatement) TokenLiteral() string {
	return v.Token.Literal
}
func (v *VarStatement) ToString() string {
	var out bytes.Buffer

	out.WriteString(v.TokenLiteral() + " ")
	out.WriteString(v.Ident.ToString())
	out.WriteString(" = ")

	if v.Value != nil {
		out.WriteString(v.Value.ToString())
	}

	out.WriteString(";")

	return out.String()
}

// --- return statement ---
type ReturnStatement struct {
	ReturnValue Expression
	Token       tokens.Token
}

func (v *ReturnStatement) statementNode() {}
func (v *ReturnStatement) TokenLiteral() string {
	return v.Token.Literal
}
func (r *ReturnStatement) ToString() string {
	var out bytes.Buffer

	out.WriteString(r.TokenLiteral() + " ")

	if r.ReturnValue != nil {
		out.WriteString(r.ReturnValue.ToString())
	}

	out.WriteString(";")

	return out.String()
}

// --- expression statement ---
type ExpressionStatement struct {
	Expression Expression
	Token      tokens.Token
}

func (v *ExpressionStatement) statementNode() {}
func (v *ExpressionStatement) TokenLiteral() string {
	return v.Token.Literal
}
func (e *ExpressionStatement) ToString() string {
	if e.Expression != nil {
		return e.Expression.ToString()
	}

	return ""
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
func (i *Identifier) ToString() string {
	return i.Value
}
