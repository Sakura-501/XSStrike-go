package scan

import (
	"errors"
	"sort"
	"strings"

	"github.com/Sakura-501/XSStrike-go/internal/config"
	"github.com/Sakura-501/XSStrike-go/internal/dom"
	"github.com/Sakura-501/XSStrike-go/internal/encoder"
	"github.com/Sakura-501/XSStrike-go/internal/requester"
	"github.com/Sakura-501/XSStrike-go/internal/utils"
)

type ParamResult struct {
	Name        string `json:"name"`
	Payload     string `json:"payload"`
	Reflections int    `json:"reflections"`
	Reflected   bool   `json:"reflected"`
	Error       string `json:"error,omitempty"`
}

type Report struct {
	Target      string        `json:"target"`
	Method      string        `json:"method"`
	Tested      int           `json:"tested"`
	Reflected   int           `json:"reflected"`
	Findings    []ParamResult `json:"findings"`
	NoParams    bool          `json:"no_params"`
	RequestBase string        `json:"request_base"`
	DOM         dom.Report    `json:"dom"`
}

type Runner struct {
	Client *requester.Client
}

func NewRunner(client *requester.Client) *Runner {
	return &Runner{Client: client}
}

func (r *Runner) Run(target string, data string, headers map[string]string, jsonData bool, encode string) (*Report, error) {
	if r.Client == nil {
		return nil, errors.New("nil requester client")
	}

	isGET := strings.TrimSpace(data) == ""
	method := "GET"
	if !isGET {
		method = "POST"
	}

	params := utils.ParseParams(target, data, jsonData)
	normalizedTarget := normalizeTarget(r.Client, target, params, headers, isGET, jsonData)
	base := utils.GetURL(normalizedTarget, isGET)
	report := &Report{Target: normalizedTarget, Method: method, RequestBase: base, DOM: dom.Report{Checked: true, Findings: []dom.Finding{}}}
	if len(params) == 0 {
		report.NoParams = true
		return report, nil
	}

	if domResp, err := baselineResponse(r.Client, base, params, headers, isGET, jsonData); err == nil {
		report.DOM = dom.Analyze(domResp.Body)
	}

	keys := sortedKeys(params)
	for _, name := range keys {
		payload := encoder.Apply(encode, config.XSSChecker)
		current := cloneMap(params)
		current[name] = payload

		var (
			resp *requester.Response
			err  error
		)
		if isGET {
			resp, err = r.Client.DoGet(base, current, headers)
		} else {
			resp, err = r.Client.DoPost(base, current, headers, jsonData)
		}

		entry := ParamResult{Name: name, Payload: payload}
		if err != nil {
			entry.Error = err.Error()
			report.Findings = append(report.Findings, entry)
			report.Tested++
			continue
		}

		plainCount := strings.Count(resp.Body, config.XSSChecker)
		payloadCount := strings.Count(resp.Body, payload)
		if payloadCount > plainCount {
			entry.Reflections = payloadCount
		} else {
			entry.Reflections = plainCount
		}
		entry.Reflected = entry.Reflections > 0
		if entry.Reflected {
			report.Reflected++
		}
		report.Findings = append(report.Findings, entry)
		report.Tested++
	}

	return report, nil
}

func baselineResponse(client *requester.Client, base string, params map[string]string, headers map[string]string, isGET bool, jsonData bool) (*requester.Response, error) {
	if isGET {
		return client.DoGet(base, params, headers)
	}
	return client.DoPost(base, params, headers, jsonData)
}

func cloneMap(in map[string]string) map[string]string {
	out := map[string]string{}
	for key, value := range in {
		out[key] = value
	}
	return out
}

func sortedKeys(in map[string]string) []string {
	keys := make([]string, 0, len(in))
	for key := range in {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
