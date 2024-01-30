package evaluator

import (
	"github.com/sl2.0/ast"
	"github.com/sl2.0/objects"
)

func (e *Evaluator) evalPrefix(exp *ast.PrefixExpression) objects.Object {
	switch exp.Operator {
	case "!":
		return e.evalBangOperator(exp)
	case "-":
		return e.evalMinusPrefix(exp)
	}

	return objects.NewError("Prefix operation not supported: %s", exp.Operator)
}

func (e *Evaluator) evalInfix(exp *ast.InfixExpression) objects.Object {
	evalLeft := e.eval(exp.Left)

	switch evalLeft.Type() {
	case objects.INTEGER_OBJ:
		return e.evalArithmeticOperations(exp)
	case objects.BOOL_OBJ:
		return e.evalBooleanExpression(exp)
	}

	return objects.NewError("Not supported infix operation: %s", exp.Operator)
}

func (e *Evaluator) evalBangOperator(exp *ast.PrefixExpression) objects.Object {
	value := e.eval(exp.Right)

	if value.Type() != objects.BOOL_OBJ {
		return objects.NewError(
			"Expected boolean expression for '!' operator. \n\tGot: %v",
			value.Inspect())
	}

	// we can compare object references because we have static true and false
	if value == true_obj {
		return false_obj
	}

	return true_obj
}

func (e *Evaluator) evalMinusPrefix(exp *ast.PrefixExpression) objects.Object {
	value := e.eval(exp.Right)

	if value.Type() != objects.INTEGER_OBJ {
		return objects.NewError(
			"Expected integer expression for '-' operator. \n\tGot: %v",
			value.Inspect())
	}

	res := value.(*objects.Integer)

	return &objects.Integer{Value: -res.Value}
}

func (e *Evaluator) evalBooleanExpression(exp *ast.InfixExpression) objects.Object {
	left := e.eval(exp.Left).(*objects.Boolean)

	evalRight := e.eval(exp.Right)

	if evalRight.Type() != objects.BOOL_OBJ {
		return objects.NewError(
			"Expected right value to be a boolean.\n\tGot: %v",
			evalRight.Inspect())
	}

	right := evalRight.(*objects.Boolean)

	switch exp.Operator {
	case "==":
		return &objects.Boolean{Value: left.Value == right.Value}
	case "!=":
		return &objects.Boolean{Value: left.Value != right.Value}
	}

	return objects.NewError(
		"Not supported operator: %s",
		exp.Operator)
}

func (e *Evaluator) evalArithmeticOperations(exp *ast.InfixExpression) objects.Object {
	left := e.eval(exp.Left).(*objects.Integer)

	evalRight := e.eval(exp.Right)

	if evalRight.Type() != objects.INTEGER_OBJ {
		return objects.NewError(
			"Expected right value of '%s' to be an integer. \n\tGot: %v",
			exp.Operator, evalRight.Inspect())
	}

	right := evalRight.(*objects.Integer)

	switch exp.Operator {
	case "+":
		return &objects.Integer{Value: left.Value + right.Value}
	case "-":
		return &objects.Integer{Value: left.Value - right.Value}
	case "*":
		return &objects.Integer{Value: left.Value * right.Value}
	case "/":
		return &objects.Integer{Value: left.Value / right.Value}
	case ">":
		return selectBoolObject(left.Value > right.Value)
	case "<":
		return selectBoolObject(left.Value < right.Value)
	case "==":
		return selectBoolObject(left.Value == right.Value)
	case "!=":
		return selectBoolObject(left.Value != right.Value)
	}

	return objects.NewError(
		"Not supported operator: %s",
		exp.Operator,
	)
}

func (e *Evaluator) evalIfExpression(exp *ast.IfExpression) objects.Object {
	condition := e.eval(exp.Condition)

	if condition.Type() != objects.BOOL_OBJ {
		return objects.NewError(
			"Expected boolean expression for 'if' condition.\n\t%v",
			condition.Inspect(),
		)
	}

	if condition == true_obj {
		return e.eval(exp.Consequence)
	}

	if exp.Alternative != nil {
		return e.eval(exp.Alternative)
	}

	return nil
}

func selectBoolObject(exp bool) *objects.Boolean {
	if exp {
		return true_obj
	}

	return false_obj
}
