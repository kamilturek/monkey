package evaluator

import (
	"fmt"

	"github.com/kamilturek/monkey/ast"
	"github.com/kamilturek/monkey/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node.Statements, env)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{
			Value: node.Value,
		}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ReturnStatement:
		returnValue := Eval(node.Value, env)
		if isError(returnValue) {
			return returnValue
		}

		return &object.ReturnValue{
			Value: returnValue,
		}
	case *ast.LetStatement:
		value := Eval(node.Value, env)
		if isError(value) {
			return value
		}

		env.Set(node.Name.Value, value)

		return NULL
	case *ast.Identifier:
		value, ok := env.Get(node.Value)
		if !ok {
			return newError("identifier not found: %s", node.Value)
		}

		return value
	case *ast.FunctionLiteral:
		return &object.Function{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)
	default:
		return nil
	}
}

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt, env)

		// If `return` or error are encountered, stop executing the block.
		switch obj := result.(type) {
		case *object.ReturnValue:
			return obj.Value
		case *object.Error:
			return obj
		}
	}

	return result
}

func evalBlockStatement(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt, env)

		// If `return` or error are encountered, stop executing the block.
		if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
			return result
		}
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return NULL
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	integer, ok := right.(*object.Integer)
	if !ok {
		return newError("unknown operator: -%s", right.Type())
	}

	return &object.Integer{
		Value: -integer.Value,
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	// Left and right must be booleans then.
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) { //nolint:gocritic
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Integer:
		return obj.Value != 0
	case *object.Boolean:
		return obj.Value
	default:
		return false
	}
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	results := []object.Object{}

	for _, exp := range exps {
		evaluated := Eval(exp, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		results = append(results, evaluated)
	}

	return results
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)
	if !ok {
		return newError("not a function: %s", fn.Type())
	}

	extendedEnv := extendFunctionEnv(function, args)
	evaluated := Eval(function.Body, extendedEnv)

	return unwrapReturnValue(evaluated)
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for i, param := range fn.Parameters {
		env.Set(param.Value, args[i])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftInteger, ok := left.(*object.Integer)
	if !ok {
		return NULL
	}

	rightInteger, ok := right.(*object.Integer)
	if !ok {
		return NULL
	}

	switch operator {
	case "+":
		return &object.Integer{Value: leftInteger.Value + rightInteger.Value}
	case "-":
		return &object.Integer{Value: leftInteger.Value - rightInteger.Value}
	case "*":
		return &object.Integer{Value: leftInteger.Value * rightInteger.Value}
	case "/":
		return &object.Integer{Value: leftInteger.Value / rightInteger.Value}
	case "<":
		return nativeBoolToBooleanObject(leftInteger.Value < rightInteger.Value)
	case ">":
		return nativeBoolToBooleanObject(leftInteger.Value > rightInteger.Value)
	case "==":
		return nativeBoolToBooleanObject(leftInteger.Value == rightInteger.Value)
	case "!=":
		return nativeBoolToBooleanObject(leftInteger.Value != rightInteger.Value)
	default:
		return NULL
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}

	return FALSE
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(format, a...),
	}
}

func isError(obj object.Object) bool {
	return obj.Type() == object.ERROR_OBJ
}
