package main

import (
	"io"
)

func newTestCommand(f muffetFactory) *commandFilter {
	return newTestCommandWithStdout(io.Discard, f)
}

func newTestCommandWithStdout(stdout io.Writer, f muffetFactory) *commandFilter {
	return newCommandFilter(
		stdout,
		io.Discard,
		false,
		f,
	)
}

/*func newTestCommandWithStderr(stderr io.Writer, f muffetFactory) *commandFilter {
	return newCommandFilter(
		io.Discard,
		stderr,
		false,
		f,
	)
}
*/
