package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/sl2.0/repl"
)

func main() {
	const colorMagenta = "\033[35m"
	const colorNone = "\033[0m"

	// Define flags
	mode := flag.String("mode", "eval", "Available modes: lexer, parser, eval(default)")
	quiet := flag.Bool("quiet", false, "Suppres unnecesary messages")

	inputFile := flag.String("file", "", "Execute the given file")
	outputFile := flag.String("o", "", "File to output the result")
	errFile := flag.String("err", "", "File to output the errors")

	// Parse command-line flags
	flag.Parse()

	builder := repl.NewReplBuilder().
		WithStderr(os.Stdout).
		WithStdout(os.Stdout).
		WithStdin(os.Stdin)

	// Check stdin for data being piped in
	stat, err := os.Stdin.Stat()
	if (stat.Mode()&os.ModeCharDevice) == 0 && err == nil {
		builder = builder.WithStdin(os.Stdin)
	} else if *inputFile != "" { // File mode
		f, err := os.Open(*inputFile)
		if err != nil {
			log.Fatal("Error opening input file: " + err.Error())
		}
		defer f.Close()
		builder = builder.WithStdin(f)
	} else { // Interactive mode
		builder = builder.Interactive()
	}

	if *outputFile != "" {
		f, err := os.Open(*outputFile)
		if err != nil {
			log.Fatal("Error opening output file: " + err.Error())
		}
		defer f.Close()
		builder = builder.WithStdout(f)
	}

	if *errFile != "" {
		f, err := os.Open(*errFile)
		if err != nil {
			log.Fatal("Error opening error file: " + err.Error())
		}
		defer f.Close()
		builder = builder.WithStderr(f)
	}

	switch *mode {
	case "lexer":
		builder = builder.WithMode(repl.LEXER)
	case "parser":
		builder = builder.WithMode(repl.PARSER)
	}

	rplInstance := builder.Build()

	// On quiet mode this lines are not printed
	if !*quiet {
		fmt.Printf("Starting REPL in %s%s%s mode...\n", colorMagenta, *mode, colorNone)
	}

	rplInstance.Run()
}
