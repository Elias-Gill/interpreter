// EVERY "declaration" is an statement. Example: a variable declaration,
// and if statement, a function declration.

package ast

import (
	"bytes"
	"strings"

	"github.com/sl2.0/tokens"
)

type VarStatement struct {
	Identifier *Identifier
	Value      Expression
	Token      tokens.Token
}

func (v *VarStatement) statementNode() {}
func (v *VarStatement) TokenLiteral() string {
	return v.Token.Literal
}

func (v *VarStatement) ToString(lvl int) string {
	var out bytes.Buffer

	indent := strings.Repeat("  ", lvl)
	out.WriteString(indent + "var statement:\n")
	out.WriteString(v.Identifier.ToString(lvl+1))
	out.WriteString(indent + "  value: \n")

	if v.Value != nil {
		out.WriteString(v.Value.ToString(lvl + 2))
	} else {
		out.WriteString("nil")
	}

	out.WriteString("\n")
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
func (r *ReturnStatement) ToString(lvl int) string {
	var out bytes.Buffer

	indent := strings.Repeat("  ", lvl)
	out.WriteString(indent + "return statement:\n")
	out.WriteString(indent + "  value: \n")

	if r.ReturnValue != nil {
		out.WriteString(r.ReturnValue.ToString(lvl + 2))
	} else {
		out.WriteString("nil")
	}

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
func (e *ExpressionStatement) ToString(lvl int) string {
	var out bytes.Buffer

	indent := strings.Repeat("  ", lvl)
	out.WriteString(indent + "expression statement:\n")
	out.WriteString(indent + " expression: \n" + e.Expression.ToString(lvl+2))

	return out.String()
}

type BlockStatement struct {
	Statements []Statement
	Token      tokens.Token // the "{" token
}

func (b *BlockStatement) statementNode() {}
func (b *BlockStatement) TokenLiteral() string {
	return b.Token.Literal
}
func (b *BlockStatement) ToString(lvl int) string {
	var buffer bytes.Buffer

	indent := strings.Repeat("  ", lvl)
	buffer.WriteString(indent + "block statement:\n")
	for _, stmt := range b.Statements {
		buffer.WriteString(stmt.ToString(lvl + 1) + "\n")
	}

	return buffer.String()
}

// Named functions
type FunctionStatement struct {
	Parameters []*Identifier
	Body       *BlockStatement
	Identifier *Identifier
	Token      tokens.Token
}

func NewFunctionStatement(t tokens.Token) *FunctionStatement {
	return &FunctionStatement{
		Token: t,
	}
}

func (f *FunctionStatement) statementNode() {}
func (f *FunctionStatement) TokenLiteral() string {
	return f.Token.Literal
}
func (f *FunctionStatement) ToString(lvl int) string {
	var buffer bytes.Buffer

	indent := strings.Repeat("  ", lvl)
	buffer.WriteString(indent + "function statement:\n")
	buffer.WriteString("  " + f.Identifier.ToString(lvl))
	buffer.WriteString(indent + "  parameters:\n")
	for _, v := range f.Parameters {
		buffer.WriteString(v.ToString(lvl+3))
	}
	buffer.WriteString(indent + "  body:\n")
	buffer.WriteString(f.Body.ToString(lvl + 2))

	return buffer.String()
}
