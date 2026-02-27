package retirejs

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"github.com/Sakura-501/XSStrike-go/internal/requester"
	"github.com/Sakura-501/XSStrike-go/internal/utils"
)

type vulnerabilityDef struct {
	Below       string                 `json:"below"`
	AtOrAbove   string                 `json:"atOrAbove"`
	Severity    string                 `json:"severity"`
	Identifiers map[string]interface{} `json:"identifiers"`
	Info        []string               `json:"info"`
}

type componentDef struct {
	Vulnerabilities []vulnerabilityDef     `json:"vulnerabilities"`
	Extractors      map[string]interface{} `json:"extractors"`
}

type detection struct {
	Component string
	Version   string
	Method    string
}

type Vulnerability struct {
	Severity string                 `json:"severity,omitempty"`
	Info     []string               `json:"info,omitempty"`
	IDs      map[string]interface{} `json:"identifiers,omitempty"`
}

type Finding struct {
	Component       string          `json:"component"`
	Version         string          `json:"version"`
	Location        string          `json:"location"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
}

type Scanner struct {
	Definitions map[string]componentDef
	Checked     map[string]struct{}
}

func NewDefault() (*Scanner, error) {
	paths := defaultDefinitionPaths()
	var lastErr error
	for _, path := range paths {
		raw, err := os.ReadFile(path)
		if err != nil {
			lastErr = err
			continue
		}
		defs := map[string]componentDef{}
		if err := json.Unmarshal(raw, &defs); err != nil {
			lastErr = err
			continue
		}
		return &Scanner{Definitions: defs, Checked: map[string]struct{}{}}, nil
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("definition not found")
	}
	return nil, lastErr
}

func New(definitions map[string]componentDef) *Scanner {
	return &Scanner{Definitions: definitions, Checked: map[string]struct{}{}}
}

func (s *Scanner) ScanScript(uri string, content string) []Finding {
	if s == nil {
		return nil
	}

	results := []detection{}
	results = append(results, s.scanByExtractor(uri, "uri")...)
	results = append(results, s.scanByExtractor(fileName(uri), "filename")...)

	fromContent := s.scanByExtractor(content, "filecontent")
	if len(fromContent) == 0 {
		fromContent = s.scanByReplacement(content)
	}
	if len(fromContent) == 0 {
		fromContent = s.scanByHash(content)
	}
	results = append(results, fromContent...)

	findings := []Finding{}
	for _, result := range uniqueDetections(results) {
		comp, ok := s.Definitions[result.Component]
		if !ok {
			continue
		}
		vulns := []Vulnerability{}
		for _, v := range comp.Vulnerabilities {
			if isVulnerableVersion(result.Version, v.Below, v.AtOrAbove) {
				vulns = append(vulns, Vulnerability{Severity: v.Severity, Info: v.Info, IDs: v.Identifiers})
			}
		}
		if len(vulns) > 0 {
			findings = append(findings, Finding{Component: result.Component, Version: result.Version, Location: uri, Vulnerabilities: vulns})
		}
	}
	return findings
}

func (s *Scanner) ScanPage(client *requester.Client, pageURL string, html string, headers map[string]string) []Finding {
	if s == nil || client == nil {
		return nil
	}

	results := []Finding{}
	scripts := utils.ExtractJSSources(html)
	for _, script := range scripts {
		uri := utils.HandleAnchor(pageURL, script)
		if _, ok := s.Checked[uri]; ok {
			continue
		}
		s.Checked[uri] = struct{}{}
		resp, err := client.DoGet(uri, map[string]string{}, headers)
		if err != nil || resp == nil {
			continue
		}
		results = append(results, s.ScanScript(uri, resp.Body)...)
	}
	return results
}

func (s *Scanner) scanByExtractor(data string, extractor string) []detection {
	results := []detection{}
	for name, comp := range s.Definitions {
		raw, ok := comp.Extractors[extractor]
		if !ok {
			continue
		}
		list, ok := raw.([]interface{})
		if !ok {
			continue
		}
		for _, item := range list {
			pattern, ok := item.(string)
			if !ok {
				continue
			}
			if version := simpleMatch(pattern, data); version != "" {
				results = append(results, detection{Component: name, Version: version, Method: extractor})
			}
		}
	}
	return results
}

func (s *Scanner) scanByReplacement(data string) []detection {
	results := []detection{}
	for name, comp := range s.Definitions {
		raw, ok := comp.Extractors["filecontentreplace"]
		if !ok {
			continue
		}
		list, ok := raw.([]interface{})
		if !ok {
			continue
		}
		for _, item := range list {
			pattern, ok := item.(string)
			if !ok {
				continue
			}
			if version := replacementMatch(pattern, data); version != "" {
				results = append(results, detection{Component: name, Version: version, Method: "filecontentreplace"})
			}
		}
	}
	return results
}

func (s *Scanner) scanByHash(data string) []detection {
	results := []detection{}
	hash := sha1.Sum([]byte(data))
	hexHash := hex.EncodeToString(hash[:])
	for name, comp := range s.Definitions {
		raw, ok := comp.Extractors["hashes"]
		if !ok {
			continue
		}
		hashes, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if value, ok := hashes[hexHash]; ok {
			if version, ok := value.(string); ok {
				results = append(results, detection{Component: name, Version: version, Method: "hash"})
			}
		}
	}
	return results
}

func defaultDefinitionPaths() []string {
	_, file, _, _ := runtime.Caller(0)
	base := filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
	return []string{
		filepath.Join(base, "db", "definitions.json"),
		filepath.Join("db", "definitions.json"),
	}
}

func simpleMatch(pattern string, data string) string {
	pattern = deJSON(pattern)
	re, err := regexp.Compile(pattern)
	if err != nil {
		return ""
	}
	match := re.FindStringSubmatch(data)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func replacementMatch(pattern string, data string) string {
	pattern = deJSON(pattern)
	re := regexp.MustCompile(`^/(.*[^\\])/([^/]+)/$`)
	parts := re.FindStringSubmatch(pattern)
	if len(parts) < 3 {
		return ""
	}
	search, err := regexp.Compile("(" + parts[1] + ")")
	if err != nil {
		return ""
	}
	match := search.FindString(data)
	if match == "" {
		return ""
	}
	repl, err := regexp.Compile(parts[1])
	if err != nil {
		return ""
	}
	return repl.ReplaceAllString(match, parts[2])
}

func isVulnerableVersion(version string, below string, atOrAbove string) bool {
	if version == "" || below == "" {
		return false
	}
	if isAtOrAbove(version, below) {
		return false
	}
	if atOrAbove != "" && !isAtOrAbove(version, atOrAbove) {
		return false
	}
	return true
}

func isAtOrAbove(v1 string, v2 string) bool {
	a := splitVersion(v1)
	b := splitVersion(v2)
	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}
	for i := 0; i < maxLen; i++ {
		p1 := partAt(a, i)
		p2 := partAt(b, i)
		n1, ok1 := toInt(p1)
		n2, ok2 := toInt(p2)
		if ok1 && ok2 {
			if n1 > n2 {
				return true
			}
			if n1 < n2 {
				return false
			}
			continue
		}
		if p1 > p2 {
			return true
		}
		if p1 < p2 {
			return false
		}
	}
	return true
}

func splitVersion(v string) []string {
	re := regexp.MustCompile(`[.-]`)
	parts := re.Split(v, -1)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p == "" {
			continue
		}
		out = append(out, p)
	}
	return out
}

func partAt(parts []string, idx int) string {
	if idx >= len(parts) {
		return "0"
	}
	return parts[idx]
}

func toInt(v string) (int, bool) {
	for _, ch := range v {
		if ch < '0' || ch > '9' {
			return 0, false
		}
	}
	if v == "" {
		return 0, false
	}
	n := 0
	for _, ch := range v {
		n = n*10 + int(ch-'0')
	}
	return n, true
}

func uniqueDetections(in []detection) []detection {
	seen := map[string]struct{}{}
	out := make([]detection, 0, len(in))
	for _, d := range in {
		key := d.Component + "|" + d.Version
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, d)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Component == out[j].Component {
			return out[i].Version < out[j].Version
		}
		return out[i].Component < out[j].Component
	})
	return out
}

func deJSON(data string) string {
	return strings.ReplaceAll(data, `\\`, `\`)
}

func fileName(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		parts := strings.Split(uri, "/")
		if len(parts) == 0 {
			return uri
		}
		return parts[len(parts)-1]
	}
	parts := strings.Split(u.Path, "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}
