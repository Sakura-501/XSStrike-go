package scan

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sakura-501/XSStrike-go/internal/requester"
)

func TestRunWAFDetectionIntegration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("xss") != "" {
			w.Header().Set("cf-ray", "edge")
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte("Attention Required! | Cloudflare"))
			return
		}
		_, _ = w.Write([]byte("q=" + r.URL.Query().Get("q")))
	}))
	defer server.Close()

	runner := NewRunner(requester.New(requester.Config{TimeoutSeconds: 5}))
	report, err := runner.Run(server.URL+"?q=1", "", map[string]string{}, false, false, "")
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if !report.WAF.Detected {
		t.Fatalf("expected waf detection")
	}
	if report.WAF.Name == "" {
		t.Fatalf("expected waf name")
	}
}
