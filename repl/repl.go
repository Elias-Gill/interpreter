package repl

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

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

	maxTime int64
}

func (r Repl) Run() {
	if r.interactive {
		r.runInteractively()
	} else {
		r.runFromFile()
	}
}

func (r Repl) runInteractively() {
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

		// Run or panic
		r.evaluateOrPanic(buffer)

		// Reset for next command
		buffer = ""
		r.rlInstance.SetPrompt(">>> ")
	}
}

func (r Repl) runFromFile() {
	inReader := bufio.NewReader(r.inFile)

	var input string
	for {
		str, err := inReader.ReadString(byte('\n'))
		if err != nil {
			break
		}
		input += str
	}

	r.evaluateOrPanic(input)
}

// Evaluate results in a separate subrutine. If the max-timeout is reached
// then kill the entire program.
func (r Repl) evaluateOrPanic(input string) {
	c := make(chan struct{})
	go func() {
		switch r.mode {
		case LEXER:
			r.lexe(input)
		case PARSER:
			r.parse(input)
		default:
			r.execute(input)
		}
		c <- struct{}{}
	}()

	// Wait for the evaluator to done or panic on timeout reached
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.maxTime)*time.Millisecond)
	defer cancel()

	select {
	case <-ctx.Done():
		fmt.Fprint(r.errFile, "Timeout execution reached\n")
		os.Exit(1)
	case <-c:
		break
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

func (r Repl) execute(in string) {
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
