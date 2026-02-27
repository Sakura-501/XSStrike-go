package crawl

import (
	"strings"

	"github.com/Sakura-501/XSStrike-go/internal/config"
	"github.com/Sakura-501/XSStrike-go/internal/requester"
)

type Finding struct {
	URL         string `json:"url"`
	Method      string `json:"method"`
	Param       string `json:"param"`
	Reflections int    `json:"reflections"`
	BlindSent   bool   `json:"blind_sent"`
	Error       string `json:"error,omitempty"`
}

type ScanSummary struct {
	Tested   int       `json:"tested"`
	Findings []Finding `json:"findings"`
}

func ScanForms(client *requester.Client, forms []Form, headers map[string]string, blindPayload string) ScanSummary {
	summary := ScanSummary{Findings: []Finding{}}
	if client == nil {
		return summary
	}

	for _, form := range forms {
		method := strings.ToLower(strings.TrimSpace(form.Method))
		if method == "" {
			method = "get"
		}

		baseParams := map[string]string{}
		for _, input := range form.Inputs {
			if strings.TrimSpace(input.Name) == "" {
				continue
			}
			baseParams[input.Name] = input.Value
		}
		if len(baseParams) == 0 {
			continue
		}

		for param := range baseParams {
			summary.Tested++
			current := cloneMap(baseParams)
			current[param] = config.XSSChecker

			entry := Finding{URL: form.Action, Method: method, Param: param}
			var (
				resp *requester.Response
				err  error
			)
			if method == "post" {
				resp, err = client.DoPost(form.Action, current, headers, false)
			} else {
				resp, err = client.DoGet(form.Action, current, headers)
			}
			if err != nil {
				entry.Error = err.Error()
				summary.Findings = append(summary.Findings, entry)
				continue
			}

			entry.Reflections = strings.Count(resp.Body, config.XSSChecker)
			if blindPayload != "" {
				blind := cloneMap(baseParams)
				blind[param] = blindPayload
				if method == "post" {
					_, _ = client.DoPost(form.Action, blind, headers, false)
				} else {
					_, _ = client.DoGet(form.Action, blind, headers)
				}
				entry.BlindSent = true
			}

			if entry.Reflections > 0 || entry.BlindSent {
				summary.Findings = append(summary.Findings, entry)
			}
		}
	}

	return summary
}

func cloneMap(in map[string]string) map[string]string {
	out := map[string]string{}
	for key, value := range in {
		out[key] = value
	}
	return out
}
