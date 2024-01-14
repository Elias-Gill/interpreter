package main

import (
	"fmt"
	"os"

	"github.com/sl2.0/repl"
)

func main() {
	fmt.Println("Starting REPL...")
	repl.Start(os.Stdin, os.Stdout)
}
