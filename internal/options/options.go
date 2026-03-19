package options

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"strings"

	"github.com/Sakura-501/XSStrike-go/internal/config"
)

// Options stores CLI arguments in a typed structure.
type Options struct {
	URL          string
	Data         string
	Encode       string
	Fuzzer       bool
	JSON         bool
	Path         bool
	Crawl        bool
	SeedsFile    string
	Level        int
	SkipDOM      bool
	Blind        bool
	BlindPayload string
	Timeout      int
	ThreadCount  int
	Delay        int
	Version      bool
	HeadersRaw   string
	Limit        int
	Proxy        string
	PayloadFile  string
	OutputJSON   string
}

func Parse(fs *flag.FlagSet, args []string) (*Options, error) {
	opts := &Options{}

	fs.StringVar(&opts.URL, "u", "", "target url")
	fs.StringVar(&opts.URL, "url", "", "target url")
	fs.StringVar(&opts.Data, "data", "", "post data")
	fs.StringVar(&opts.Encode, "encode", "", "encode payloads")
	fs.BoolVar(&opts.Fuzzer, "fuzzer", false, "run fuzzer mode")
	fs.BoolVar(&opts.JSON, "json", false, "treat post data as json")
	fs.BoolVar(&opts.Path, "path", false, "inject payloads in path")
	fs.BoolVar(&opts.Crawl, "crawl", false, "enable crawl mode")
	fs.StringVar(&opts.SeedsFile, "seeds", "", "load crawling seeds from file")
	fs.IntVar(&opts.Level, "level", 2, "crawl depth level")
	fs.BoolVar(&opts.SkipDOM, "skip-dom", false, "skip DOM analysis")
	fs.BoolVar(&opts.Blind, "blind", false, "inject blind payload while crawling")
	fs.StringVar(&opts.BlindPayload, "blind-payload", "", "blind payload content")
	fs.IntVar(&opts.Timeout, "timeout", config.DefaultTimeout, "http timeout in seconds")
	fs.IntVar(&opts.ThreadCount, "threads", config.DefaultThreadCount, "number of worker threads")
	fs.IntVar(&opts.Delay, "delay", config.DefaultDelay, "delay between requests in seconds")
	fs.BoolVar(&opts.Version, "v", false, "show version")
	fs.BoolVar(&opts.Version, "version", false, "show version")
	fs.StringVar(&opts.HeadersRaw, "headers", "", "custom headers string")
	fs.IntVar(&opts.Limit, "limit", 5, "max vectors to print in fuzzer mode")
	fs.StringVar(&opts.Proxy, "proxy", "", "proxy url (example: http://127.0.0.1:8080)")
	fs.StringVar(&opts.PayloadFile, "f", "", "load payloads from a file (use 'default' for built-ins)")
	fs.StringVar(&opts.PayloadFile, "file", "", "load payloads from a file (use 'default' for built-ins)")
	fs.StringVar(&opts.OutputJSON, "output", "", "write scan report json to file")
	fs.StringVar(&opts.OutputJSON, "output-json", "", "write scan report json to file")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if err := opts.normalizeAndValidate(); err != nil {
		return nil, err
	}
	return opts, nil
}

func (o *Options) normalizeAndValidate() error {
	o.URL = strings.TrimSpace(o.URL)
	o.Data = strings.TrimSpace(o.Data)
	o.Encode = strings.ToLower(strings.TrimSpace(o.Encode))
	o.SeedsFile = strings.TrimSpace(o.SeedsFile)
	o.BlindPayload = strings.TrimSpace(o.BlindPayload)
	o.HeadersRaw = strings.TrimSpace(o.HeadersRaw)
	o.Proxy = strings.TrimSpace(o.Proxy)
	o.PayloadFile = strings.TrimSpace(o.PayloadFile)
	o.OutputJSON = strings.TrimSpace(o.OutputJSON)

	if o.Timeout <= 0 {
		return errors.New("timeout must be greater than 0")
	}
	if o.ThreadCount <= 0 {
		return errors.New("threads must be greater than 0")
	}
	if o.Delay < 0 {
		return errors.New("delay cannot be negative")
	}
	if o.Level < 0 {
		return errors.New("level cannot be negative")
	}
	if o.Limit < 0 {
		o.Limit = 0
	}
	if o.Encode != "" && o.Encode != "base64" {
		return fmt.Errorf("unsupported encode mode %q", o.Encode)
	}
	if o.Proxy != "" {
		parsed, err := url.Parse(o.Proxy)
		if err != nil || parsed.Scheme == "" || parsed.Host == "" {
			return fmt.Errorf("invalid proxy url %q", o.Proxy)
		}
	}
	if o.Crawl && o.URL == "" && o.SeedsFile == "" {
		return errors.New("crawl mode requires --url or --seeds")
	}
	if o.Blind {
		if !o.Crawl && o.SeedsFile == "" {
			return errors.New("blind mode requires crawl mode or --seeds")
		}
		if o.BlindPayload == "" {
			return errors.New("blind mode requires --blind-payload")
		}
	}
	return nil
}
