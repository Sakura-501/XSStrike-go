package scan

import (
	"strings"

	"github.com/Sakura-501/XSStrike-go/internal/requester"
)

func normalizeTarget(client *requester.Client, target string, params map[string]string, headers map[string]string, isGET bool, jsonData bool) string {
	trimmed := strings.TrimSpace(target)
	if strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://") {
		return trimmed
	}

	httpsCandidate := "https://" + trimmed
	if _, err := baselineResponse(client, httpsCandidate, params, headers, isGET, jsonData); err == nil {
		return httpsCandidate
	}

	httpCandidate := "http://" + trimmed
	if _, err := baselineResponse(client, httpCandidate, params, headers, isGET, jsonData); err == nil {
		return httpCandidate
	}

	return httpCandidate
}
