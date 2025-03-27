package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
)

func newErrorForMissingField(fieldName, theStruct interface{}) error {
	return fmt.Errorf("missing required field: '%s' for type: %s, %+v", fieldName, reflect.TypeOf(theStruct).Name(), theStruct)
}

type UrlSuccessLink struct {
	Url    string `json:"url"`
	Status int    `json:"status"`
}

func (successLink *UrlSuccessLink) validate() error {
	// make sure required fields exist in the ErrorLink
	if successLink.Url == "" {
		return newErrorForMissingField("Url", *successLink)
	} else if successLink.Status == 0 {
		return newErrorForMissingField("Status", *successLink)
	}
	return nil
}

type UrlErrorLink struct {
	Url   string `json:"url"`
	Error string `json:"error"`
}

func (errorLink *UrlErrorLink) isMatch(linkPatternToIgnore UrlErrorLink) bool {
	if errorLink.Url != linkPatternToIgnore.Url {
		// check for regex match on url
		match, _ := regexp.MatchString(linkPatternToIgnore.Url, errorLink.Url)
		if !match {
			return false
		}
	}
	// if we got this far, the urls match, so now check the error message

	if errorLink.Error == linkPatternToIgnore.Error {
		return true
	}
	match, _ := regexp.MatchString(linkPatternToIgnore.Error, errorLink.Error)
	return match
}

func (errorLink *UrlErrorLink) validate() error {
	// make sure required fields exist in the ErrorLink
	if errorLink.Url == "" {
		return newErrorForMissingField("Url", *errorLink)
	} else if errorLink.Error == "" {
		return newErrorForMissingField("Error", *errorLink)
	}
	return nil
}

type UrlToCheck struct {
	Url   string        `json:"url"`
	Links []interface{} `json:"links"`
}
type Report struct {
	UrlsToCheck []UrlToCheck
}

//goland:noinspection SpellCheckingInspection
type parseResponse struct {
	rawdata string
}

func (r *parseResponse) loadReport(args *arguments) (report Report, err error) {

	var raw []json.RawMessage
	if err = json.Unmarshal([]byte(r.rawdata), &raw); err != nil {
		return
	}

	for i := 0; i < len(raw); i++ {
		var urlToCheck UrlToCheck
		if err = json.Unmarshal(raw[i], &urlToCheck); err != nil {
			return
		}
		// convert untyped Links map to specific link types
		for j := 0; j < len(urlToCheck.Links); j++ {
			rawLink := urlToCheck.Links[j]

			var jsonLink []byte
			if jsonLink, err = json.Marshal(rawLink); err != nil {
				return
			}

			if mapLink, ok := rawLink.(map[string]interface{}); !ok {
				err = errors.New("invalid link map")
				return
			} else if _, ok := mapLink["status"]; ok {
				// must be a Success link
				var urlSuccessLink UrlSuccessLink
				if err = json.Unmarshal(jsonLink, &urlSuccessLink); err == nil {
					// make sure required fields exist in the SuccessLink
					if err = urlSuccessLink.validate(); err != nil {
						return
					}

					// replace link in report struct
					urlToCheck.Links[j] = urlSuccessLink
				} else {
					return
				}
			} else {
				// try using UrlErrorLink
				var urlErrorLink UrlErrorLink
				if err = json.Unmarshal(jsonLink, &urlErrorLink); err == nil {
					// make sure required fields exist in the ErrorLink
					if err = urlErrorLink.validate(); err != nil {
						if args.Verbose {
							fmt.Printf("invalid error returned from muffet: %s, urlToCheck: %s\n", jsonLink, urlToCheck)
						}

						if args.IgnoreEmptyErrUrl {
							if urlErrorLink.Url == "" {
								urlErrorLink.Url = "empty"
								err = nil
								// don't continue, allow replacement of link in report struct below
							}
						} else {
							return
						}
					}

					// replace link in report struct
					urlToCheck.Links[j] = urlErrorLink
					continue
				} else {
					return
				}
			}
		}
		report.UrlsToCheck = append(report.UrlsToCheck, urlToCheck)
	}
	return
}

func (rep *Report) filter(errorsToIgnore []UrlErrorLink, isVerbose bool) (filteredReport Report, err error) {
	var tempUrlsToCheck []UrlToCheck
	for _, urlToCheck := range rep.UrlsToCheck {
		tempUrlToCheck := UrlToCheck{Url: urlToCheck.Url}
		for _, link := range urlToCheck.Links {
			switch v := link.(type) {
			case UrlErrorLink:
				if !isErrorIgnored(v, errorsToIgnore) {
					tempUrlToCheck.Links = append(tempUrlToCheck.Links, link)
				} else if isVerbose {
					fmt.Printf("skipping urlError: %+v on UrlToCheck: %s\n", link, urlToCheck.Url)
				}
			case UrlSuccessLink:
				// do nothing here, as we leave success links alone for now
				// maybe later we could decide to add a "quiet" mode, where success links get removed
				tempUrlToCheck.Links = append(tempUrlToCheck.Links, link)
			default:
				err = fmt.Errorf("unexpected url error type %T", v)
				return
			}
		}
		// add UrlToCheck if links exist
		if len(tempUrlToCheck.Links) > 0 {
			tempUrlsToCheck = append(tempUrlsToCheck, tempUrlToCheck)
		}
	}
	filteredReport = Report{UrlsToCheck: tempUrlsToCheck}
	return
}

func isErrorIgnored(urlError UrlErrorLink, errorsToIgnore []UrlErrorLink) bool {
	for _, errToIgnore := range errorsToIgnore {
		if urlError.isMatch(errToIgnore) {
			return true
		}
	}
	return false
}

func doesFileExist(fileToCheck string) (itExists bool, err error) {
	_, err = os.Stat(fileToCheck)
	if err == nil {
		itExists = true
	} else if errors.Is(err, os.ErrNotExist) {
		itExists = false
	}
	return
}

func loadIgnoreList(args *arguments) (ignoreUrlErrors []UrlErrorLink, err error) {
	var ignoreListFile string
	if args.IgnoresJson != "" {
		ignoreListFile = args.IgnoresJson
		var itExists bool
		if itExists, err = doesFileExist(ignoreListFile); !itExists {
			// a non-default file was specified, so it is an error if that specified file is missing
			return
		}
	} else {
		// next, we look for ignores file in the current working directory
		pwd, _ := os.Getwd()
		ignoreListFile = getDefaultIgnoresFile(pwd)
		var itExists bool
		if itExists, _ = doesFileExist(ignoreListFile); !itExists {
			// check user home dir for ignores file
			homeDir, _ := getUserHomeDir()
			ignoreListFile = getDefaultIgnoresFile(homeDir)
		}
	}

	var ignoreListRaw []byte
	ignoreListRaw, err = os.ReadFile(ignoreListFile)
	if err != nil {
		if args.Verbose {
			fmt.Printf("ignoring missing ignores file: %s", ignoreListFile)
		}
		err = nil
		return
	}

	err = json.Unmarshal(ignoreListRaw, &ignoreUrlErrors)
	if err != nil {
		fmt.Printf("error loading ignore list file: %s, error: %v", ignoreListFile, err)
	}
	return
}
