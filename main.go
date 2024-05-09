package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sl2.0/evaluator"
	"github.com/sl2.0/objects"
	"github.com/sl2.0/parser"
	"github.com/sl2.0/repl"
)

func main() {
	const colorMagenta = "\033[35m"
	const colorNone = "\033[0m"

	// Define flags
	parserFlag := flag.Bool("parser", false, "Enable parser flag")
	lexerFlag := flag.Bool("lexer", false, "Enable lexer flag")
	execute := flag.String("exec", "", "Execute the given file")

	// Parse command-line flags
	flag.Parse()

	if *execute != "" {
		runProgram(*execute)
		return
	}

	fmt.Print("\nStarting REPL ")

	// Check if either --parser or --lexer flag is provided
	if *parserFlag {
		fmt.Printf("in%s parser%s mode:\n", colorMagenta, colorNone)
		repl.StartLiveParser(os.Stdin, os.Stdout)
		return
	}

	if *lexerFlag {
		fmt.Printf("in%s lexer%s mode:\n", colorMagenta, colorNone)
		repl.StartLiveLexer(os.Stdin, os.Stdout)
		return
	}

	repl.StartREPL(os.Stdin, os.Stdout)
}

func runProgram(file string) {
	f, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("Error:" + err.Error())
		return
	}
	input := string(f)

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
		for _, msg := range p.Errors() {
			fmt.Println("\t" + msg + "\n")
		}
		return
	}

	if evaluated != nil {
		fmt.Println(evaluated.Inspect())
	} else {
		fmt.Println("Not returned values")
	}
}
