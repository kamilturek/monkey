package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		store: map[string]Object{},
	}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	return &Environment{
		store: map[string]Object{},
		outer: outer,
	}
}

func (e *Environment) Get(name string) (Object, bool) {
	value, ok := e.store[name]
	if !ok && e.outer != nil {
		value, ok = e.outer.Get(name)
	}

	return value, ok
}

func (e *Environment) Set(name string, value Object) Object {
	e.store[name] = value

	return value
}
