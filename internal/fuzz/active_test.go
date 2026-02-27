package fuzz

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
		false,
		map[string]string{},
		[]string{"A", "B"},
		"",
	)
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if report.Tested != 2 {
		t.Fatalf("expected tested=2, got %d", report.Tested)
	}
	if report.Hits == 0 {
		t.Fatalf("expected reflected hits")
	}
}

func TestRunPathMode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("path=" + r.URL.Path))
	}))
	defer server.Close()

	report, err := Run(
		requester.New(requester.Config{TimeoutSeconds: 5}),
		server.URL+"/a/b",
		"",
		false,
		true,
		map[string]string{},
		[]string{"ZZ"},
		"",
	)
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if report.NoParams {
		t.Fatalf("expected path params")
	}
	if report.Tested != 2 {
		t.Fatalf("expected tested=2 in path mode")
	}
}
