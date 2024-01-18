// THE CONTENT of a statement is called expression. So a statement is a
// TREE OF EXPRESSIONS.

package ast

import (
	"bytes"
	"strconv"

	"github.com/sl2.0/tokens"
)

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

type Boolean struct {
	Value bool
	Token tokens.Token
}

func NewBoolean(t tokens.Token) *Boolean {
	b := &Boolean{
		Token: t,
	}

	b.Value = t.Type == tokens.TRUE

	return b
}
func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}
func (b *Boolean) ToString() string {
	return b.TokenLiteral()
}
