package encoder

import "testing"

func TestBase64(t *testing.T) {
	got := Base64("test")
	if got != "dGVzdA==" {
		t.Fatalf("unexpected base64 value: %q", got)
	}
}

func TestApply(t *testing.T) {
	if Apply("base64", "x") != "eA==" {
		t.Fatalf("expected base64 encoded payload")
	}
	if Apply("unknown", "x") != "x" {
		t.Fatalf("expected passthrough for unknown encoding")
	}
}
