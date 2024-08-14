package interpreter

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"

	"github.com/kamilturek/monkey/evaluator"
	"github.com/kamilturek/monkey/lexer"
	"github.com/kamilturek/monkey/object"
	"github.com/kamilturek/monkey/parser"
)

const PROMPT = ">> "

type interpreter struct {
	input  io.Reader
	output io.Writer
	repl   bool
}

type option func(interpreter *interpreter) error

func NewInterpreter(opts ...option) (*interpreter, error) {
	interp := &interpreter{
		input:  os.Stdin,
		output: os.Stdout,
		repl:   true,
	}

	for _, opt := range opts {
		if err := opt(interp); err != nil {
			return nil, err
		}
	}

	return interp, nil
}

func WithInputFromArgs(args []string) option {
	return func(interp *interpreter) error {
		if len(args) < 1 {
			return nil
		}

		f, err := os.Open(args[0])
		if err != nil {
			return err
		}

		interp.input = f
		interp.repl = false

		return nil
	}
}

func (interp *interpreter) Start() {
	if interp.repl {
		interp.startRepl()
	} else {
		interp.executeFile()
	}
}

func (interp *interpreter) executeFile() {
	env := object.NewEnvironment()

	rawBody, err := io.ReadAll(interp.input)
	if err != nil {
		panic(err)
	}

	body := string(rawBody)
	l := lexer.NewLexer(body)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	errors := p.Errors()
	if len(errors) != 0 {
		interp.printParserErrors(errors)

		return
	}

	evaluated := evaluator.Eval(program, env)

	if err, ok := evaluated.(*object.Error); ok {
		fmt.Fprintf(interp.output, "%s\n", err.Inspect())
	}
}

func (interp *interpreter) startRepl() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")

	scanner := bufio.NewScanner(interp.input)
	env := object.NewEnvironment()

	for {
		fmt.Fprint(interp.output, PROMPT)

		if !scanner.Scan() {
			return
		}

		line := scanner.Text()
		l := lexer.NewLexer(line)
		p := parser.NewParser(l)
		program := p.ParseProgram()

		errors := p.Errors()
		if len(errors) != 0 {
			interp.printParserErrors(errors)

			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			fmt.Fprintf(interp.output, "%s\n", evaluated.Inspect())
		}
	}
}

func (interpreter *interpreter) printParserErrors(errors []string) {
	fmt.Fprint(interpreter.output, "Parser Errors:\n")

	for _, msg := range errors {
		fmt.Fprintf(interpreter.output, "\t%s\n", msg)
	}
}
