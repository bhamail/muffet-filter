package main

import (
	"encoding/json"
	"errors"
)

type UrlSuccessLink struct {
	Url    string `json:"url"`
	Status int    `json:"status"`
}
type UrlErrorLink struct {
	Url   string `json:"url"`
	Error string `json:"error"`
}
type UrlToCheck struct {
	Url   string        `json:"url"`
	Links []interface{} `json:"links"`
}
type Report struct {
	UrlsToCheck []UrlToCheck
}

type parseResponse struct {
	rawdata string
}

func (r *parseResponse) loadReport() (Report, error) {
	var report Report

	var raw []json.RawMessage
	if err := json.Unmarshal([]byte(r.rawdata), &raw); err != nil {
		return report, err
	}

	for i := 0; i < len(raw); i++ {
		var urlToCheck UrlToCheck
		if err := json.Unmarshal(raw[i], &urlToCheck); err != nil {
			return report, err
		}
		// convert untyped Links map to specific link types
		for j := 0; j < len(urlToCheck.Links); j++ {
			rawLink := urlToCheck.Links[j]

			var jsonLink []byte
			if jsonRaw, err := json.Marshal(rawLink); err != nil {
				return report, err
			} else {
				jsonLink = jsonRaw
			}

			if mapLink, ok := rawLink.(map[string]interface{}); !ok {
				return report, errors.New("invalid link map")
			} else if _, ok := mapLink["status"]; ok {
				// must be a Success link
				var urlSuccessLink UrlSuccessLink
				if err := json.Unmarshal(jsonLink, &urlSuccessLink); err == nil {
					// replace link in report struct
					urlToCheck.Links[j] = urlSuccessLink
				} else {
					return report, err
				}
			} else {
				// try using UrlErrorLink
				var urlErrorLink UrlErrorLink
				if err := json.Unmarshal(jsonLink, &urlErrorLink); err == nil {
					// replace link in report struct
					urlToCheck.Links[j] = urlErrorLink
					continue
				} else {
					return report, err
				}
			}
		}
		report.UrlsToCheck = append(report.UrlsToCheck, urlToCheck)
	}

	return report, nil
}
