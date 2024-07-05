package object

type Environment struct {
	store map[string]Object
}

func NewEnvironment() *Environment {
	return &Environment{
		store: map[string]Object{},
	}
}

func (e *Environment) Get(name string) (Object, bool) {
	value, ok := e.store[name]

	return value, ok
}

func (e *Environment) Set(name string, value Object) Object {
	e.store[name] = value

	return value
}
