package main

import (
	"encoding/json"
	"fmt"
	"github.com/logrusorgru/aurora/v3"
	"io"
	"os"
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
		c.print(agentName, " version \"", version, "\" ", date)
		return true, nil
	}

	// call muffet to generate json response
	options := muffetOptions{arguments: defaultOptions}
	if len(args.MuffetArg) != 0 {
		options.arguments = append(options.arguments, args.MuffetArg...)
	}
	options.arguments = append(options.arguments, args.URL)
	muffetExec := c.factory.Create(options)
	jsonReport, err := muffetExec.Check(args)
	if err != nil {
		return false, err
	}

	// read json file into struct
	parseReport := parseResponse{jsonReport}
	report, err := parseReport.loadReport(args)
	if err != nil {
		return false, err
	}

	// filter out matching errors
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
