package repl

import (
	"io"
	"os"

	"github.com/chzyer/readline"
	"github.com/sl2.0/objects"
)

type ReplBuilder struct {
	repl Repl
}

func NewReplBuilder() ReplBuilder {
	return ReplBuilder{
		repl: Repl{
			inFile:      os.Stdin,
			outFile:     os.Stdout,
			errFile:     os.Stderr,
			env:         objects.NewStorage(),
			mode:        EVAL,
			interactive: false,
		},
	}
}

func (r ReplBuilder) WithStdin(file io.ReadCloser) ReplBuilder {
	r.repl.inFile = file
	return r
}

func (r ReplBuilder) WithStdout(file io.WriteCloser) ReplBuilder {
	r.repl.outFile = file
	return r
}

func (r ReplBuilder) WithStderr(file io.WriteCloser) ReplBuilder {
	r.repl.errFile = file
	return r
}

func (r ReplBuilder) WithMode(mode mode) ReplBuilder {
	r.repl.mode = mode
	return r
}

func (r ReplBuilder) Interactive() ReplBuilder {
	r.repl.interactive = true
	return r
}

func (b ReplBuilder) Build() Repl {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          ">>> ",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		Stdin:           b.repl.inFile,
	})
	if err != nil {
		panic(err.Error())
	}

	b.repl.rlInstance = rl

	return b.repl
}
