package reflection

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sakura-501/XSStrike-go/internal/config"
	"github.com/Sakura-501/XSStrike-go/internal/requester"
)

func TestCheck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf("x=%s", r.URL.Query().Get("x"))))
	}))
	defer server.Close()

	client := requester.New(requester.Config{TimeoutSeconds: 5})
	params := map[string]string{"x": config.XSSChecker}
	eff := Check(client, server.URL, params, map[string]string{}, true, false, "<", []int{1}, "")
	if len(eff) == 0 {
		t.Fatalf("expected efficiencies")
	}
	if eff[0] == 0 {
		t.Fatalf("expected non-zero efficiency")
	}
}

func TestFilterCheck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf("<a href=\"%s\">x</a>", r.URL.Query().Get("x"))))
	}))
	defer server.Close()

	client := requester.New(requester.Config{TimeoutSeconds: 5})
	occ := Occurrences{
		10: &Occurrence{Position: 10, Context: "attribute", Details: Details{Name: "href", Type: "value", Quote: "\""}, Score: map[string]int{}},
	}
	params := map[string]string{"x": config.XSSChecker}
	out := FilterCheck(client, server.URL, params, map[string]string{}, true, false, occ, "")
	if len(out[10].Score) == 0 {
		t.Fatalf("expected scored environments")
	}
	if _, ok := out[10].Score["<"]; !ok {
		t.Fatalf("expected score for <")
	}
}
