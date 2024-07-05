package evaluator

import (
	"github.com/kamilturek/monkey/ast"
	"github.com/kamilturek/monkey/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.BlockStatement:
		return evalStatements(node.Statements)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{
			Value: node.Value,
		}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)

		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)

		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		condition := Eval(node.Condition)

		if condition.Type() == object.INTEGER_OBJ {
			integer, ok := condition.(*object.Integer)
			if !ok {
				return nil
			}

			if integer.Value != 0 {
				return Eval(node.Consequence)
			} else if node.Alternative != nil {
				return Eval(node.Alternative)
			} else {
				return NULL
			}
		}

		if condition.Type() == object.BOOLEAN_OBJ {
			boolean, ok := condition.(*object.Boolean)
			if !ok {
				return nil
			}

			if boolean.Value {
				return Eval(node.Consequence)
			} else if node.Alternative != nil {
				return Eval(node.Alternative)
			} else {
				return NULL
			}
		}

		return nil
	default:
		return nil
	}
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)
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
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}

	integer, ok := right.(*object.Integer)
	if !ok {
		return NULL
	}

	return &object.Integer{
		Value: -integer.Value,
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	// Left and right must be booleans then.
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return NULL
	}
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)

	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
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
