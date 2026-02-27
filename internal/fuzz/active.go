package fuzz

import (
	"strings"

	"github.com/Sakura-501/XSStrike-go/internal/encoder"
	"github.com/Sakura-501/XSStrike-go/internal/requester"
	"github.com/Sakura-501/XSStrike-go/internal/utils"
)

type Entry struct {
	Param       string `json:"param"`
	Payload     string `json:"payload"`
	Reflections int    `json:"reflections"`
	Reflected   bool   `json:"reflected"`
	Error       string `json:"error,omitempty"`
}

type Report struct {
	Target   string  `json:"target"`
	Tested   int     `json:"tested"`
	Hits     int     `json:"hits"`
	NoParams bool    `json:"no_params"`
	Results  []Entry `json:"results"`
}

func Run(client *requester.Client, target string, data string, jsonData bool, pathMode bool, headers map[string]string, payloads []string, encodeMode string) (Report, error) {
	report := Report{Target: target, Results: []Entry{}}
	if client == nil {
		return report, nil
	}

	isGET := strings.TrimSpace(data) == "" || pathMode
	base := utils.GetURL(target, isGET)
	params := map[string]string{}
	if pathMode {
		params = utils.URLPathToMap(target)
		base = target
	} else {
		params = utils.ParseParams(target, data, jsonData)
	}
	if len(params) == 0 {
		report.NoParams = true
		return report, nil
	}

	for param := range params {
		for _, payload := range payloads {
			report.Tested++
			entry := Entry{Param: param, Payload: payload}
			current := cloneMap(params)
			current[param] = encoder.Apply(encodeMode, payload)

			var (
				resp *requester.Response
				err  error
			)
			if pathMode {
				requestURL := utils.MapToURLPath(base, current)
				resp, err = client.DoGet(requestURL, map[string]string{}, headers)
			} else if isGET {
				resp, err = client.DoGet(base, current, headers)
			} else {
				resp, err = client.DoPost(base, current, headers, jsonData)
			}
			if err != nil {
				entry.Error = err.Error()
				report.Results = append(report.Results, entry)
				continue
			}

			count := strings.Count(resp.Body, payload)
			encodedPayload := encoder.Apply(encodeMode, payload)
			if encodedPayload != payload {
				encodedCount := strings.Count(resp.Body, encodedPayload)
				if encodedCount > count {
					count = encodedCount
				}
			}
			entry.Reflections = count
			entry.Reflected = count > 0
			if entry.Reflected {
				report.Hits++
			}
			report.Results = append(report.Results, entry)
		}
	}

	return report, nil
}

func cloneMap(in map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range in {
		out[k] = v
	}
	return out
}
