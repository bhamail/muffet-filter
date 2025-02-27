package main

import (
	"fmt"
	"github.com/bradleyjkemp/cupaloy"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetArguments(t *testing.T) {
	for _, ss := range [][]string{
		{"my-file.json"},
		{"-v", "my-file.json"},
		{"--verbose", "my-file.json"},
		{"--version"},
	} {
		_, err := getArguments(ss)
		assert.Nil(t, err)
	}
}

func TestGetArgumentsHelp(t *testing.T) {
	for _, ss := range [][]string{
		{"-h"},
		{"--help"},
	} {
		_, err := getArguments(ss)
		assert.ErrorContains(t, err, "Application Options:\n  -m, --muffet-path=          Path to muffet executable\n")
	}
}

func TestGetArgumentsErrorArgsCount(t *testing.T) {
	for _, ss := range [][]string{
		{},
		{"foo", "my-file.json"},
	} {
		_, err := getArguments(ss)
		assert.NotNil(t, err)
		assert.EqualError(t, err, fmt.Sprintf("invalid number of arguments\n\n%s", help()))
	}
}

func TestGetArgumentsMuffetArg(t *testing.T) {
	ss := []string{"--muffet-arg=--one-page-only", "--muffet-arg=--max-connections-per-host=10",
		"--muffet-arg=--verbose", "--muffet-arg=--verbose", "my-url"}

	args, err := getArguments(ss)

	assert.Nil(t, err)
	assert.Equal(t, "--one-page-only", args.MuffetArg[0])
	assert.Equal(t, "--max-connections-per-host=10", args.MuffetArg[1])
	assert.Equal(t, "--verbose", args.MuffetArg[2])
	assert.Equal(t, "--verbose", args.MuffetArg[3])
	assert.Equal(t, "my-url", args.URL)
}

func TestGetArgumentsErrorUnknownFlag(t *testing.T) {
	for _, ss := range [][]string{
		{"--bogusArg"},
	} {
		_, err := getArguments(ss)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "unknown flag `bogusArg'")
	}
}

func TestHelp(t *testing.T) {
	cupaloy.SnapshotT(t, help())
}

func TestGetUserHomeDir(t *testing.T) {
	homeDir, err := getUserHomeDir()
	assert.Nil(t, err)
	assert.NotNil(t, homeDir)
}

func TestGetUserHomeDirError(t *testing.T) {
	// force error by removing expected env var
	origHome := os.Getenv("HOME")
	defer func() {
		_ = os.Setenv("HOME", origHome)
	}()
	_ = os.Unsetenv("HOME")

	homeDir, err := getUserHomeDir()
	assert.EqualError(t, err, "$HOME is not defined")
	assert.Equal(t, "", homeDir)
}
