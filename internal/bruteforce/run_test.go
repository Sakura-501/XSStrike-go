package bruteforce

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sakura-501/XSStrike-go/internal/requester"
)

func TestRunGET(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf("q=%s", r.URL.Query().Get("q"))))
	}))
	defer server.Close()

	report, err := Run(
		requester.New(requester.Config{TimeoutSeconds: 5}),
		server.URL+"?q=1",
		"",
		false,
		map[string]string{},
		[]string{"AAA", "BBB"},
		"",
	)
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if report.NoParams {
		t.Fatalf("expected params")
	}
	if len(report.Hits) != 1 {
		t.Fatalf("expected one hit, got %d", len(report.Hits))
	}
}

func TestRunNoParams(t *testing.T) {
	report, err := Run(
		requester.New(requester.Config{TimeoutSeconds: 5}),
		"https://example.com",
		"",
		false,
		map[string]string{},
		[]string{"AAA"},
		"",
	)
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if !report.NoParams {
		t.Fatalf("expected no params")
	}
}
