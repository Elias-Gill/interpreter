// EVERY "declaration" is an statement. Example: a variable declaration,
// and if statement, a function declration.

package ast

import (
	"bytes"

	"github.com/sl2.0/tokens"
)

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
