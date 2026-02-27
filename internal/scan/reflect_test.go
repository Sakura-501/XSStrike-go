package scan

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sakura-501/XSStrike-go/internal/requester"
)

func TestRunGETReflections(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		page := r.URL.Query().Get("page")
		_, _ = w.Write([]byte(fmt.Sprintf("q=%s,page=%s", q, page)))
	}))
	defer server.Close()

	runner := NewRunner(requester.New(requester.Config{TimeoutSeconds: 5}))
	report, err := runner.Run(server.URL+"?q=hello&page=1", "", map[string]string{}, false, "")
	if err != nil {
		t.Fatalf("unexpected run error: %v", err)
	}
	if report.Tested != 2 {
		t.Fatalf("unexpected tested count: %d", report.Tested)
	}
	if report.Reflected == 0 {
		t.Fatalf("expected at least one reflection")
	}
}

func TestRunPOSTJSONReflections(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"ref":"v3dm0s"}`))
	}))
	defer server.Close()

	runner := NewRunner(requester.New(requester.Config{TimeoutSeconds: 5}))
	report, err := runner.Run(server.URL, `{"name":"alice"}`, map[string]string{}, true, "")
	if err != nil {
		t.Fatalf("unexpected run error: %v", err)
	}
	if report.Method != "POST" {
		t.Fatalf("unexpected method: %s", report.Method)
	}
	if report.Reflected != 1 {
		t.Fatalf("expected reflected count 1, got %d", report.Reflected)
	}
}

func TestRunNoParams(t *testing.T) {
	runner := NewRunner(requester.New(requester.Config{TimeoutSeconds: 5}))
	report, err := runner.Run("https://example.com", "", map[string]string{}, false, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !report.NoParams {
		t.Fatalf("expected no params report")
	}
}
