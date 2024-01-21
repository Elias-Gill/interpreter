package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sl2.0/repl"
)

func main() {
	const colorMagenta = "\033[35m"
	const colorNone = "\033[0m"

	// Define flags
	parserFlag := flag.Bool("parser", false, "Enable parser flag")
	lexerFlag := flag.Bool("lexer", false, "Enable lexer flag")

	// Parse command-line flags
	flag.Parse()

	fmt.Print("\nStarting REPL ")

	// Check if either --parser or --lexer flag is provided
	if *parserFlag {
		fmt.Printf("in%s parser%s mode:\n", colorMagenta, colorNone)
		repl.StartParserREPL(os.Stdin, os.Stdout)
		return
	}

	if *lexerFlag {
		fmt.Printf("in%s lexer%s mode:\n", colorMagenta, colorNone)
		repl.StartLexerREPL(os.Stdin, os.Stdout)
		return
	}
}
