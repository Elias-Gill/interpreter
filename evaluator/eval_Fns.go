package evaluator

import (
	"github.com/sl2.0/ast"
	"github.com/sl2.0/objects"
)

func (e *Evaluator) evalPrefix(exp *ast.PrefixExpression, env *objects.Storage) objects.Object {
	switch exp.Operator {
	case "!":
		return e.evalBangOperator(exp, env)
	case "-":
		return e.evalMinusPrefix(exp, env)
	}

	return objects.NewError("Prefix operation not supported: %s", exp.Operator)
}

func (e *Evaluator) evalInfix(exp *ast.InfixExpression, env *objects.Storage) objects.Object {
	evalLeft := e.eval(exp.Left, env)

	switch evalLeft.Type() {
	case objects.INTEGER_OBJ:
		return e.evalArithmeticOperations(exp, env)
	case objects.BOOL_OBJ:
		return e.evalBooleanExpression(exp, env)
	case objects.STRING_OBJ:
		return e.evalStringExpression(exp, env)
	}

	return objects.NewError("Not supported infix operation: %s", exp.Operator)
}

func (e *Evaluator) evalBangOperator(exp *ast.PrefixExpression, env *objects.Storage) objects.Object {
	value := e.eval(exp.Right, env)

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

func (e *Evaluator) evalMinusPrefix(exp *ast.PrefixExpression, env *objects.Storage) objects.Object {
	value := e.eval(exp.Right, env)

	if value.Type() != objects.INTEGER_OBJ {
		return objects.NewError(
			"Expected integer expression for '-' operator. \n\tGot: %v",
			value.Inspect())
	}

	res := value.(*objects.Integer)

	return &objects.Integer{Value: -res.Value}
}

func (e *Evaluator) evalBooleanExpression(exp *ast.InfixExpression, env *objects.Storage) objects.Object {
	left := e.eval(exp.Left, env).(*objects.Boolean)

	evalRight := e.eval(exp.Right, env)

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

func (e *Evaluator) evalStringExpression(exp *ast.InfixExpression, env *objects.Storage) objects.Object {
	left := e.eval(exp.Left, env).(*objects.String)

	evalRight := e.eval(exp.Right, env)

	if evalRight.Type() != objects.STRING_OBJ {
		return objects.NewError(
			"Expected right value to be a String.\n\tGot: %v",
			evalRight.Inspect())
	}

	right := evalRight.(*objects.String)

	switch exp.Operator {
	case "==":
		return &objects.Boolean{Value: left.Value == right.Value}
	case "!=":
		return &objects.Boolean{Value: left.Value != right.Value}
	case "+":
		return &objects.String{Value: left.Value + right.Value}
	}

	return objects.NewError(
		"Not supported operator: %s",
		exp.Operator)
}

func (e *Evaluator) evalArithmeticOperations(exp *ast.InfixExpression, env *objects.Storage) objects.Object {
	left := e.eval(exp.Left, env).(*objects.Integer)

	evalRight := e.eval(exp.Right, env)

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

func (e *Evaluator) evalIfExpression(exp *ast.IfExpression, env *objects.Storage) objects.Object {
	condition := e.eval(exp.Condition, env)

	if condition.Type() != objects.BOOL_OBJ {
		return objects.NewError(
			"Expected boolean expression for 'if' condition.\n\t%v",
			condition.Inspect(),
		)
	}

	if condition == true_obj {
		return e.eval(exp.Consequence, env)
	}

	if exp.Alternative != nil {
		return e.eval(exp.Alternative, env)
	}

	return nil
}

func (e *Evaluator) evalFunctionCall(fun *ast.FunctionCall, env *objects.Storage) objects.Object {
	f, ok := e.eval(fun.Identifier, env).(*objects.FunctionObject)
	if !ok {
		return objects.NewError("Function '%s' not found", fun.Identifier.ToString(0))
	}

	// check argument list size
	if len(fun.Arguments) != len(f.Parameters) {
		return objects.NewError("Number of Arguments mismatch with number of Parameters")
	}

	// eval every argument
	args := e.evalExpressions(fun.Arguments, env)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}

	// Create a local scope (with maximum recurssion level)
	localEnv, err := objects.NewEnclosedStorage(env)
	if err != nil {
		return objects.NewError("%s", err.Error())
	}

	for i, param := range f.Parameters {
		localEnv.Set(param.Value, args[i])
	}

	// unwrap the returned value
	result := e.eval(f.Body, localEnv)
	unwrapped, ok := result.(*objects.ReturnObject)
	if ok {
		return unwrapped.Value
	} else {
		return result
	}
}

func (e *Evaluator) evalForLoop(exp *ast.ForLoop, env *objects.Storage) objects.Object {
	var value objects.Object
	for i := 0; i < int(exp.Iterations.Value); i++ {
		value = e.evalBlockStatement(exp.Body, env)
	}
	return value
}

func selectBoolObject(exp bool) *objects.Boolean {
	if exp {
		return true_obj
	}

	return false_obj
}
