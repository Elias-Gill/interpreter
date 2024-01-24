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

type IntegerLiteral struct {
	Value int64
	Token tokens.Token
}

func NewInteger(t tokens.Token) *IntegerLiteral {
	value, err := strconv.ParseInt(t.Literal, 0, 64)
	if err != nil {
		return nil
	}

	return &IntegerLiteral{
		Value: value,
		Token: t,
	}
}
func (i *IntegerLiteral) expressionNode() {}
func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}
func (i *IntegerLiteral) ToString() string {
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

	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.ToString())
	out.WriteString(")")

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

type IfExpression struct {
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
	Token       tokens.Token
}

func NewIfExpression(t tokens.Token) *IfExpression {
	return &IfExpression{
		Token: t,
	}
}

func (i *IfExpression) expressionNode() {}
func (i *IfExpression) TokenLiteral() string {
	return i.Token.Literal
}
func (i *IfExpression) ToString() string {
	var buffer bytes.Buffer

	buffer.WriteString(i.TokenLiteral())
	buffer.WriteString(i.Condition.ToString())
	buffer.WriteString(i.Consequence.ToString())

	if i.Alternative != nil {
		buffer.WriteString("sino" + i.Alternative.ToString())
	}

	return buffer.String()
}

// Unnamed functions
type FunctionLiteral struct {
	Paramenters []*Identifier
	Body        *BlockStatement
	Token       tokens.Token
}

func NewFunctionLiteral(t tokens.Token) *FunctionLiteral {
	return &FunctionLiteral{
		Token: t,
	}
}

func (f *FunctionLiteral) expressionNode() {}
func (f *FunctionLiteral) TokenLiteral() string {
	return f.Token.Literal
}
func (f *FunctionLiteral) ToString() string {
	var buffer bytes.Buffer

	buffer.WriteString(f.TokenLiteral())
	buffer.WriteString("(")

	for _, v := range f.Paramenters {
		buffer.WriteString(v.ToString() + ", ")
	}

	buffer.WriteString(")")
	buffer.WriteString(f.Body.ToString())

	return buffer.String()
}

type FunctionCall struct {
	Arguments []Expression
	Ident     Expression
	Token     tokens.Token
}

func NewFunctionCall(t tokens.Token, i Expression) *FunctionCall {
	return &FunctionCall{
		Token: t,
		Ident: i,
	}
}

func (f *FunctionCall) expressionNode() {}
func (f *FunctionCall) TokenLiteral() string {
	return f.Token.Literal
}
func (f *FunctionCall) ToString() string {
	var buffer bytes.Buffer

	buffer.WriteString(f.Ident.ToString())
	buffer.WriteString(f.TokenLiteral())

	for i := 0; i < len(f.Arguments)-1; i++ {
		v := f.Arguments[i]
		buffer.WriteString(v.ToString() + ", ")
	}

	buffer.WriteString(f.Arguments[len(f.Arguments)-1].ToString())
	buffer.WriteString(")")

	return buffer.String()
}

type ForLoop struct {
	Iterations IntegerLiteral
	Body       *BlockStatement
	Token      tokens.Token
}

func NewForLoop(t tokens.Token) *ForLoop {
	return &ForLoop{
		Token: t,
	}
}

func (f *ForLoop) expressionNode() {}
func (f *ForLoop) TokenLiteral() string {
	return f.Token.Literal
}
func (f *ForLoop) ToString() string {
	var buffer bytes.Buffer

	buffer.WriteString(f.TokenLiteral() + " ")
	buffer.WriteString(f.Iterations.ToString())
	buffer.WriteString(f.Body.ToString())

	return buffer.String()
}
