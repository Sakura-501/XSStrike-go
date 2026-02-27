package bruteforce

import (
	"strings"

	"github.com/Sakura-501/XSStrike-go/internal/encoder"
	"github.com/Sakura-501/XSStrike-go/internal/requester"
	"github.com/Sakura-501/XSStrike-go/internal/utils"
)

type Hit struct {
	Param       string `json:"param"`
	Payload     string `json:"payload"`
	Reflections int    `json:"reflections"`
}

type Report struct {
	Target   string `json:"target"`
	Tested   int    `json:"tested"`
	Hits     []Hit  `json:"hits"`
	NoParams bool   `json:"no_params"`
}

func Run(client *requester.Client, target string, data string, jsonData bool, headers map[string]string, payloads []string, encode string) (Report, error) {
	report := Report{Target: target, Hits: []Hit{}}
	if client == nil {
		return report, nil
	}

	isGET := strings.TrimSpace(data) == ""
	base := utils.GetURL(target, isGET)
	params := utils.ParseParams(target, data, jsonData)
	if len(params) == 0 {
		report.NoParams = true
		return report, nil
	}

	for param := range params {
		for _, payload := range payloads {
			report.Tested++
			current := cloneMap(params)
			encodedPayload := encoder.Apply(encode, payload)
			current[param] = encodedPayload

			var (
				resp *requester.Response
				err  error
			)
			if isGET {
				resp, err = client.DoGet(base, current, headers)
			} else {
				resp, err = client.DoPost(base, current, headers, jsonData)
			}
			if err != nil {
				continue
			}

			countPayload := strings.Count(resp.Body, encodedPayload)
			countRaw := strings.Count(resp.Body, payload)
			if countPayload > 0 || countRaw > 0 {
				count := countPayload
				if countRaw > count {
					count = countRaw
				}
				report.Hits = append(report.Hits, Hit{Param: param, Payload: payload, Reflections: count})
				break
			}
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
