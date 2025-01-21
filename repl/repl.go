package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"

	"github.com/sl2.0/evaluator"
	"github.com/sl2.0/lexer"
	"github.com/sl2.0/objects"
	"github.com/sl2.0/parser"
	"github.com/sl2.0/tokens"
)

// Starts a new live lexer session with multi-line input support
func StartLiveLexer(in io.Reader, out io.Writer) {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          ">>> ",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	buffer := ""
	for {
		line, err := rl.Readline()
		if err == readline.ErrInterrupt { // Handle Ctrl+C
			if buffer != "" {
				buffer = "" // Clear buffer
				fmt.Fprintln(out, "Input cleared.")
				continue
			}
			break
		} else if err != nil { // Handle EOF or errors
			break
		}

		// Detect multi-line input with Shift+Enter or empty buffer with Enter
		if strings.HasSuffix(line, "\\") {
			buffer += strings.TrimSuffix(line, "\\") + "\n"
			rl.SetPrompt("... ") // Change prompt for multi-line input
			continue
		}

		buffer += line
		if buffer == "exit" { // Exit condition
			break
		}

		// Tokenize and output results
		l := lexer.NewLexer(buffer)
		for token := l.NexToken(); token.Type != tokens.EOF; token = l.NexToken() {
			fmt.Fprintf(out, "\n[Type: %v, Literal: '%v']", token.Type, token.Literal)
		}
		fmt.Fprintln(out)

		// Reset for next command
		buffer = ""
		rl.SetPrompt(">>> ")
	}
}

// Starts a new live parser session with multi-line input support
func StartLiveParser(in io.Reader, out io.Writer) {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          ">>> ",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	buffer := ""
	for {
		line, err := rl.Readline()
		if err == readline.ErrInterrupt { // Handle Ctrl+C
			if buffer != "" {
				buffer = "" // Clear buffer
				fmt.Fprintln(out, "Input cleared.")
				continue
			}
			break
		} else if err != nil { // Handle EOF or errors
			break
		}

		// Detect multi-line input with Shift+Enter or empty buffer with Enter
		if strings.HasSuffix(line, "\\") {
			buffer += strings.TrimSuffix(line, "\\") + "\n"
			rl.SetPrompt("... ") // Change prompt for multi-line input
			continue
		}

		buffer += line
		if buffer == "exit" { // Exit condition
			break
		}

		// Parse and output results
		p := parser.NewParser(buffer)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printErrors(out, p.Errors()) // Print errors if any
			buffer = ""                  // Clear buffer after errors
			rl.SetPrompt(">>> ")         // Reset prompt
			continue
		}

		io.WriteString(out, program.ToString())
		io.WriteString(out, "\n")

		// Reset for next command
		buffer = ""
		rl.SetPrompt(">>> ")
	}
}

func StartREPL(in io.Reader, out io.Writer) {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          ">>> ",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	env := objects.NewStorage()
	buffer := ""

	for {
		line, err := rl.Readline()
		if err == readline.ErrInterrupt { // Handle Ctrl+C
			if buffer != "" {
				buffer = "" // Clear the buffer
				fmt.Fprintln(rl.Stderr(), "\nInput cleared.")
				continue
			}
			break
		} else if err == io.EOF { // Handle EOF (Ctrl+D)
			break
		}

		// Detect multi-line input
		if strings.HasSuffix(line, "\\") { // Use `\\` to continue to the next line
			buffer += strings.TrimSuffix(line, "\\") + "\n"
			rl.SetPrompt("... ") // Change prompt for multi-line input
			continue
		}

		buffer += line
		if buffer == "exit" {
			break
		}

		// Parse and evaluate the complete input
		p := parser.NewParser(buffer)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printErrors(out, p.Errors())
			buffer = ""          // Clear the buffer after errors
			rl.SetPrompt(">>> ") // Reset prompt
			continue
		}

		ev := evaluator.NewFromProgram(program)
		evaluated := ev.EvalProgram(env)

		if ev.HasErrors() {
			printErrors(out, ev.Errors())
		} else if evaluated != nil {
			fmt.Fprintln(out, evaluated.Inspect())
		} else {
			fmt.Fprintln(out, "Feature not implemented")
		}

		buffer = ""          // Clear the buffer after evaluation
		rl.SetPrompt(">>> ") // Reset prompt
	}
}

func EvaluateProgram(in *os.File, out *os.File) {
	inReader := bufio.NewReader(in)

	var input string
	for {
		str, err := inReader.ReadString(byte('\n'))
		if err != nil {
			break
		}
		input += str
	}

	p := parser.NewParser(input)
	program := p.ParseProgram()

	if p.HasErrors() {
		for _, msg := range p.Errors() {
			fmt.Println("\t" + msg + "\n")
		}
		return
	}

	ev := evaluator.NewFromProgram(program)
	evaluated := ev.EvalProgram(objects.NewStorage())

	if ev.HasErrors() {
		printErrors(out, ev.Errors())
		return
	}

	if evaluated != nil {
		fmt.Println(evaluated.Inspect())
	} else {
		fmt.Println("Not returned values")
	}
}

func printErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
