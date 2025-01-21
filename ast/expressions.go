// THE CONTENT of a statement is called expression. So a statement is a
// TREE OF EXPRESSIONS.

package ast

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

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
func (i *Identifier) ToString(lvl int) string {
	indent := strings.Repeat("  ", lvl)
	return fmt.Sprintf("%s%s", indent, i.Value)
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
func (i *IntegerLiteral) ToString(lvl int) string {
	indent := strings.Repeat("  ", lvl)
	return fmt.Sprintf("%s%s", indent, i.TokenLiteral())
}

type StringLiteral struct {
	Value string
	Token tokens.Token
}

func NewString(t tokens.Token) *StringLiteral {
	return &StringLiteral{
		Value: t.Literal,
		Token: t,
	}
}
func (i *StringLiteral) expressionNode() {}
func (i *StringLiteral) TokenLiteral() string {
	return i.Token.Literal
}
func (i *StringLiteral) ToString(lvl int) string {
	indent := strings.Repeat("  ", lvl)
	return fmt.Sprintf("%s%s", indent, i.TokenLiteral())
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
func (p *PrefixExpression) ToString(lvl int) string {
	var out bytes.Buffer

	indent := strings.Repeat("  ", lvl)
	out.WriteString(indent + "prefix expression:\n")
	out.WriteString(indent + " operator: " + p.Operator + "\n")
	out.WriteString(indent + " right:\n")
	out.WriteString(p.Right.ToString(lvl + 2)) // Increase indentation for the right expression

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
func (i *InfixExpression) ToString(lvl int) string {
	var out bytes.Buffer

	indent := strings.Repeat("  ", lvl)
	out.WriteString(indent + "infix expression:\n")
	out.WriteString(indent + " left:\n")
	out.WriteString(i.Left.ToString(lvl+2) + "\n") // Increase indentation for the left expression
	out.WriteString(indent + " operator: " + i.Operator + "\n")
	out.WriteString(indent + " right:\n")
	out.WriteString(i.Right.ToString(lvl + 2)) // Increase indentation for the right expression

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
func (b *Boolean) ToString(lvl int) string {
	indent := strings.Repeat("  ", lvl)
	return fmt.Sprintf("%s%s", indent, b.TokenLiteral())
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
func (i *IfExpression) ToString(lvl int) string {
	var buffer bytes.Buffer

	indent := strings.Repeat("  ", lvl)
	buffer.WriteString(indent + "if expression:\n")
	buffer.WriteString(indent + "  condition:\n")
	buffer.WriteString(i.Condition.ToString(lvl + 2)) // Increase indentation for the condition
	buffer.WriteString(indent + "  consequence:\n")
	buffer.WriteString(i.Consequence.ToString(lvl + 2)) // Increase indentation for the consequence

	if i.Alternative != nil {
		buffer.WriteString(indent + "  alternative:\n")
		buffer.WriteString(i.Alternative.ToString(lvl + 2)) // Increase indentation for the alternative
	}

	return buffer.String()
}

type AnonymousFunction struct {
	Parameters []*Identifier
	Body       *BlockStatement
	Token      tokens.Token
}

func NewAnonymousFunction(t tokens.Token) *AnonymousFunction {
	return &AnonymousFunction{
		Token: t,
	}
}

func (f *AnonymousFunction) expressionNode() {}
func (f *AnonymousFunction) TokenLiteral() string {
	return f.Token.Literal
}
func (f *AnonymousFunction) ToString(lvl int) string {
	var buffer bytes.Buffer

	indent := strings.Repeat("  ", lvl)
	buffer.WriteString(indent + "anonymous function:\n")
	buffer.WriteString(indent + "  parameters:\n")
	for _, v := range f.Parameters {
		buffer.WriteString(indent + "    " + v.ToString(lvl) + "\n")
	}
	buffer.WriteString(indent + "  body:\n")
	buffer.WriteString(f.Body.ToString(lvl + 2)) // Increase indentation for the body

	return buffer.String()
}

type FunctionCall struct {
	Arguments  []Expression
	Identifier Expression
	Token      tokens.Token
}

func NewFunctionCall(t tokens.Token, i Expression) *FunctionCall {
	return &FunctionCall{
		Token:      t,
		Identifier: i,
	}
}

func (f *FunctionCall) expressionNode() {}
func (f *FunctionCall) TokenLiteral() string {
	return f.Token.Literal
}
func (f *FunctionCall) ToString(lvl int) string {
	var buffer bytes.Buffer

	indent := strings.Repeat("  ", lvl)
	buffer.WriteString(indent + "function call:\n")
	buffer.WriteString(indent + "  identifier: " + f.Identifier.ToString(0) + "\n")

	// Print the arguments
	buffer.WriteString(indent + "  arguments:\n")
	for _, arg := range f.Arguments {
		buffer.WriteString(arg.ToString(lvl+2) + "\n") // Increase indentation for arguments
	}

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
func (f *ForLoop) ToString(lvl int) string {
	var buffer bytes.Buffer

	indent := strings.Repeat("  ", lvl)

	buffer.WriteString(indent + "for loop:\n")
	buffer.WriteString(indent + " iterations: " + f.Iterations.ToString(0) + "\n")
	buffer.WriteString(indent + " body:\n")
	buffer.WriteString(f.Body.ToString(lvl + 2))

	return buffer.String()
}
