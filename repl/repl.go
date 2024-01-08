package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/sl2.0/lexer"
	"github.com/sl2.0/tokens"
)

func Start(in io.Reader, out io.Writer) {
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
