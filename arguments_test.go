package main

import (
	"github.com/bradleyjkemp/cupaloy"
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
		assert.EqualError(t, err, "invalid number of arguments")
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
