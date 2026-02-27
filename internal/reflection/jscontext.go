package reflection

import (
	"regexp"
	"strings"

	"github.com/Sakura-501/XSStrike-go/internal/config"
	"github.com/Sakura-501/XSStrike-go/internal/utils"
)

var reStripJSBlocks = regexp.MustCompile(`(?s)\{.*?\}|\(.*?\)|".*?"|'.*?'`)

func JSContexter(script string) string {
	broken := strings.Split(script, config.XSSChecker)
	if len(broken) == 0 {
		return ""
	}
	pre := reStripJSBlocks.ReplaceAllString(broken[0], "")
	breaker := ""
	for i := 0; i < len(pre); i++ {
		char := pre[i]
		switch char {
		case '{':
			breaker += "}"
		case '(':
			breaker += ";)"
		case '[':
			breaker += "]"
		case '/':
			if i+1 < len(pre) && pre[i+1] == '*' {
				breaker += "/*"
			}
		case '}':
			breaker = utils.Stripper(breaker, "}", "right")
		case ')':
			breaker = utils.Stripper(breaker, ")", "right")
		case ']':
			breaker = utils.Stripper(breaker, "]", "right")
		}
	}
	return reverse(breaker)
}

func reverse(input string) string {
	r := []rune(input)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
