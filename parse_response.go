package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

type UrlSuccessLink struct {
	Url    string `json:"url"`
	Status int    `json:"status"`
}
type UrlErrorLink struct {
	Url   string `json:"url"`
	Error string `json:"error"`
}

func (errorLink *UrlErrorLink) isMatch(link UrlErrorLink) bool {
	if (errorLink.Url == link.Url) && (errorLink.Error == link.Error) {
		return true
	}
	return false
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

func (rep *Report) filter(errorsToIgnore []UrlErrorLink) (Report, error) {
	//temp := s[:0]
	for indexUrl, urlToCheck := range rep.UrlsToCheck {
		for indexLink, link := range urlToCheck.Links {
			switch v := link.(type) {
			case UrlErrorLink:
				for _, errToIgnore := range errorsToIgnore {
					if v.isMatch(errToIgnore) {
						// remove this error link
						rep.UrlsToCheck[indexUrl].Links = append(urlToCheck.Links[:indexLink], urlToCheck.Links[indexLink+1:]...)
					}
				}
			default:
				return *rep, errors.New(fmt.Sprintf("I don't know about type %T!\n", v))
			}
		}
	}
	return *rep, nil
}
