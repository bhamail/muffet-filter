package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
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

const jsonUrlErrorLInk = `{
    "url": "https://help.sonatype.com/index.html",
    "links": [
      {
        "url": "https://help.sonatype.com/index.html#content-wrapper",
        "error": "id #content-wrapper not found"
      }
    ]
  }`

func TestLoadUrlErrorLink(t *testing.T) {
	var urlToCheck UrlToCheck
	err := json.Unmarshal([]byte(jsonUrlErrorLInk), &urlToCheck)
	assert.Nil(t, err)
	assert.Equal(t, UrlToCheck{Url: "https://help.sonatype.com/index.html",
		Links: []interface{}{
			map[string]interface{}{"url": "https://help.sonatype.com/index.html#content-wrapper", "error": "id #content-wrapper not found"},
		}},
		urlToCheck)
}

const jsonReportOneError = `[
  ` + jsonUrlErrorLInk + `
]`

var expectedFirstUrlErrorToCheck = UrlToCheck{
	Url: "https://help.sonatype.com/index.html",
	Links: []interface{}{
		UrlErrorLink{
			Url:   "https://help.sonatype.com/index.html#content-wrapper",
			Error: "id #content-wrapper not found",
		}},
}

func TestLoadReportOneError(t *testing.T) {
	resp := parseResponse{jsonReportOneError}
	report, err := resp.loadReport()
	assert.Nil(t, err)
	assert.Equal(t, Report{UrlsToCheck: []UrlToCheck{expectedFirstUrlErrorToCheck}}, report)
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
	assert.Nil(t, err, "json: cannot unmarshal number into Go struct field UrlToCheck.url of type string")
	assert.Equal(t, Report{UrlsToCheck: []UrlToCheck{{
		Url: "myUrl",
		Links: []interface{}{
			UrlErrorLink{
				Url:   "",
				Error: "",
			}},
	}}}, report)
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
	assert.Equal(t, expectedFirstUrlErrorToCheck, report.UrlsToCheck[0])
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
func TestReportFilterOneErrorNoMatch(t *testing.T) {
	resp := parseResponse{jsonReportOneError}
	report, err := resp.loadReport()
	assert.Nil(t, err)

	reportFiltered, err := report.filter(nil)
	assert.Equal(t, 1, len(reportFiltered.UrlsToCheck[0].Links))
	assert.Equal(t, Report{UrlsToCheck: []UrlToCheck{expectedFirstUrlErrorToCheck}}, report)
}
func TestReportFilterOneErrorMatch(t *testing.T) {
	resp := parseResponse{jsonReportOneError}
	report, err := resp.loadReport()
	assert.Nil(t, err)

	reportFiltered, err := report.filter([]UrlErrorLink{
		{Url: "https://help.sonatype.com/index.html#content-wrapper", Error: "id #content-wrapper not found"},
	})
	assert.Equal(t, 0, len(reportFiltered.UrlsToCheck[0].Links))
}
func TestReportFilterTwoErrorMatch(t *testing.T) {
	resp := parseResponse{jsonReportOneError}
	report, err := resp.loadReport()
	assert.Nil(t, err)

	keptErrLink := UrlErrorLink{"urlNoMatch", "errorNoMatch"}
	report.UrlsToCheck[0].Links = append(report.UrlsToCheck[0].Links, keptErrLink)

	reportFiltered, err := report.filter([]UrlErrorLink{
		{Url: "https://help.sonatype.com/index.html#content-wrapper", Error: "id #content-wrapper not found"},
	})
	assert.Equal(t, 1, len(reportFiltered.UrlsToCheck[0].Links))
	assert.Equal(t, keptErrLink, reportFiltered.UrlsToCheck[0].Links[0])
	assert.Equal(t, keptErrLink, report.UrlsToCheck[0].Links[0])
}
