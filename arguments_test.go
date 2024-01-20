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

func TestGetArgumentsError(t *testing.T) {
	for _, ss := range [][]string{
		{},
		{"foo", "my-file.json"},
	} {
		_, err := getArguments(ss)
		assert.NotNil(t, err)
	}
}

func TestHelp(t *testing.T) {
	cupaloy.SnapshotT(t, help())
}
