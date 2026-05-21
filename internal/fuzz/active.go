package fuzz

import (
	"sort"
	"strings"
	"sync"

	"github.com/Sakura-501/XSStrike-go/internal/encoder"
	"github.com/Sakura-501/XSStrike-go/internal/requester"
	"github.com/Sakura-501/XSStrike-go/internal/utils"
)

type Config struct {
	Threads int
}

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
	return RunWithConfig(client, target, data, jsonData, pathMode, headers, payloads, encodeMode, Config{Threads: 1})
}

func RunWithConfig(client *requester.Client, target string, data string, jsonData bool, pathMode bool, headers map[string]string, payloads []string, encodeMode string, cfg Config) (Report, error) {
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

	tasks := make([]task, 0, len(params)*len(payloads))
	for _, param := range sortedKeys(params) {
		for _, payload := range payloads {
			tasks = append(tasks, task{Index: len(tasks), Param: param, Payload: payload})
		}
	}

	if len(tasks) == 0 {
		return report, nil
	}

	threads := cfg.Threads
	if threads <= 0 {
		threads = 1
	}
	if threads > len(tasks) {
		threads = len(tasks)
	}

	results := make([]Entry, len(tasks))
	taskCh := make(chan task)
	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for currentTask := range taskCh {
				results[currentTask.Index] = runTask(client, base, params, headers, currentTask, isGET, jsonData, pathMode, encodeMode)
			}
		}()
	}
	for _, currentTask := range tasks {
		taskCh <- currentTask
	}
	close(taskCh)
	wg.Wait()

	report.Tested = len(results)
	report.Results = results
	for _, entry := range results {
		if entry.Reflected {
			report.Hits++
		}
	}
	return report, nil
}

type task struct {
	Index   int
	Param   string
	Payload string
}

func runTask(client *requester.Client, base string, params map[string]string, headers map[string]string, currentTask task, isGET bool, jsonData bool, pathMode bool, encodeMode string) Entry {
	entry := Entry{Param: currentTask.Param, Payload: currentTask.Payload}
	current := cloneMap(params)
	encodedPayload := encoder.Apply(encodeMode, currentTask.Payload)
	current[currentTask.Param] = encodedPayload

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
		return entry
	}

	count := strings.Count(resp.Body, currentTask.Payload)
	if encodedPayload != currentTask.Payload {
		encodedCount := strings.Count(resp.Body, encodedPayload)
		if encodedCount > count {
			count = encodedCount
		}
	}
	entry.Reflections = count
	entry.Reflected = count > 0
	return entry
}

func cloneMap(in map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range in {
		out[k] = v
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
