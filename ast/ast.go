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

type Program struct {
	Statements []Statement
}

func (p *Program) ToString() string {
	var buffer bytes.Buffer

	for _, stmt := range p.Statements {
		buffer.WriteString(stmt.ToString())
	}

	return buffer.String()
}

func (p *Program) TokenLiteral() string {
    if len(p.Statements) > 0 {
        return p.Statements[0].TokenLiteral()
    } else {
        return ""
    }
}
