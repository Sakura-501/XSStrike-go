package options

import (
	"flag"
	"testing"

	"github.com/Sakura-501/XSStrike-go/internal/config"
)

func TestParseDefaults(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	opts, err := Parse(fs, []string{"-u", "https://example.com"})
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if opts.Timeout != config.DefaultTimeout {
		t.Fatalf("timeout mismatch: got %d want %d", opts.Timeout, config.DefaultTimeout)
	}
	if opts.ThreadCount != config.DefaultThreadCount {
		t.Fatalf("thread mismatch: got %d want %d", opts.ThreadCount, config.DefaultThreadCount)
	}
	if opts.Delay != config.DefaultDelay {
		t.Fatalf("delay mismatch: got %d want %d", opts.Delay, config.DefaultDelay)
	}
	if opts.Proxy != "" {
		t.Fatalf("proxy should default to empty")
	}
}

func TestParseFlags(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	opts, err := Parse(fs, []string{
		"--url", "https://target.local",
		"--data", "a=1&b=2",
		"--json",
		"--path",
		"--fuzzer",
		"--timeout", "20",
		"--threads", "3",
		"--delay", "1",
		"--encode", "base64",
		"--proxy", "http://127.0.0.1:8080",
	})
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if opts.URL != "https://target.local" {
		t.Fatalf("url mismatch: got %q", opts.URL)
	}
	if !opts.JSON || !opts.Path || !opts.Fuzzer {
		t.Fatalf("boolean flags were not parsed correctly: %+v", opts)
	}
	if opts.Timeout != 20 || opts.ThreadCount != 3 || opts.Delay != 1 {
		t.Fatalf("numeric flags mismatch: %+v", opts)
	}
	if opts.Encode != "base64" {
		t.Fatalf("encode mismatch: got %q", opts.Encode)
	}
	if opts.Proxy != "http://127.0.0.1:8080" {
		t.Fatalf("proxy mismatch: got %q", opts.Proxy)
	}
}
