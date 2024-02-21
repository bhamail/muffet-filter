package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestCommandRun(t *testing.T) {
	ok := newTestCommand(
		newFakeMuffetFactory("hello", nil),
	).Run([]string{"https://foo.com"})

	assert.True(t, ok)
}

func TestExecCommandLs(t *testing.T) {
	textOut, textErr, exitCode, err := executeCommand(false, "ls", "-alh")

	assert.True(t, strings.Contains(textOut, ".."), "textOut missing `..`:\n%s", textOut)
	assert.True(t, strings.Contains(textOut, "\n"), "textOut missing newline:\n%s", textOut)
	assert.Equal(t, "", textErr)
	assert.Equal(t, 0, exitCode)
	assert.Nil(t, err)
}
func TestExecCommandExitCodeNonZero(t *testing.T) {
	textOut, textErr, exitCode, err := executeCommand(false, "egrep", "hello", "./main.go")

	assert.Equal(t, "", textOut)
	assert.Equal(t, "", textErr)
	assert.Equal(t, 1, exitCode)
	assert.EqualError(t, err, "exit status 1")
}
func TestExecCommandBadCommand(t *testing.T) {
	textOut, textErr, exitCode, err := executeCommand(false, "noSuchCommand", "-alh")

	assert.Equal(t, "", textOut)
	assert.Equal(t, "", textErr)
	assert.Equal(t, -1, exitCode)
	assert.EqualError(t, err, "exec: \"noSuchCommand\": executable file not found in $PATH")
}
func TestGetMuffetMuffetPathInvalid(t *testing.T) {
	isDownloaded, muffetPath, err := getMuffet(&arguments{MuffetPath: "noSuchMuffetExecutable"})
	assert.EqualError(t, err, "stat noSuchMuffetExecutable: no such file or directory")
	assert.Equal(t, "", muffetPath)
	assert.False(t, isDownloaded)
}
func TestGetMuffetMuffetPathValid(t *testing.T) {
	origTempMuffet := os.TempDir() + getExtractedExecutableName(muffetExecutableBaseName)
	tempMuffetAlreadyExists, _ := doesFileExist(origTempMuffet)
	if !tempMuffetAlreadyExists {
		// download muffet executable
		isDownloaded, muffetPath, err := getMuffet(&arguments{Verbose: true})
		defer func() {
			_ = os.Remove(origTempMuffet)
		}()
		assert.Nil(t, err)
		assert.Equal(t, origTempMuffet, muffetPath)
		assert.True(t, isDownloaded)
	}

	isDownloaded, muffetPath, err := getMuffet(&arguments{MuffetPath: origTempMuffet})
	assert.Nil(t, err)
	assert.Equal(t, origTempMuffet, muffetPath)
	assert.False(t, isDownloaded)
}
