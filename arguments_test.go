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
		{"-h"},
		{"--help"},
		{"--version"},
	} {
		_, err := getArguments(ss)
		assert.Nil(t, err)
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
