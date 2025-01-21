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
	liveParser := *flag.Bool("parser", false, "Enable parser flag")
	liveLexer := *flag.Bool("lexer", false, "Enable lexer flag")
	inputFile := *flag.String("file", "", "Execute the given file")

	// Parse command-line flags
	flag.Parse()

	// Check stdin for data being piped in
	stat, err := os.Stdin.Stat()
	if (stat.Mode()&os.ModeCharDevice) == 0 && err == nil {
		repl.EvaluateProgram(os.Stdin, os.Stdout)
		return
	}

	// If a file is given
	if inputFile != "" {
		f, err := os.Open(inputFile)
		if err != nil {
			fmt.Println("Error opening file: " + err.Error())
			return
		}

		repl.EvaluateProgram(f, os.Stdout)
		return
	}

	fmt.Print("\nStarting REPL ")

	// Check if either --parser or --lexer flag is provided
	if liveParser {
		fmt.Printf("in%s parser%s mode:\n", colorMagenta, colorNone)
		repl.StartLiveParser(os.Stdin, os.Stdout)
		return
	}

	if liveLexer {
		fmt.Printf("in%s lexer%s mode:\n", colorMagenta, colorNone)
		repl.StartLiveLexer(os.Stdin, os.Stdout)
		return
	}

	fmt.Printf("in%s eval%s mode:\n", colorMagenta, colorNone)
	repl.StartREPL(os.Stdin, os.Stdout)
}
