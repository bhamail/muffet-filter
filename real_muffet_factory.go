package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

var defaultOptions = []string{"--buffer-size=8192", "--max-connections=10", "--color=always", "--format=json"}

type realMuffetFactory struct {
}

func newRealMuffetFactory() *realMuffetFactory {
	return &realMuffetFactory{}
}

func (f *realMuffetFactory) Create(options muffetOptions) muffetExecutor {
	muffetExec := realMuffetExecutor{options}
	return &muffetExec
}

type realMuffetExecutor struct {
	options muffetOptions
}

func (r *realMuffetExecutor) Check() (string, error) {
	cmd := exec.Command("/Users/bhamail/sonatype/sasq/link-checker/muffet/muffet", r.options.arguments...)
	cmdStdOut, err := cmd.StdoutPipe()
	cmdStdErr, err := cmd.StderrPipe()
	defer cmdStdOut.Close()
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
		os.Stdout.Write(line)
		os.Stdout.Write([]byte{'\n'})
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
		os.Stderr.Write(line)
		os.Stderr.Write([]byte{'\n'})
	}
	cmd.Wait()
	// we ignore exit code because failed links result in non-zero exit code
	fmt.Println(cmd.ProcessState.ExitCode())

	return jsonReportOut, nil
}
