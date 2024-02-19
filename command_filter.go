package main

import (
	"encoding/json"
	"fmt"
	"github.com/logrusorgru/aurora/v3"
	"io"
	"os"
	"os/exec"
	"strings"
)

type commandFilter struct {
	stdout, stderr io.Writer
	terminal       bool
	factory        muffetFactory
}

func newCommandFilter(stdout, stderr io.Writer, terminal bool, f muffetFactory) *commandFilter {
	return &commandFilter{stdout, stderr, terminal, f}
}

func (c *commandFilter) getMuffet(args *arguments) (muffetPath string, err error) {
	var muffetExec string
	if args.MuffetPath != "" {
		muffetExec = args.MuffetPath
		var itExists bool
		if itExists, err = doesFileExist(muffetExec); !itExists {
			// a non-default file was specified, so it is an error if that specified file is missing
			return
		}
	}

	// fetch muffet if not on path
	cmd := exec.Command("muffet", "--version")
	err = cmd.Start()
	if err != nil {
		if err.Error() == "exec: \"muffet\": executable file not found in $PATH" {
			// try to fetch muffet to local temp dir
			// todo Install muffet
			/*
				# Install muffet if needed
				if ! hash muffet &> /dev/null; then

				if [ "${MUFFET_VERSION}" = "latest" ]; then
				MUFFET_URL=$(wget -qO- https://api.github.com/repos/raviqqe/muffet/releases/latest | grep "browser_download_url.*muffet_linux_amd64.tar.gz" | cut -d \" -f 4)
				else
				MUFFET_URL="https://github.com/raviqqe/muffet/releases/download/v${MUFFET_VERSION}/muffet_linux_amd64.tar.gz"
				fi

				wget -qO- "${MUFFET_URL}" | $sudo_cmd tar xzf - -C /usr/local/bin/ muffet
				fi
			*/
			cmd = exec.Command("go", "install", "muffet@latest")
			err = cmd.Start()
			if err != nil {
				return
			}
		}
	} else {
		muffetPath = cmd.Path
	}

	return
}

func (c *commandFilter) Run(args []string) bool {
	ok, err := c.runWithError(args)
	if err != nil {
		c.printError(err)
	}

	return ok
}

func (c *commandFilter) runWithError(ss []string) (bool, error) {
	args, err := getArguments(ss)
	if err != nil {
		return false, err
	} else if args.Help {
		c.print(help())
		return true, nil
	} else if args.Version {
		c.print(version)
		return true, nil
	}

	// call muffet to generate json response
	options := muffetOptions{arguments: defaultOptions}
	options.arguments = append(options.arguments, args.URL)
	muffetExec := c.factory.Create(options)
	jsonReport, err := muffetExec.Check(args)
	if err != nil {
		return false, err
	}

	// read json file into struct
	parseReport := parseResponse{jsonReport}
	report, err := parseReport.loadReport()
	if err != nil {
		return false, err
	}

	// filter out matching errors
	// todo Build this next
	// load errorsToIgnore from on disk config and/or args
	errorsToIgnore, err := loadIgnoreList(args)
	if err != nil {
		return false, err
	}

	reportFiltered, err := report.filter(errorsToIgnore, args.Verbose)
	if err != nil {
		return false, err
	}
	if len(reportFiltered.UrlsToCheck) > 0 {
		prettyJson, err := json.MarshalIndent(reportFiltered, "", "  ")
		if err != nil {
			return false, err
		}
		_, _ = os.Stdout.Write(prettyJson)
		fmt.Println()
		//fmt.Printf("%v\n", report)
		return false, nil
	}

	return true, nil
}

func (c *commandFilter) print(xs ...any) {
	if _, err := fmt.Fprintln(c.stdout, strings.TrimSpace(fmt.Sprint(xs...))); err != nil {
		panic(err)
	}
}

func (c *commandFilter) printError(xs ...any) {
	s := fmt.Sprint(xs...)

	// Do not check --color option here because this can be used on argument parsing errors.
	if c.terminal {
		s = aurora.Red(s).String()
	}

	if _, err := fmt.Fprintln(c.stderr, s); err != nil {
		panic(err)
	}
}
