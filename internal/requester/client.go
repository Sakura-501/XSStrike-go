package requester

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var defaultUserAgents = []string{
	"Mozilla/5.0 (X11; Linux i686; rv:60.0) Gecko/20100101 Firefox/60.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36 OPR/43.0.2442.991",
}

type Config struct {
	TimeoutSeconds int
	DelaySeconds   int
	Proxy          string
}

type Client struct {
	httpClient *http.Client
	cfg        Config
	rng        *rand.Rand
}

type Response struct {
	StatusCode int
	Body       string
	Headers    map[string]string
}

func New(cfg Config) *Client {
	transport := &http.Transport{}
	if cfg.Proxy != "" {
		if parsed, err := url.Parse(cfg.Proxy); err == nil {
			transport.Proxy = http.ProxyURL(parsed)
		}
	}
	if cfg.TimeoutSeconds <= 0 {
		cfg.TimeoutSeconds = 10
	}
	return &Client{
		httpClient: &http.Client{Timeout: time.Duration(cfg.TimeoutSeconds) * time.Second, Transport: transport},
		cfg:        cfg,
		rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (c *Client) DoGet(rawURL string, params map[string]string, headers map[string]string) (*Response, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	for key, value := range params {
		query.Set(key, value)
	}
	u.RawQuery = query.Encode()
	return c.doRequest(http.MethodGet, u.String(), nil, headers, false)
}

func (c *Client) DoPost(rawURL string, data map[string]string, headers map[string]string, asJSON bool) (*Response, error) {
	if asJSON {
		raw, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		if headers == nil {
			headers = map[string]string{}
		}
		headers["Content-Type"] = "application/json"
		return c.doRequest(http.MethodPost, rawURL, bytes.NewReader(raw), headers, false)
	}

	form := url.Values{}
	for key, value := range data {
		form.Set(key, value)
	}
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	return c.doRequest(http.MethodPost, rawURL, strings.NewReader(form.Encode()), headers, false)
}

func (c *Client) doRequest(method, rawURL string, body io.Reader, headers map[string]string, _ bool) (*Response, error) {
	if c.cfg.DelaySeconds > 0 {
		time.Sleep(time.Duration(c.cfg.DelaySeconds) * time.Second)
	}

	req, err := http.NewRequest(method, rawURL, body)
	if err != nil {
		return nil, err
	}

	resolvedHeaders := cloneHeaders(headers)
	ua := resolvedHeaders["User-Agent"]
	if ua == "" || ua == "$" {
		resolvedHeaders["User-Agent"] = defaultUserAgents[c.rng.Intn(len(defaultUserAgents))]
	}
	for key, value := range resolvedHeaders {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	outHeaders := map[string]string{}
	for key, values := range resp.Header {
		if len(values) > 0 {
			outHeaders[key] = values[0]
		}
	}

	return &Response{StatusCode: resp.StatusCode, Body: string(rawBody), Headers: outHeaders}, nil
}

func cloneHeaders(in map[string]string) map[string]string {
	out := map[string]string{}
	for key, value := range in {
		out[key] = value
	}
	return out
}
