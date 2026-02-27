package scan

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Sakura-501/XSStrike-go/internal/requester"
)

func TestNormalizeTargetKeepsScheme(t *testing.T) {
	client := requester.New(requester.Config{TimeoutSeconds: 5})
	out := normalizeTarget(client, "https://example.com/path", map[string]string{}, map[string]string{}, true, false)
	if out != "https://example.com/path" {
		t.Fatalf("expected same target, got %q", out)
	}
}

func TestNormalizeTargetFallbackHTTP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	hostOnly := strings.TrimPrefix(server.URL, "http://")
	client := requester.New(requester.Config{TimeoutSeconds: 5})
	out := normalizeTarget(client, hostOnly, map[string]string{}, map[string]string{}, true, false)
	if !strings.HasPrefix(out, "http://") {
		t.Fatalf("expected http fallback, got %q", out)
	}
}
