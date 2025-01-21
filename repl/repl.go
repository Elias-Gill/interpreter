package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/chzyer/readline"

	"github.com/sl2.0/evaluator"
	"github.com/sl2.0/lexer"
	"github.com/sl2.0/objects"
	"github.com/sl2.0/parser"
	"github.com/sl2.0/tokens"
)

type mode int

const (
	EVAL = iota
	PARSER
	LEXER
)

type Repl struct {
	inFile  io.ReadCloser
	outFile io.WriteCloser
	errFile io.WriteCloser

	mode        mode
	interactive bool

	rlInstance *readline.Instance
	env        *objects.Storage
}

func (r Repl) Run() {
	if r.interactive {
		r.runInteractive()
	} else {
		r.fromFile()
	}
}

func (r Repl) runInteractive() {
	defer r.rlInstance.Close()
	buffer := ""
	for {
		line, err := r.rlInstance.Readline()
		if err == readline.ErrInterrupt { // Handle Ctrl+C
			if buffer != "" {
				buffer = "" // Clear buffer
				fmt.Fprintln(r.errFile, "Kill Signal Recieved")
				continue
			}
			break
		} else if err != nil { // Handle EOF or errors
			break
		}

		// Detect multi-line input with Shift+Enter or empty buffer with Enter
		if strings.HasSuffix(line, "\\") {
			buffer += strings.TrimSuffix(line, "\\") + "\n"
			r.rlInstance.SetPrompt("... ") // Change prompt for multi-line input
			continue
		}

		buffer += line
		if buffer == "exit" { // Exit condition
			break
		}

		// Evaluate results
		switch r.mode {
		case LEXER:
			r.lexe(buffer)
		case PARSER:
			r.parse(buffer)
		default:
			r.evaluate(buffer)
		}

		// Reset for next command
		buffer = ""
		r.rlInstance.SetPrompt(">>> ")
	}
}

func (r Repl) fromFile() {
	inReader := bufio.NewReader(r.inFile)

	var input string
	for {
		str, err := inReader.ReadString(byte('\n'))
		if err != nil {
			break
		}
		input += str
	}

	switch r.mode {
	case LEXER:
		r.lexe(input)
		return
	case PARSER:
		r.parse(input)
		return
	default:
		r.evaluate(input)
	}
}

func (r Repl) lexe(in string) {
	l := lexer.NewLexer(in)
	for token := l.NexToken(); token.Type != tokens.EOF; token = l.NexToken() {
		fmt.Fprintf(r.outFile, "[Type: %v, Literal: '%v']\n", token.Type, token.Literal)
	}
	fmt.Fprintln(r.outFile)
}

func (r Repl) parse(in string) {
	// Parse and output results
	p := parser.NewParser(in)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printErrors(r.errFile, p.Errors()) // Print errors if any
	} else {
		fmt.Fprintf(r.outFile, "%v", program.ToString(0))
		fmt.Fprintln(r.outFile)
	}
}

func (r Repl) evaluate(in string) {
	// Parse and evaluate the complete input
	p := parser.NewParser(in)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printErrors(r.outFile, p.Errors())
	} else {
		ev := evaluator.NewFromProgram(program)
		evaluated := ev.EvalProgram(r.env)

		if ev.HasErrors() {
			if p.HasErrors() {
				printErrors(r.errFile, p.Errors())
			}

			ev := evaluator.NewFromProgram(program)
			evaluated := ev.EvalProgram(objects.NewStorage())

			if ev.HasErrors() {
				printErrors(r.errFile, p.Errors())
				return
			}

			if evaluated != nil {
				fmt.Fprintln(r.outFile, evaluated.Inspect())
			} else {
				fmt.Fprintln(r.outFile, "No returned values")
			}
			printErrors(r.errFile, ev.Errors())
		} else if evaluated != nil {
			fmt.Fprintln(r.outFile, evaluated.Inspect())
		} else {
			fmt.Fprintln(r.outFile)
		}
	}
}

func printErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
