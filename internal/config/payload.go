package config

var DefaultTags = []string{"html", "d3v", "a", "details"}

var DefaultFillings = []string{"%09", "%0a", "%0d", "/+/"}

var DefaultEFillings = []string{"%09", "%0a", "%0d", "+"}

var DefaultLFillings = []string{"", "%0dx"}

var DefaultEventHandlers = map[string][]string{
	"ontoggle":       {"details"},
	"onpointerenter": {"d3v", "details", "html", "a"},
	"onmouseover":    {"a", "html", "d3v"},
}

var DefaultFunctions = []string{
	"[8].find(confirm)",
	"confirm()",
	"(confirm)()",
}

var DefaultEnds = []string{">", "//"}
