package options

import (
	"flag"

	"github.com/Sakura-501/XSStrike-go/internal/config"
)

// Options stores CLI arguments in a typed structure.
type Options struct {
	URL         string
	Data        string
	Encode      string
	Fuzzer      bool
	JSON        bool
	Path        bool
	Timeout     int
	ThreadCount int
	Delay       int
	Version     bool
	HeadersRaw  string
	Limit       int
	Proxy       string
	PayloadFile string
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

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	return opts, nil
}
