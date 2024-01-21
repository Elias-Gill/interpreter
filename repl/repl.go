package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/sl2.0/lexer"
	"github.com/sl2.0/parser"
	"github.com/sl2.0/tokens"
)

func StartLexerREPL(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	// read with an infinit loop
	for {
		fmt.Print(">>> ")

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

func StartParserREPL(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(">>> ")

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
			printParserErrors(out, p.Errors())
			continue
		}

		io.WriteString(out, program.ToString())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
