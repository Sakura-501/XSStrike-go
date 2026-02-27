package waf

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sakura-501/XSStrike-go/internal/requester"
)

func TestNewDefault(t *testing.T) {
	detector, err := NewDefault()
	if err != nil {
		t.Fatalf("expected default signatures to load: %v", err)
	}
	if len(detector.Signatures) == 0 {
		t.Fatalf("expected non-empty signatures")
	}
}

func TestDetectCloudflareLike(t *testing.T) {
	signer := &Detector{Signatures: map[string]Signature{
		"CloudFlare": {
			Code:    "403",
			Page:    "Attention Required! \\| Cloudflare",
			Headers: "cf-ray",
		},
	}}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("cf-ray", "abcd")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("Attention Required! | Cloudflare"))
	}))
	defer server.Close()

	client := requester.New(requester.Config{TimeoutSeconds: 5})
	result := signer.Detect(client, server.URL, map[string]string{"q": "1"}, map[string]string{}, true, false)
	if !result.Detected {
		t.Fatalf("expected WAF detection result")
	}
	if result.Name != "CloudFlare" {
		t.Fatalf("unexpected detected name: %s", result.Name)
	}
}

func TestDetectNoMatch(t *testing.T) {
	signer := &Detector{Signatures: map[string]Signature{
		"AnyWAF": {Code: "403", Page: "blocked", Headers: "x-waf"},
	}}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("all good"))
	}))
	defer server.Close()

	client := requester.New(requester.Config{TimeoutSeconds: 5})
	result := signer.Detect(client, server.URL, map[string]string{"q": "1"}, map[string]string{}, true, false)
	if result.Detected {
		t.Fatalf("did not expect detection")
	}
}
