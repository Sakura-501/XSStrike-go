package options

import (
	"flag"
	"strings"
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
	if opts.OutputJSON != "" {
		t.Fatalf("output json should default to empty")
	}
	if opts.Level != 2 {
		t.Fatalf("crawl level default mismatch: %d", opts.Level)
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
		"--output", "out.json",
		"--crawl",
		"--seeds", "seeds.txt",
		"--level", "3",
		"--blind",
		"--blind-payload", "xss.example/abc",
	})
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if opts.URL != "https://target.local" {
		t.Fatalf("url mismatch: got %q", opts.URL)
	}
	if !opts.JSON || !opts.Path || !opts.Fuzzer || !opts.Crawl || !opts.Blind {
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
	if opts.OutputJSON != "out.json" {
		t.Fatalf("output json mismatch: got %q", opts.OutputJSON)
	}
	if opts.SeedsFile != "seeds.txt" || opts.Level != 3 {
		t.Fatalf("crawl args mismatch: %+v", opts)
	}
	if opts.BlindPayload != "xss.example/abc" {
		t.Fatalf("blind payload mismatch: %q", opts.BlindPayload)
	}
}

func TestParseNormalizesEncodeAndLimit(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	opts, err := Parse(fs, []string{
		"--url", "https://target.local",
		"--encode", " BASE64 ",
		"--limit", "-2",
	})
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if opts.Encode != "base64" {
		t.Fatalf("encode should be normalized, got %q", opts.Encode)
	}
	if opts.Limit != 0 {
		t.Fatalf("negative limit should normalize to 0, got %d", opts.Limit)
	}
}

func TestParseRejectsInvalidValues(t *testing.T) {
	testCases := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "unsupported encode",
			args:    []string{"--url", "https://target.local", "--encode", "hex"},
			wantErr: `unsupported encode mode "hex"`,
		},
		{
			name:    "invalid proxy",
			args:    []string{"--url", "https://target.local", "--proxy", "127.0.0.1:8080"},
			wantErr: `invalid proxy url "127.0.0.1:8080"`,
		},
		{
			name:    "blind without payload",
			args:    []string{"--crawl", "--url", "https://target.local", "--blind"},
			wantErr: "blind mode requires --blind-payload",
		},
		{
			name:    "crawl without seed",
			args:    []string{"--crawl"},
			wantErr: "crawl mode requires --url or --seeds",
		},
		{
			name:    "negative timeout",
			args:    []string{"--url", "https://target.local", "--timeout", "0"},
			wantErr: "timeout must be greater than 0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fs := flag.NewFlagSet("test", flag.ContinueOnError)
			_, err := Parse(fs, tc.args)
			if err == nil {
				t.Fatalf("expected error containing %q", tc.wantErr)
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
