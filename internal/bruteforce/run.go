package bruteforce

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

func Run(client *requester.Client, target string, data string, jsonData bool, pathMode bool, headers map[string]string, payloads []string, encode string) (Report, error) {
	return RunWithConfig(client, target, data, jsonData, pathMode, headers, payloads, encode, Config{Threads: 1})
}

func RunWithConfig(client *requester.Client, target string, data string, jsonData bool, pathMode bool, headers map[string]string, payloads []string, encode string, cfg Config) (Report, error) {
	report := Report{Target: target, Hits: []Hit{}}
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

	results := make([]taskResult, len(tasks))
	taskCh := make(chan task)
	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for currentTask := range taskCh {
				results[currentTask.Index] = runTask(client, base, params, headers, currentTask, isGET, jsonData, pathMode, encode)
			}
		}()
	}
	for _, currentTask := range tasks {
		taskCh <- currentTask
	}
	close(taskCh)
	wg.Wait()

	report.Tested = len(results)
	for _, result := range results {
		if result.Hit {
			report.Hits = append(report.Hits, Hit{Param: result.Param, Payload: result.Payload, Reflections: result.Reflections})
		}
	}

	return report, nil
}

type task struct {
	Index   int
	Param   string
	Payload string
}

type taskResult struct {
	Param       string
	Payload     string
	Reflections int
	Hit         bool
}

func runTask(client *requester.Client, base string, params map[string]string, headers map[string]string, currentTask task, isGET bool, jsonData bool, pathMode bool, encodeMode string) taskResult {
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
		return taskResult{Param: currentTask.Param, Payload: currentTask.Payload}
	}

	countPayload := strings.Count(resp.Body, encodedPayload)
	countRaw := strings.Count(resp.Body, currentTask.Payload)
	if countPayload == 0 && countRaw == 0 {
		return taskResult{Param: currentTask.Param, Payload: currentTask.Payload}
	}
	count := countPayload
	if countRaw > count {
		count = countRaw
	}
	return taskResult{Param: currentTask.Param, Payload: currentTask.Payload, Reflections: count, Hit: true}
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
