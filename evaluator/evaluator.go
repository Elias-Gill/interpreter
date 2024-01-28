package evaluator

import (
	"fmt"

	"github.com/sl2.0/ast"
	"github.com/sl2.0/objects"
	"github.com/sl2.0/parser"
	"github.com/sl2.0/tokens"
)

var (
	true_obj  = &objects.Boolean{Value: true}
	false_obj = &objects.Boolean{Value: false}
	null_obj  = &objects.NULL{}
)

type Evaluator struct {
	errors  []string
	program *ast.Program
}

func NewFromInput(input string) *Evaluator {
	eval := &Evaluator{}
	pars := parser.NewParser(input)

	if pars == nil {
		eval.errors = append(eval.errors, "Parser returned a nil value")
		return nil
	}

	eval.program = pars.ParseProgram()

	if pars.HasErrors() {
		eval.errors = pars.Errors()
		return nil
	}

	return eval
}

func NewFromProgram(ast *ast.Program) *Evaluator {
	eval := &Evaluator{}

	if ast == nil {
		eval.errors = append(eval.errors, "Submited an empty(nil) ast")
		return nil
	}

	eval.program = ast

	return eval
}

func (e *Evaluator) Errors() []string {
	return e.errors
}

func (e *Evaluator) HasErrors() bool {
	return len(e.errors) != 0
}

func (e *Evaluator) EvalProgram() objects.Object {
	return e.eval(e.program)
}

func (e *Evaluator) eval(node ast.Node) objects.Object {
	switch node := node.(type) {
	case *ast.IntegerLiteral:
		return &objects.Integer{Value: node.Value}

	case *ast.Boolean:
		if node.Token.Type == tokens.TRUE {
			return true_obj
		}
		return false_obj

	case *ast.Program:
		return e.evalStatements(node.Statements)

	case *ast.ExpressionStatement:
		return e.eval(node.Expression)

	case *ast.PrefixExpression:
		return e.evalPrefix(node)

	case *ast.InfixExpression:
		return e.evalInfix(node)

	case *ast.IfExpression:
		return e.evalIfExpression(node)

	case *ast.BlockStatement:
		return e.evalStatements(node.Statements)
	}

	e.errors = append(e.errors,
		fmt.Sprintf("Cannot evaluate node: %v", node.ToString()))

	return null_obj
}

func (e *Evaluator) evalStatements(stmts []ast.Statement) objects.Object {
	var res objects.Object

	for _, value := range stmts {
		res = e.eval(value)
	}

	return res
}

