package main

import (
	"os"

	"github.com/kamilturek/monkey/interpreter"
)

func main() {
	interp, _ := interpreter.NewInterpreter(
		interpreter.WithInputFromArgs(os.Args[1:]),
	)
	interp.Start()
}
