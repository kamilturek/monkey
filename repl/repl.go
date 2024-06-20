package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/kamilturek/monkey/lexer"
	"github.com/kamilturek/monkey/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)

		if !scanner.Scan() {
			return
		}

		line := scanner.Text()
		l := lexer.NewLexer(line)
		p := parser.NewParser(l)
		program := p.ParseProgram()

		errors := p.Errors()
		if len(errors) != 0 {
			printParserErrors(out, errors)

			continue
		}

		fmt.Fprint(out, program.String(), "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	fmt.Fprint(out, "Parser Errors:\n")

	for _, msg := range errors {
		fmt.Fprintf(out, "\t%s\n", msg)
	}
}
