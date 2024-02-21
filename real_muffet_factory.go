package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"slices"
)

var defaultOptions = []string{"--buffer-size=8192", "--max-connections=10", "--color=always", "--format=json"}

type realMuffetFactory struct {
}

func newRealMuffetFactory() *realMuffetFactory {
	return &realMuffetFactory{}
}

func (f *realMuffetFactory) Create(options muffetOptions) muffetExecutor {
	return &realMuffetExecutor{options}
}

type realMuffetExecutor struct {
	options muffetOptions
}

func (r *realMuffetExecutor) Check(args *arguments) (string, error) {
	isDownloaded, muffetPath, err := getMuffet(args)
	// todo Decide if we want to delete the downloaded executable here, maybe add a flag?
	/*
		if isDownloaded {
			defer func() {
				_ = os.Remove(muffetPath)
			}()
		}
	*/
	if isDownloaded {
		log.Println("muffet was downloaded to: " + muffetPath)
	}

	cmd := exec.Command(muffetPath, r.options.arguments...)
	cmdStdOut, err := cmd.StdoutPipe()
	defer func(cmdStdOut io.ReadCloser) {
		_ = cmdStdOut.Close()
	}(cmdStdOut)
	if err != nil {
		log.Fatalf("command failed with %s\n", err)
	}
	cmdStdErr, err := cmd.StderrPipe()
	defer func(cmdStdErr io.ReadCloser) {
		_ = cmdStdErr.Close()
	}(cmdStdErr)
	if err != nil {
		log.Fatalf("command failed with %s\n", err)
	}

	stdoutReader := bufio.NewReader(cmdStdOut)
	stderrReader := bufio.NewReader(cmdStdErr)
	err = cmd.Start()
	if err != nil {
		return "", err
	}
	// Read stdout
	// capture in string
	jsonReportOut := ""
	for {
		line, err := stdoutReader.ReadBytes('\n')
		if err == io.EOF {
			if len(line) == 0 {
				// remove last newline
				break
			}
		} else {
			if err != nil {
				log.Fatal(err)
			}
			line = line[:(len(line) - 1)]
		}
		jsonReportOut = jsonReportOut + string(line)
		if slices.Contains(r.options.arguments, "--verbose") || slices.Contains(r.options.arguments, "-v") {
			_, _ = os.Stdout.Write(line)
			_, _ = os.Stdout.Write([]byte{'\n'})
		}
	}
	// Read stderr
	for {
		line, err := stderrReader.ReadBytes('\n')
		if err == io.EOF {
			if len(line) == 0 {
				break
			}
		} else {
			if err != nil {
				log.Fatal(err)
			}
			line = line[:(len(line) - 1)]
		}
		_, _ = os.Stderr.Write(line)
		_, _ = os.Stderr.Write([]byte{'\n'})
	}
	// we ignore exit code because failed links result in non-zero exit code
	_ = cmd.Wait()
	if args.Verbose {
		fmt.Printf("called muffet: %s, exit status: %d\n", cmd.Path, cmd.ProcessState.ExitCode())
	}

	return jsonReportOut, nil
}
