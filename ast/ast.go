package ast

import (
	"bytes"
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

func (a *Ast) ToString() string {
	var buffer bytes.Buffer

	for _, stmt := range a.Statements {
		buffer.WriteString(stmt.ToString())
	}

	return buffer.String()
}
