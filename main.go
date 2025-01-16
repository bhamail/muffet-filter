package main

import (
	"os"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

//goland:noinspection GoUnusedGlobalVariable
var (
	version   = "dev"
	commit    = "none"
	date      = "unknown"
	agentName = "muffet-filter"
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
