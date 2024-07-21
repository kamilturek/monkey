package evaluator

import "github.com/kamilturek/monkey/object"

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{
					Value: int64(len(arg.Value)),
				}
			case *object.Array:
				return &object.Integer{
					Value: int64(len(arg.Elements)),
				}
			default:
				return NewError("argument to `len` not supported, got=%s", args[0].Type())
			}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) == 0 {
					return NULL
				}

				return arg.Elements[0]
			default:
				return NewError("argument to `first` not supported, got=%s", args[0].Type())
			}
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) == 0 {
					return NULL
				}

				return arg.Elements[len(arg.Elements)-1]
			default:
				return NewError("argument to `last` not supported, got=%s", args[0].Type())
			}
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) == 0 {
					return NULL
				}

				newLen := len(arg.Elements) - 1
				newElements := make([]object.Object, newLen)
				copy(newElements, arg.Elements[1:])

				return &object.Array{
					Elements: newElements,
				}
			default:
				return NewError("argument to `rest` not supported, got=%s", args[0].Type())
			}
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return NewError("wrong number of arguments. got=%d, want=2", len(args))
			}

			newElement := args[1]

			switch arg := args[0].(type) {
			case *object.Array:
				newLen := len(arg.Elements) + 1
				newElements := make([]object.Object, newLen-1, newLen)
				copy(newElements, arg.Elements)
				newElements = append(newElements, newElement)

				return &object.Array{
					Elements: newElements,
				}
			default:
				return NewError("argument to `push` not supported, got=%s", args[0].Type())
			}
		},
	},
}
