package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestLoadReportBlankIsError(t *testing.T) {
	resp := parseResponse{""}
	report, err := resp.loadReport()
	assert.Error(t, err)
	assert.Equal(t, Report{}, report)
}
func TestLoadReportEmpty(t *testing.T) {
	resp := parseResponse{"{}"}
	report, err := resp.loadReport()
	assert.Error(t, err)
	assert.Equal(t, Report{}, report)
}

const urlToCheckUrl string = "https://help.sonatype.com/index.html"
const urlErrorLinkUrl string = "https://help.sonatype.com/index.html#content-wrapper"
const urlErrorLinkError string = "id #content-wrapper not found"

const jsonUrlErrorLInk = `{
    "url": "` + urlToCheckUrl + `",
    "links": [
      {
        "url": "` + urlErrorLinkUrl + `",
        "error": "` + urlErrorLinkError + `"
      }
    ]
  }`

func TestLoadUrlErrorLink(t *testing.T) {
	var urlToCheck UrlToCheck
	err := json.Unmarshal([]byte(jsonUrlErrorLInk), &urlToCheck)
	assert.Nil(t, err)
	assert.Equal(t, UrlToCheck{Url: urlToCheckUrl,
		Links: []interface{}{
			map[string]interface{}{"url": urlErrorLinkUrl, "error": urlErrorLinkError},
		}},
		urlToCheck)
}

func TestUrlSuccessLinkValidate(t *testing.T) {
	{
		badLink := UrlSuccessLink{Status: -1}
		assert.EqualError(t, badLink.validate(), newErrorForMissingField("Url", badLink).Error())
	}

	{
		badLink := UrlSuccessLink{Url: "myUrl"}
		assert.EqualError(t, badLink.validate(), newErrorForMissingField("Status", badLink).Error())
	}

	{
		badLink := UrlSuccessLink{}
		assert.EqualError(t, badLink.validate(), newErrorForMissingField("Url", badLink).Error())
	}
}
func TestUrlErrorLinkValidate(t *testing.T) {
	{
		badLink := UrlErrorLink{Error: "myError"}
		assert.EqualError(t, badLink.validate(), newErrorForMissingField("Url", badLink).Error())
	}

	{
		badLink := UrlErrorLink{Url: "myUrl"}
		assert.EqualError(t, badLink.validate(), newErrorForMissingField("Error", badLink).Error())
	}

	{
		badLink := UrlErrorLink{}
		assert.EqualError(t, badLink.validate(), newErrorForMissingField("Url", badLink).Error())
	}
}

const jsonReportOneError = `[
  ` + jsonUrlErrorLInk + `
]`

var expectedFirstUrlToCheckError = UrlToCheck{
	Url: urlToCheckUrl,
	Links: []interface{}{
		UrlErrorLink{
			Url:   urlErrorLinkUrl,
			Error: urlErrorLinkError,
		}},
}

func TestLoadReportErrorLinkValue(t *testing.T) {
	jsonUrlErrorLInkBadVal := `[{
    "url": "` + urlToCheckUrl + `",
    "links": [
      {
        "urlBad": "` + urlErrorLinkUrl + `",
        "errorBad": "` + urlErrorLinkError + `"
      }
    ]
  }]`
	resp := parseResponse{jsonUrlErrorLInkBadVal}
	report, err := resp.loadReport()
	assert.EqualError(t, err, "missing required field: 'Url' for type: UrlErrorLink, {Url: Error:}")
	assert.Equal(t, Report{}, report)
}
func TestLoadReportOneError(t *testing.T) {
	resp := parseResponse{jsonReportOneError}
	report, err := resp.loadReport()
	assert.Nil(t, err)
	assert.Equal(t, Report{UrlsToCheck: []UrlToCheck{expectedFirstUrlToCheckError}}, report)
}

func TestLoadReportErrorParsingUrlToCheck(t *testing.T) {
	resp := parseResponse{`[{"url":9}]`}
	report, err := resp.loadReport()
	assert.EqualError(t, err, "json: cannot unmarshal number into Go struct field UrlToCheck.url of type string")
	assert.Equal(t, Report{}, report)
}

func TestLoadReportEmptyLinks(t *testing.T) {
	resp := parseResponse{`[{"url":"myUrl", "links": [{}]}]`}
	report, err := resp.loadReport()
	assert.EqualError(t, err, newErrorForMissingField("Url", UrlErrorLink{}).Error())
	assert.Equal(t, Report{}, report)
}

var expectedNearLast159UrlErrorToCheck = UrlToCheck{
	Url: "https://help.sonatype.com/en/nexus-repository-3-37-0---3-37-3-release-notes.html",
	Links: []interface{}{
		UrlErrorLink{
			Url:   "https://ossindex.sonatype.org/vulnerability/f0ac54b6-9b81-45bb-99a4-e6cb54749f9d",
			Error: "404",
		}},
}

func TestLoadReportBigErrorsOnly(t *testing.T) {
	// NOTE: This file was generated via muffet using:
	// $ ./muffet --buffer-size=8192 --max-connections=10 --color=always --format=json https://help.sonatype.com  > reportErrorsOnly.json
	err, report := loadTestReportFromFile(t, "testdata/reportErrorsOnly.json")
	assert.Nil(t, err)

	assert.NotNil(t, report)
	assert.Equal(t, 162, len(report.UrlsToCheck))
	assert.Equal(t, expectedFirstUrlToCheckError, report.UrlsToCheck[0])
	assert.Equal(t, expectedNearLast159UrlErrorToCheck, report.UrlsToCheck[159])
}

func loadTestReportFromFile(t *testing.T, filePath string) (error, Report) {
	fileContent, err := os.ReadFile(filePath)
	assert.Nil(t, err)
	// Convert []byte to string
	bigReport := string(fileContent)

	resp := parseResponse{bigReport}
	report, err := resp.loadReport()
	return err, report
}

const jsonUrlSuccessLInk = `{
    "url": "https://help.sonatype.com/index.html",
    "links": [
      {
        "url": "https://help.sonatype.com/favicon.ico",
        "status": 200
      }
    ]
  }`

func TestLoadUrlSuccessLink(t *testing.T) {
	var urlToCheck UrlToCheck
	err := json.Unmarshal([]byte(jsonUrlSuccessLInk), &urlToCheck)
	assert.Nil(t, err)
	assert.Equal(t, UrlToCheck{Url: "https://help.sonatype.com/index.html",
		Links: []interface{}{
			map[string]interface{}{"url": "https://help.sonatype.com/favicon.ico", "status": 200.0},
		}},
		urlToCheck)
}

const jsonReportOneSuccess = `[
  ` + jsonUrlSuccessLInk + `
]`

var expectedFirstUrlSuccessToCheck = UrlToCheck{
	Url: "https://help.sonatype.com/index.html",
	Links: []interface{}{
		UrlSuccessLink{
			Url:    "https://help.sonatype.com/favicon.ico",
			Status: 200.0,
		}},
}

func TestLoadReportOneSuccess(t *testing.T) {
	resp := parseResponse{jsonReportOneSuccess}
	report, err := resp.loadReport()
	assert.Nil(t, err)
	assert.Equal(t, Report{UrlsToCheck: []UrlToCheck{expectedFirstUrlSuccessToCheck}}, report)
}

func TestLoadReportBigSuccessOnly(t *testing.T) {
	err, report := loadTestReportFromFile(t, "testdata/reportSuccessOnly.json")
	assert.Nil(t, err)

	assert.NotNil(t, report)
	assert.Equal(t, 1, len(report.UrlsToCheck))
	assert.Equal(t, 27, len(report.UrlsToCheck[0].Links))
	assert.Equal(t, "https://www.google.com/", report.UrlsToCheck[0].Url)
}
func TestLoadReportBigSuccessAndError(t *testing.T) {
	err, report := loadTestReportFromFile(t, "testdata/reportSuccessAndError.json")
	assert.Nil(t, err)

	assert.NotNil(t, report)
	assert.Equal(t, 1, len(report.UrlsToCheck))
	assert.Equal(t, 72, len(report.UrlsToCheck[0].Links))
	assert.Equal(t, UrlSuccessLink{Url: "https://help.sonatype.com/css/sm-simple.css", Status: 200}, report.UrlsToCheck[0].Links[0])
	assert.Equal(t, UrlErrorLink{Url: "https://help.sonatype.com/index.html#content-wrapper", Error: "id #content-wrapper not found"}, report.UrlsToCheck[0].Links[71])
}

func TestUrlErrorIsMatch(t *testing.T) {
	errLink := UrlErrorLink{"a", "b"}
	assert.Equal(t, false, errLink.isMatch(UrlErrorLink{"x", "y"}))
	assert.Equal(t, false, errLink.isMatch(UrlErrorLink{"a", "y"}))
	assert.Equal(t, false, errLink.isMatch(UrlErrorLink{"x", "b"}))
	assert.Equal(t, true, errLink.isMatch(UrlErrorLink{"a", "b"}))
}

func TestLoadIgnoreListFromTestdata(t *testing.T) {
	args := arguments{IgnoresJson: "testdata/urlErrorIgnore.json"}
	ignores, err := loadIgnoreList(&args)
	assert.Nil(t, err)
	assert.NotNil(t, ignores)
}
func TestLoadIgnoreListBadArg(t *testing.T) {
	args := arguments{IgnoresJson: "bad-ignore-file.json"}
	ignores, err := loadIgnoreList(&args)
	assert.EqualError(t, err, "stat bad-ignore-file.json: no such file or directory")
	assert.Nil(t, ignores)
}
func TestLoadIgnoreListMissingAllDefaults(t *testing.T) {
	// override default ignores file to non-existent file/path
	origDefaultIgnoresSuffix := defaultIgnoresSuffix
	defer func() {
		defaultIgnoresSuffix = origDefaultIgnoresSuffix
	}()
	defaultIgnoresSuffix = "bogusIgnoresTestPathSuffix"

	args := arguments{Verbose: true}
	ignores, err := loadIgnoreList(&args)
	assert.Nil(t, err)
	assert.Nil(t, ignores)
}
func TestLoadIgnoreListDefaultCurrentDir(t *testing.T) {
	// override default ignores file to non-existent file/path
	origDefaultIgnoresSuffix := defaultIgnoresSuffix
	defer func() {
		defaultIgnoresSuffix = origDefaultIgnoresSuffix
	}()
	defaultIgnoresSuffix = "testdata/urlErrorIgnore.json"

	args := arguments{Verbose: true}
	ignores, err := loadIgnoreList(&args)
	assert.Nil(t, err)
	assert.NotNil(t, ignores)
}
func TestLoadIgnoreListInvalidDefault(t *testing.T) {
	// override default ignores file to non-existent file/path
	origDefaultIgnoresSuffix := defaultIgnoresSuffix
	defer func() {
		defaultIgnoresSuffix = origDefaultIgnoresSuffix
	}()
	pwd, _ := os.Getwd()
	userHome, _ := getUserHomeDir()
	indexProjectPath := strings.Index(pwd, userHome)
	if indexProjectPath == -1 {
		t.Skipf("skip test: %s, because our project is not under user home directory", t.Name())
		return
	}
	projectPath := pwd[indexProjectPath+1+len(userHome):]
	badTestFile := projectPath + "/testdata/bad.json"
	defaultIgnoresSuffix = badTestFile

	args := arguments{Verbose: true}
	ignores, err := loadIgnoreList(&args)
	assert.EqualError(t, err, "json: cannot unmarshal string into Go value of type []main.UrlErrorLink")
	assert.Nil(t, ignores)
}

func TestReportFilterOneErrorNoMatch(t *testing.T) {
	resp := parseResponse{jsonReportOneError}
	report, err := resp.loadReport()
	assert.Nil(t, err)

	reportFiltered, err := report.filter(nil, false)
	assert.Nil(t, err)
	assert.Equal(t, report.UrlsToCheck[0], reportFiltered.UrlsToCheck[0])
	assert.Equal(t, 1, len(reportFiltered.UrlsToCheck[0].Links))
	assert.Equal(t, Report{UrlsToCheck: []UrlToCheck{expectedFirstUrlToCheckError}}, report)
}
func TestReportFilterOneErrorMatch(t *testing.T) {
	resp := parseResponse{jsonReportOneError}
	report, err := resp.loadReport()
	assert.Nil(t, err)

	reportFiltered, err := report.filter([]UrlErrorLink{
		{Url: "https://help.sonatype.com/index.html#content-wrapper", Error: "id #content-wrapper not found"},
	}, false)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(reportFiltered.UrlsToCheck))
	assert.Equal(t, 1, len(report.UrlsToCheck))
}
func TestReportFilterTwoErrorMatch(t *testing.T) {
	resp := parseResponse{jsonReportOneError}
	report, err := resp.loadReport()
	assert.Nil(t, err)

	keptErrLink := UrlErrorLink{"urlNoMatch", "errorNoMatch"}
	report.UrlsToCheck[0].Links = append(report.UrlsToCheck[0].Links, keptErrLink)

	reportFiltered, err := report.filter([]UrlErrorLink{
		{Url: "https://help.sonatype.com/index.html#content-wrapper", Error: "id #content-wrapper not found"},
	}, false)
	assert.Equal(t, 1, len(reportFiltered.UrlsToCheck[0].Links))
	assert.Equal(t, keptErrLink, reportFiltered.UrlsToCheck[0].Links[0])
	assert.Equal(t, keptErrLink, report.UrlsToCheck[0].Links[1])
}
func TestReportFilterErrorMatchAndSuccessLink(t *testing.T) {
	resp := parseResponse{jsonReportOneError}
	report, err := resp.loadReport()
	assert.Nil(t, err)

	keptSuccessLink := UrlSuccessLink{"urlSuccess", 200}
	report.UrlsToCheck[0].Links = append(report.UrlsToCheck[0].Links, keptSuccessLink)

	reportFiltered, err := report.filter([]UrlErrorLink{
		{Url: "https://help.sonatype.com/index.html#content-wrapper", Error: "id #content-wrapper not found"},
	}, false)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(reportFiltered.UrlsToCheck[0].Links))
	assert.Equal(t, keptSuccessLink, reportFiltered.UrlsToCheck[0].Links[0])
	assert.Equal(t, keptSuccessLink, report.UrlsToCheck[0].Links[1])
}
func TestReportFilterUnknownLinkInterface(t *testing.T) {
	resp := parseResponse{jsonReportOneError}
	report, err := resp.loadReport()
	assert.Nil(t, err)

	unknownLinkType := "unknown type"
	report.UrlsToCheck[0].Links = append(report.UrlsToCheck[0].Links, unknownLinkType)

	_, err = report.filter(nil, false)
	assert.EqualError(t, err, "I don't know about type string!\n")
}
