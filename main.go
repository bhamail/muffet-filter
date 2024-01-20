package main

import (
	"os"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

func main() {
	ok := newCommandFilter(
		colorable.NewColorableStdout(),
		os.Stderr,
		isatty.IsTerminal(os.Stdout.Fd()),
		newRealMuffetFactory(),
	).Run(os.Args[1:])

	if !ok {
		os.Exit(1)
	}
}
