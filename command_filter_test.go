package main

import (
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
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

func newTestCommandWithStderr(stderr io.Writer, f muffetFactory) *commandFilter {
	return newCommandFilter(
		io.Discard,
		stderr,
		false,
		f,
	)
}

func TestCommandRun(t *testing.T) {
	ok := newTestCommand(
		newFakeMuffetFactory("hello", nil),
	).Run([]string{"https://foo.com"})

	assert.True(t, ok)
}
