package ast

import (
	"bytes"
	"strconv"

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

/*
An expression statement is a expression which is not assosiated to a variable
declaration like: -(5+5)
*/
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

type Integer struct {
	Value int64
	Token tokens.Token
}

func NewInteger(t tokens.Token) *Integer {
	value, err := strconv.ParseInt(t.Literal, 0, 64)
	if err != nil {
		return nil
	}

	return &Integer{
		Value: value,
		Token: t,
	}
}
func (i *Integer) expressionNode() {}
func (i *Integer) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Integer) ToString() string {
	return i.TokenLiteral()
}

type PrefixExpression struct {
	Right    Expression
	Operator string
	Token    tokens.Token
}

func (p *PrefixExpression) expressionNode() {}
func (p *PrefixExpression) TokenLiteral() string {
	return p.Token.Literal
}
func (p *PrefixExpression) ToString() string {
	var out bytes.Buffer

	out.WriteString(p.TokenLiteral() + " (")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.ToString())
	out.WriteString(") ")

	return out.String()
}

type InfixExpression struct {
	Right    Expression
	Left     Expression
	Operator string
	Token    tokens.Token
}

func (i *InfixExpression) expressionNode() {}
func (i *InfixExpression) TokenLiteral() string {
	return i.Token.Literal
}
func (i *InfixExpression) ToString() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(i.Left.ToString())
	out.WriteString(i.Operator)
	out.WriteString(i.Right.ToString())
	out.WriteString(")")

	return out.String()
}
