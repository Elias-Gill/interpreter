package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/sl2.0/evaluator"
	"github.com/sl2.0/lexer"
	"github.com/sl2.0/parser"
	"github.com/sl2.0/tokens"
)

// Starts a new live lexer session
func StartLiveLexer(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	// read with an infinit loop
	for {
		fmt.Print("\n>>> ")

		hasNext := scanner.Scan()
		if !hasNext {
			return
		}

		line := scanner.Text()
		if line == "exit" {
			return
		}

		l := lexer.NewLexer(line)

		for token := l.NexToken(); token.Type != tokens.EOF; token = l.NexToken() {
			fmt.Printf("\n[Type: %v, Literal: '%v']", token.Type, token.Literal)
		}

		fmt.Print("\n\n")
	}
}

// Starts a new live parsing session
func StartLiveParser(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print("\n>>> ")

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "exit" {
			return
		}

		p := parser.NewParser(line)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printErrors(out, p.Errors())
			continue
		}

		io.WriteString(out, program.ToString())
		io.WriteString(out, "\n")
	}
}

// Starts a new REPL
func StartREPL(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print("\n>>> ")

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "exit" {
			return
		}

		p := parser.NewParser(line)
		program := p.ParseProgram()

		if p.HasErrors() {
			printErrors(out, p.Errors())
			continue
		}

		ev := evaluator.NewFromProgram(program)

		evaluated := ev.EvalProgram()

        if ev.HasErrors() {
            printErrors(out, ev.Errors())
            continue
        }

		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		} else {
			io.WriteString(out, "Feature not implemented\n")
		}
	}
}

func printErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
