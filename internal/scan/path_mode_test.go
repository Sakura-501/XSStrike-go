package scan

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Sakura-501/XSStrike-go/internal/config"
	"github.com/Sakura-501/XSStrike-go/internal/requester"
)

func TestRunPathMode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("path=" + r.URL.Path))
	}))
	defer server.Close()

	runner := NewRunner(requester.New(requester.Config{TimeoutSeconds: 5}))
	report, err := runner.Run(server.URL+"/a/b", "", map[string]string{}, false, true, "")
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if report.Method != "GET" {
		t.Fatalf("expected GET for path mode")
	}
	if report.Tested != 2 {
		t.Fatalf("expected two tested path segments, got %d", report.Tested)
	}
	if report.Reflected == 0 {
		t.Fatalf("expected reflected path payload")
	}
	found := false
	for _, item := range report.Findings {
		if item.Reflected && strings.Contains(item.Payload, config.XSSChecker) {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected reflected finding with marker payload")
	}
}
