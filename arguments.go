package main

import (
	"bytes"
	"fmt"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
)

func getUserHomeDir() (dirName string, err error) {
	dirName, err = os.UserHomeDir()
	if err != nil {
		log.Print(err)
	}
	return
}

const configDir = ".muffet-filter"
const ignoresFilename = "ignores.json"

var defaultIgnoresSuffix = configDir + "/" + ignoresFilename

func getDefaultIgnoresFile(prefix string) string {
	return prefix + "/" + defaultIgnoresSuffix
}

type arguments struct {
	MuffetPath        string   `short:"m" long:"muffet-path" description:"Path to muffet executable"`
	MuffetJson        string   `short:"j" long:"input-json" description:"Path to muffet link check output file in json format"`
	IgnoresJson       string   `short:"i" long:"ignores" description:"File containing url errors to ignore in json format. Defaults: .muffet-filter/ignores.json, ~/.muffet-filter/ignores.json"`
	Verbose           bool     `short:"v" long:"verbose" description:"Show more output"`
	Help              bool     `short:"h" long:"help" description:"Show this help"`
	Version           bool     `long:"version" description:"Show version"`
	MuffetArg         []string `long:"muffet-arg" description:"Additional argument passed to muffet executable."`
	IgnoreEmptyErrUrl bool     `long:"ignore-empty-err-url" description:"Ignore empty URL field in error links (only use for special cases)"`
	URL               string
}

func getArguments(ss []string) (*arguments, error) {
	args := arguments{}
	//ss, err := flags.NewParser(&args, flags.PassDoubleDash).ParseArgs(ss)
	ss, err := flags.ParseArgs(&args, ss)

	// TODO: Help arg parsing is borked. Check if replacing ParseArgs with NewParser fixes it.
	if err != nil {
		return nil, err
	} else if args.Version || args.Help {
		return &args, nil
	} else if len(ss) != 1 {
		return nil, fmt.Errorf("invalid number of arguments\n\n%s", help())
	}

	args.URL = ss[0]

	return &args, nil
}

func help() string {
	p := flags.NewParser(&arguments{}, flags.PassDoubleDash)
	p.Usage = "[options] <url of website to check>"

	// Parse() is run here to show default values in help.
	// This seems to be a bug in go-flags. Was this fixed???
	//_, _ = p.Parse() // nolint:errcheck

	b := &bytes.Buffer{}
	p.WriteHelp(b)
	return b.String()
}
