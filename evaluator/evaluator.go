package evaluator

import (
	"github.com/sl2.0/ast"
	"github.com/sl2.0/objects"
	"github.com/sl2.0/tokens"
)

var (
	true_obj  = &objects.Boolean{Value: true}
	false_obj = &objects.Boolean{Value: false}
)

func Eval(node ast.Node) objects.Object {
	switch node := node.(type) {
	case *ast.IntegerLiteral:
		return &objects.Integer{Value: node.Value}

	case *ast.Boolean:
		if node.Token.Type == tokens.TRUE {
			return true_obj
		}
		return false_obj

	case *ast.Program:
		return parseStatements(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.PrefixExpression:
		return evalPrefix(node)
	}

	return nil
}

func parseStatements(stmts []ast.Statement) objects.Object {
	var res objects.Object

	for _, value := range stmts {
		res = Eval(value)
	}

	return res
}

func evalPrefix(exp *ast.PrefixExpression) objects.Object {
	switch exp.Operator {
	case "!":
		return evalBangOperator(exp)
	}

	return nil
}

func evalBangOperator(exp *ast.PrefixExpression) objects.Object {
	value := Eval(exp.Right)

	if value.Type() != objects.BOOL_OBJ {
		// TODO: add error messages
		return nil
	}

    // TODO: rethink these lines
    if value.Inspect() == true_obj.Inspect() {
        return false_obj
    }

	return true_obj
}
