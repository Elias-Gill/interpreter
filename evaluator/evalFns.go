package evaluator

import (
	"fmt"

	"github.com/sl2.0/ast"
	"github.com/sl2.0/objects"
)

func (e *Evaluator) evalPrefix(exp *ast.PrefixExpression) objects.Object {
	switch exp.Operator {
	case "!":
		return e.evalBangOperator(exp)
    case "-":
        return e.evalMinusOperator(exp)
	}

	return null_obj
}

func (e *Evaluator) evalBangOperator(exp *ast.PrefixExpression) objects.Object {
	value := e.eval(exp.Right)

	if value.Type() != objects.BOOL_OBJ {
		e.errors = append(e.errors,
			fmt.Sprintf(
				"Expected boolean expression for '!' operator. \nGot: %v",
				value.Type(),
			),
		)

		return null_obj
	}

	// we can compare object references because we have static true and false
	if value == true_obj {
		return false_obj
	}

	return true_obj
}

func (e *Evaluator) evalMinusOperator(exp *ast.PrefixExpression) objects.Object {
    value := e.eval(exp.Right)

    if value.Type() != objects.INTEGER_OBJ {
        e.errors = append(e.errors,
            fmt.Sprintf(
                "Expected integer expression for '-' operator. \nGot: %v",
                value.Type(),
            ),
        )

        return null_obj
    }

    res := value.(*objects.Integer)

    return &objects.Integer{Value: -res.Value}
}

func (e *Evaluator) evalInfix(exp *ast.InfixExpression) objects.Object {
	evalLeft := e.eval(exp.Left)

	switch evalLeft.Type() {
	case objects.INTEGER_OBJ:
		return e.evalArithmeticOperations(exp)
	case objects.BOOL_OBJ:
		return e.evalBooleanExpression(exp)
	}

	e.errors = append(e.errors,
		fmt.Sprintf(
			"Not supported infix operation: %s",
			exp.Operator,
		),
	)

	return null_obj
}

func (e *Evaluator) evalBooleanExpression(exp *ast.InfixExpression) objects.Object {
	left := e.eval(exp.Left).(*objects.Boolean)

	evalRight := e.eval(exp.Right)

	if evalRight.Type() != objects.BOOL_OBJ {
		e.errors = append(e.errors,
			fmt.Sprintf(
				"Expected right value on infix to be a boolean. \nGot: %v",
				string(evalRight.Type()),
			),
		)

		return null_obj
	}

	right := evalRight.(*objects.Boolean)

	switch exp.Operator {
	case "==":
		return &objects.Boolean{Value: left.Value == right.Value}
	case "!=":
		return &objects.Boolean{Value: left.Value != right.Value}
	}

	e.errors = append(e.errors,
		fmt.Sprintf(
			"Not supported operator: %s",
			exp.Operator,
		),
	)

	return null_obj
}

func (e *Evaluator) evalArithmeticOperations(exp *ast.InfixExpression) objects.Object {
	left := e.eval(exp.Left).(*objects.Integer)

	evalRight := e.eval(exp.Right)

	if evalRight.Type() != objects.INTEGER_OBJ {
		e.errors = append(e.errors,
			fmt.Sprintf(
				"Expected right value of '%s' to be an integer. \nGot: %v",
				exp.Operator, string(evalRight.Type()),
			),
		)

		return null_obj
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

	e.errors = append(e.errors,
		fmt.Sprintf(
			"Not supported operator: %s",
			exp.Operator,
		),
	)

	return null_obj
}

func (e *Evaluator) evalIfExpression(exp *ast.IfExpression) objects.Object {
	condition := e.eval(exp.Condition)

	if condition.Type() != objects.BOOL_OBJ {
		e.errors = append(e.errors,
			fmt.Sprintf(
				"Expected boolean expression for 'if' condition. \nGot: %v",
				condition.Type(),
			),
		)

		return null_obj
	}

	if condition == true_obj {
		return e.eval(exp.Consequence)
	}

	if exp.Alternative != nil {
		return e.eval(exp.Alternative)
	}

	return null_obj
}

func selectBoolObject(exp bool) *objects.Boolean {
    if exp {
        return true_obj
    }

    return false_obj
}
