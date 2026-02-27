package utils

import (
	"net/url"
	"regexp"
	"strings"
)

var (
	reScriptBody = regexp.MustCompile(`(?is)<script.*?>(.*?)</script>`)
	reScriptSrc  = regexp.MustCompile(`(?is)<(?:script).*?(?:src)=([^\s>]+)`)
)

// ExtractReflectedScripts returns inline script blocks containing checker marker.
func ExtractReflectedScripts(response, checker string) []string {
	result := []string{}
	lowerChecker := strings.ToLower(checker)
	for _, match := range reScriptBody.FindAllStringSubmatch(strings.ToLower(response), -1) {
		if len(match) < 2 {
			continue
		}
		content := match[1]
		if strings.Contains(content, lowerChecker) {
			result = append(result, content)
		}
	}
	return result
}

// ExtractJSSources returns script src entries from HTML response.
func ExtractJSSources(response string) []string {
	result := []string{}
	matches := reScriptSrc.FindAllStringSubmatch(response, -1)
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		src := strings.Trim(match[1], "\"'`")
		if src != "" {
			result = append(result, src)
		}
	}
	return result
}

// HandleAnchor normalizes relative links against the parent URL.
func HandleAnchor(parentURL, anchor string) string {
	parent, err := url.Parse(parentURL)
	if err != nil {
		return anchor
	}

	if strings.HasPrefix(anchor, "http") {
		return anchor
	}

	if strings.HasPrefix(anchor, "//") {
		return parent.Scheme + ":" + anchor
	}

	resolved, err := parent.Parse(anchor)
	if err != nil {
		return anchor
	}
	return resolved.String()
}
