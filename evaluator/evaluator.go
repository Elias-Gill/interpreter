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

/*
eval evaluates every statement or expression within the program recursivelly

eval recieves an storage environment, which is ONLY local to the execution scope.
So there are no global values. For statements scope dependant (like functions or for loops),
a new env has to be created an passed to the eval function.
*/
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
		return e.evalBlockStatement(node)

	case *ast.ReturnStatement:
		val := e.eval(node.ReturnValue)
		return &objects.ReturnObject{Value: val}
	}

	return objects.NewError(
		fmt.Sprintf("Cannot evaluate node: %v", node.ToString()))
}

func (e *Evaluator) evalBlockStatement(node *ast.BlockStatement) objects.Object {
	var res objects.Object

	for _, value := range node.Statements {
		res = e.eval(value)

		if res != nil {
			rt := res.Type()
			if rt == objects.RETURN_OBJ || rt == objects.ERROR_OBJ {
				return res
			}
		}
	}

	return res
}

func (e *Evaluator) evalStatements(stmts []ast.Statement) objects.Object {
	var res objects.Object

	for _, value := range stmts {
		res = e.eval(value)

		switch res := res.(type) {
		case *objects.ReturnObject:
			return res.Value

		case *objects.ErrorObject:
			return res
		}
	}

	return res
}
