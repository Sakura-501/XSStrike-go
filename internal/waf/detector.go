package waf

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"github.com/Sakura-501/XSStrike-go/internal/requester"
)

type Signature struct {
	Code    string `json:"code"`
	Page    string `json:"page"`
	Headers string `json:"headers"`
}

type Detector struct {
	Signatures map[string]Signature
}

type Result struct {
	Detected bool   `json:"detected"`
	Name     string `json:"name,omitempty"`
	Score    int    `json:"score"`
}

func NewDefault() (*Detector, error) {
	candidates := defaultSignaturePaths()
	var lastErr error
	for _, path := range candidates {
		raw, err := os.ReadFile(path)
		if err != nil {
			lastErr = err
			continue
		}
		signs := map[string]Signature{}
		if err := json.Unmarshal(raw, &signs); err != nil {
			lastErr = err
			continue
		}
		return &Detector{Signatures: signs}, nil
	}
	if lastErr == nil {
		lastErr = errors.New("no signature path found")
	}
	return nil, fmt.Errorf("load waf signatures failed: %w", lastErr)
}

func (d *Detector) Detect(client *requester.Client, url string, params map[string]string, headers map[string]string, isGET bool, jsonData bool) Result {
	if d == nil || len(d.Signatures) == 0 || client == nil {
		return Result{}
	}

	probe := cloneMap(params)
	probe["xss"] = `<script>alert("XSS")</script>`

	var (
		resp *requester.Response
		err  error
	)
	if isGET {
		resp, err = client.DoGet(url, probe, headers)
	} else {
		resp, err = client.DoPost(url, probe, headers, jsonData)
	}
	if err != nil || resp == nil {
		return Result{}
	}

	status := fmt.Sprintf("%d", resp.StatusCode)
	if resp.StatusCode < 400 {
		return Result{}
	}

	headersText := headersToString(resp.Headers)
	page := resp.Body

	bestName := ""
	bestScore := 0
	for name, sign := range d.Signatures {
		score := 0
		if matches(sign.Page, page) {
			score += 2
		}
		if matches(sign.Code, status) {
			score += 1
		}
		if matches(sign.Headers, headersText) {
			score += 2
		}
		if score > bestScore {
			bestScore = score
			bestName = name
		}
	}

	if bestScore == 0 {
		return Result{}
	}
	return Result{Detected: true, Name: bestName, Score: bestScore}
}

func matches(pattern, target string) bool {
	if strings.TrimSpace(pattern) == "" {
		return false
	}
	re, err := regexp.Compile("(?i)" + pattern)
	if err != nil {
		return strings.Contains(strings.ToLower(target), strings.ToLower(pattern))
	}
	return re.FindStringIndex(target) != nil
}

func defaultSignaturePaths() []string {
	_, file, _, _ := runtime.Caller(0)
	base := filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
	return []string{
		filepath.Join(base, "db", "wafSignatures.json"),
		filepath.Join("db", "wafSignatures.json"),
	}
}

func cloneMap(in map[string]string) map[string]string {
	out := map[string]string{}
	for key, value := range in {
		out[key] = value
	}
	return out
}

func headersToString(in map[string]string) string {
	if len(in) == 0 {
		return ""
	}
	keys := make([]string, 0, len(in))
	for key := range in {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, key+": "+in[key])
	}
	return strings.Join(parts, "\n")
}
