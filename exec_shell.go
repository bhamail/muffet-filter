package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

func executeCommand(isVerbose bool, name string, options ...string) (textOut string, textErr string, exitCode int, err error) {
	exitCode = -1
	cmd := exec.Command(name, options...)

	var cmdStdOutCloser io.ReadCloser
	cmdStdOutCloser, err = cmd.StdoutPipe()
	defer func(cmdStdOutCloser io.Closer) {
		_ = cmdStdOutCloser.Close()
	}(cmdStdOutCloser)
	if err != nil {
		log.Fatalf("command failed with %s\n", err)
		return
	}

	var cmdStdErrCloser io.ReadCloser
	cmdStdErrCloser, err = cmd.StderrPipe()
	defer func(cmdStdErrCloser io.Closer) {
		_ = cmdStdErrCloser.Close()
	}(cmdStdErrCloser)
	if err != nil {
		log.Fatalf("command failed with %s\n", err)
		return
	}

	bufOutReader := bufio.NewReader(cmdStdOutCloser)
	bufErrReader := bufio.NewReader(cmdStdErrCloser)

	err = cmd.Start()
	if err != nil {
		return
	}
	// Read stdout
	// capture in string
	for {
		line, err := bufOutReader.ReadBytes('\n')
		if err == io.EOF {
			if len(line) == 0 {
				break
			}
		} else {
			if err != nil {
				log.Fatal(err)
			}
		}
		textOut = textOut + string(line)
		if isVerbose {
			_, _ = os.Stdout.Write(line)
		}
	}
	// Read stderr
	for {
		line, err := bufErrReader.ReadBytes('\n')
		if err == io.EOF {
			if len(line) == 0 {
				break
			}
		} else {
			if err != nil {
				log.Fatal(err)
			}
		}
		textErr = textErr + string(line)
		if isVerbose {
			_, _ = os.Stderr.Write(line)
		}
	}
	// caller may choose to ignore exit code, e.g. because failed links result in non-zero exit code
	err = cmd.Wait()
	exitCode = cmd.ProcessState.ExitCode()
	if isVerbose {
		fmt.Printf("exec: %s, args: %s, exit status: %d\n", cmd.Path, cmd.Args, exitCode)
	}
	return
}

const muffetExecutableBaseName = "muffet"

func getMuffet(args *arguments) (isDownloaded bool, muffetPath string, err error) {
	var muffetExec string
	if args.MuffetPath != "" {
		muffetExec = args.MuffetPath
		var itExists bool
		if itExists, err = doesFileExist(muffetExec); !itExists {
			// a non-default file was specified, so it is an error if that specified file is missing
			return
		}
	} else {
		muffetExec = muffetExecutableBaseName
	}

	// fetch muffet if not on path
	_, _, _, err = executeCommand(args.Verbose, muffetExec, "--version")
	if err != nil {
		if err.Error() == "exec: \""+muffetExecutableBaseName+"\": executable file not found in $PATH" {

			// see if we've already downloaded the executable to the temp dir
			extractedExecutableName := getExtractedExecutableName(muffetExecutableBaseName)
			tempDir := getTempDirWTrailingSlash()
			muffetPath = tempDir + extractedExecutableName
			var itExists bool
			if itExists, err = doesFileExist(muffetPath); itExists {
				// show muffet version of found temp file
				var versionOut string
				versionOut, _, _, err = executeCommand(args.Verbose, muffetPath, "--version")
				if err != nil {
					return
				}
				log.Printf("using muffet: %s, version: %s", muffetPath, versionOut)
				return
			}

			// fetch muffet to local temp dir

			// reviewed logic from: https://github.com/ruzickap/action-my-broken-link-checker/blob/main/entrypoint.sh#L67

			// get latest muffet release via GitHub API
			var textOut, textErr string
			var exitCode int
			// this command is really noisy, so don't be verbose even if args ask for it.
			textOut, textErr, exitCode, err = executeCommand(false, "curl", "-s", "https://api.github.com/repos/raviqqe/muffet/releases/latest")
			if err != nil {
				printCommandResults(textOut, textErr, exitCode)
				return
			}

			// find URL to release bundle to download
			// "browser_download_url": "https://github.com/raviqqe/muffet/releases/download/v2.9.3/muffet_darwin_amd64.tar.gz"
			r := regexp.MustCompile("browser_download_url\": \"(.*/releases/download/.*/(muffet_" + runtime.GOOS + "_" + runtime.GOARCH + ".*))\"")
			theRegExMatches := r.FindStringSubmatch(textOut)
			if theRegExMatches == nil {
				err = errors.New("could not find download url")
				return
			} else if len(theRegExMatches) != 3 {
				err = errors.New("unexpected download url regex result")
				return
			}
			latestVersionUrl := theRegExMatches[1]
			if args.Verbose {
				log.Println("fetching: " + latestVersionUrl)
			}

			// download to temp directory
			textOut, textErr, exitCode, err = executeCommand(args.Verbose, "wget", "--no-verbose", "--directory-prefix="+tempDir, latestVersionUrl)
			if err != nil {
				printCommandResults(textOut, textErr, exitCode)
				return
			}

			// extract executable
			gzipFile := tempDir + theRegExMatches[2]
			textOut, textErr, exitCode, err = executeCommand(args.Verbose, "tar", "-C", tempDir, "-xf", gzipFile, extractedExecutableName)
			if err != nil {
				printCommandResults(textOut, textErr, exitCode)
				return
			}
			isDownloaded = true
			if args.Verbose {
				log.Println("downloaded muffet to: " + muffetPath)
			}

			// delete downloaded zip
			err = os.Remove(gzipFile)
			if err != nil {
				printCommandResults(textOut, textErr, exitCode)
				return
			}
			if args.Verbose {
				log.Println("deleted download bundle: " + gzipFile)
			}
		} else {
			log.Printf("error attempting to find 'muffet': %+v", err)
			return
		}
	} else {
		// muffet was found on the path, so use it
		muffetPath = muffetExec
	}

	return
}

func getTempDirWTrailingSlash() string {
	tempDir := os.TempDir()
	if !strings.HasSuffix(tempDir, "/") {
		tempDir = tempDir + "/"
	}
	return tempDir
}

func getExtractedExecutableName(baseName string) string {
	var extractedExecutableName string
	if //goland:noinspection GoBoolExpressions
	runtime.GOOS == "windows" {
		extractedExecutableName = baseName + ".exe"
	} else {
		extractedExecutableName = baseName
	}
	return extractedExecutableName
}

func printCommandResults(textOut string, textErr string, exitCode int) {
	log.Println("Out:\n" + textOut)
	log.Println("Err:\n" + textErr)
	log.Println("ExitCode:\n" + strconv.Itoa(exitCode))
}
