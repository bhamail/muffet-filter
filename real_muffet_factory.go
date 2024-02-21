package main

import (
	"fmt"
	"log"
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
	isDownloaded, muffetPath, _ := getMuffet(args)
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

	textOut, _, exitCode, _ := executeCommand(args.Verbose, muffetPath, r.options.arguments...)
	// we ignore exit code because failed links result in non-zero exit code
	if args.Verbose {
		fmt.Printf("called muffet: %s, exit status: %d\n", muffetPath, exitCode)
	}

	return textOut, nil
}
