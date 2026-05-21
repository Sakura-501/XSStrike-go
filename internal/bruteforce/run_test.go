package bruteforce

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

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
		[]string{"AAA", "BBB"},
		"",
	)
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if report.NoParams {
		t.Fatalf("expected params")
	}
	if len(report.Hits) != 2 {
		t.Fatalf("expected all reflected payloads, got %d", len(report.Hits))
	}
	if report.Hits[0].Payload != "AAA" || report.Hits[1].Payload != "BBB" {
		t.Fatalf("expected deterministic payload order, got %+v", report.Hits)
	}
}

func TestRunNoParams(t *testing.T) {
	report, err := Run(
		requester.New(requester.Config{TimeoutSeconds: 5}),
		"https://example.com",
		"",
		false,
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
		[]string{"PAYLOAD"},
		"",
	)
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if report.NoParams {
		t.Fatalf("expected path params")
	}
	if report.Tested != 2 {
		t.Fatalf("expected two path entries tested, got %d", report.Tested)
	}
}

func TestRunWithConfigUsesParallelWorkers(t *testing.T) {
	var mu sync.Mutex
	active := 0
	maxActive := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		active++
		if active > maxActive {
			maxActive = active
		}
		mu.Unlock()

		time.Sleep(50 * time.Millisecond)
		_, _ = w.Write([]byte(fmt.Sprintf("q=%s", r.URL.Query().Get("q"))))

		mu.Lock()
		active--
		mu.Unlock()
	}))
	defer server.Close()

	report, err := RunWithConfig(
		requester.New(requester.Config{TimeoutSeconds: 5}),
		server.URL+"?q=1",
		"",
		false,
		false,
		map[string]string{},
		[]string{"A", "B", "C", "D", "E", "F", "G", "H"},
		"",
		Config{Threads: 4},
	)
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if report.Tested != 8 {
		t.Fatalf("expected tested=8, got %d", report.Tested)
	}
	if len(report.Hits) != 8 {
		t.Fatalf("expected hits=8, got %d", len(report.Hits))
	}

	mu.Lock()
	gotMaxActive := maxActive
	mu.Unlock()
	if gotMaxActive < 2 {
		t.Fatalf("expected parallel requests, max active=%d", gotMaxActive)
	}
}
