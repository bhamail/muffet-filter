package main

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockMuffetExecutor struct {
	result string
	err    error
}

func (m *mockMuffetExecutor) Check(args *arguments) (string, error) {
	return m.result, m.err
}

type mockMuffetFactory struct {
	executor *mockMuffetExecutor
}

func (m *mockMuffetFactory) Create(options muffetOptions) muffetExecutor {
	return m.executor
}

func TestCommandFilter_Help(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cf := newCommandFilter(stdout, stderr, false, &mockMuffetFactory{})

	ok := cf.Run([]string{"--help"})

	assert.True(t, ok)
	assert.Contains(t, stdout.String(), "Usage:")
	assert.Empty(t, stderr.String())
}

func TestCommandFilter_Version(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cf := newCommandFilter(stdout, stderr, false, &mockMuffetFactory{})

	// Save original values
	origVersion := version
	origCommit := commit
	origDate := date
	defer func() {
		version = origVersion
		commit = origCommit
		date = origDate
	}()

	version = "1.2.3"
	commit = "abc"
	date = "2023-01-01"

	ok := cf.Run([]string{"--version"})

	assert.True(t, ok)
	assert.Contains(t, stdout.String(), "muffet-filter version \"1.2.3\" 2023-01-01 (abc)")
	assert.Empty(t, stderr.String())
}

func TestCommandFilter_ArgumentError(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cf := newCommandFilter(stdout, stderr, false, &mockMuffetFactory{})

	// Too many arguments
	ok := cf.Run([]string{"url1", "url2"})

	assert.False(t, ok)
	assert.Contains(t, stderr.String(), "invalid number of arguments")
}

func TestCommandFilter_MuffetError(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	mockExec := &mockMuffetExecutor{err: errors.New("muffet failed")}
	cf := newCommandFilter(stdout, stderr, false, &mockMuffetFactory{executor: mockExec})

	ok := cf.Run([]string{"http://example.com"})

	assert.False(t, ok)
	assert.Contains(t, stderr.String(), "muffet failed")
}

func TestCommandFilter_Success(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	// Empty array of results means no broken links found by muffet
	mockExec := &mockMuffetExecutor{result: "[]"}
	cf := newCommandFilter(stdout, stderr, false, &mockMuffetFactory{executor: mockExec})

	ok := cf.Run([]string{"http://example.com"})

	assert.True(t, ok)
	assert.Empty(t, stdout.String())
	assert.Empty(t, stderr.String())
}

func TestCommandFilter_BrokenLinks(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	// Sample JSON report from muffet
	jsonReport := `[
		{
			"url": "http://example.com",
			"links": [
				{
					"url": "http://broken.com",
					"error": "404 Not Found"
				}
			]
		}
	]`

	mockExec := &mockMuffetExecutor{result: jsonReport}
	cf := newCommandFilter(stdout, stderr, false, &mockMuffetFactory{executor: mockExec})

	ok := cf.Run([]string{"http://example.com"})

	assert.False(t, ok)
	assert.Contains(t, stdout.String(), "http://broken.com")
}

func TestCommandFilter_Run(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	mockExec := &mockMuffetExecutor{err: errors.New("execution error")}
	cf := newCommandFilter(stdout, stderr, false, &mockMuffetFactory{executor: mockExec})

	ok := cf.Run([]string{"http://example.com"})

	assert.False(t, ok)
	assert.Contains(t, stderr.String(), "execution error")
}

func TestCommandFilter_PrintErrorTerminal(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cf := newCommandFilter(stdout, stderr, true, &mockMuffetFactory{})

	cf.printError("some error")

	assert.Contains(t, stderr.String(), "some error")
	// Should contain ANSI escape codes for red color
	assert.Contains(t, stderr.String(), "\x1b[31m")
}

func TestCommandFilter_IgnoreListError(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	mockExec := &mockMuffetExecutor{result: "[]"}
	cf := newCommandFilter(stdout, stderr, false, &mockMuffetFactory{executor: mockExec})

	ok := cf.Run([]string{"--ignores", "non-existent.json", "http://example.com"})

	assert.False(t, ok)
	assert.Contains(t, stderr.String(), "no such file or directory")
}

func TestCommandFilter_LoadReportError(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	mockExec := &mockMuffetExecutor{result: "invalid json"}
	cf := newCommandFilter(stdout, stderr, false, &mockMuffetFactory{executor: mockExec})

	ok := cf.Run([]string{"http://example.com"})

	assert.False(t, ok)
	assert.Contains(t, stderr.String(), "invalid character")
}
