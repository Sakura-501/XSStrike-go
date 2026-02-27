package reflection

import (
	"regexp"
	"strings"

	"github.com/Sakura-501/XSStrike-go/internal/config"
	"github.com/Sakura-501/XSStrike-go/internal/encoder"
	"github.com/Sakura-501/XSStrike-go/internal/utils"
)

var (
	reHTMLComment      = regexp.MustCompile(`(?s)<!--[\s\S]*?-->`)
	reScriptTag        = regexp.MustCompile(`(?is)<script[^>]*>(.*?)</script>`)
	reAttrContext      = regexp.MustCompile(`(?is)<[^>]*?(v3dm0s)[^>]*?>`)
	reCommentContext   = regexp.MustCompile(`(?is)<!--[\s\S]*?(v3dm0s)[\s\S]*?-->`)
	reScriptLineMarker = regexp.MustCompile(`(?s)(v3dm0s.*?)$`)
)

var nonExecTags = []string{"style", "template", "textarea", "title", "noembed", "noscript"}

func Parse(responseBody string, encodeMode string) Occurrences {
	body := responseBody
	if encodeMode != "" {
		encoded := encoder.Apply(encodeMode, config.XSSChecker)
		if encoded != "" && encoded != config.XSSChecker {
			body = strings.ReplaceAll(body, encoded, config.XSSChecker)
		}
	}

	reflections := strings.Count(body, config.XSSChecker)
	if reflections == 0 {
		return Occurrences{}
	}

	positionContext := map[int]string{}
	envDetails := map[int]Details{}

	cleanResponse := reHTMLComment.ReplaceAllString(body, "")

	scriptMatches := reScriptTag.FindAllStringSubmatchIndex(cleanResponse, -1)
	for _, scriptMatch := range scriptMatches {
		if len(scriptMatch) < 4 {
			continue
		}
		scriptBody := cleanResponse[scriptMatch[2]:scriptMatch[3]]
		offset := 0
		for {
			local := strings.Index(scriptBody[offset:], config.XSSChecker)
			if local == -1 {
				break
			}
			localPos := offset + local
			thisPosition := scriptMatch[2] + localPos
			positionContext[thisPosition] = "script"
			detail := Details{Quote: ""}
			fragment := scriptBody[localPos:]
			for i := 0; i < len(fragment); i++ {
				current := fragment[i]
				if (current == '/' || current == '\'' || current == '`' || current == '"') && !utils.Escaped(i, fragment) {
					detail.Quote = string(current)
				} else if (current == ')' || current == ']' || current == '}') && !utils.Escaped(i, fragment) {
					break
				}
			}
			envDetails[thisPosition] = detail
			offset = localPos + len(config.XSSChecker)
		}
	}

	if len(positionContext) < reflections {
		matches := reAttrContext.FindAllStringSubmatchIndex(cleanResponse, -1)
		for _, match := range matches {
			if len(match) < 4 {
				continue
			}
			thisPosition := match[2]
			if _, ok := positionContext[thisPosition]; ok {
				continue
			}
			tagText := cleanResponse[match[0]:match[1]]
			parts := strings.Fields(tagText)
			tag := ""
			if len(parts) > 0 {
				tag = strings.TrimPrefix(parts[0], "<")
			}

			detail := Details{Tag: tag}
			for _, part := range parts {
				if !strings.Contains(part, config.XSSChecker) {
					continue
				}
				if strings.Contains(part, "=") {
					detail.Type = "value"
					chunks := strings.SplitN(part, "=", 2)
					detail.Name = chunks[0]
					rawValue := chunks[1]
					if strings.HasPrefix(rawValue, "\"") || strings.HasPrefix(rawValue, "'") || strings.HasPrefix(rawValue, "`") {
						detail.Quote = rawValue[:1]
					}
					rawValue = strings.TrimSuffix(rawValue, ">")
					rawValue = strings.Trim(rawValue, "\"'`")
					detail.Value = rawValue
					if detail.Name == config.XSSChecker {
						detail.Type = "name"
					}
				} else {
					detail.Type = "flag"
				}
			}

			positionContext[thisPosition] = "attribute"
			envDetails[thisPosition] = detail
		}
	}

	if len(positionContext) < reflections {
		offset := 0
		for {
			idx := strings.Index(cleanResponse[offset:], config.XSSChecker)
			if idx == -1 {
				break
			}
			pos := offset + idx
			if _, ok := positionContext[pos]; !ok {
				positionContext[pos] = "html"
				envDetails[pos] = Details{}
			}
			offset = pos + len(config.XSSChecker)
		}
	}

	if len(positionContext) < reflections {
		matches := reCommentContext.FindAllStringSubmatchIndex(body, -1)
		for _, match := range matches {
			if len(match) < 4 {
				continue
			}
			pos := match[2]
			if _, ok := positionContext[pos]; ok {
				continue
			}
			positionContext[pos] = "comment"
			envDetails[pos] = Details{}
		}
	}

	database := Occurrences{}
	for _, pos := range sortedKeys(positionContext) {
		database[pos] = &Occurrence{Position: pos, Context: positionContext[pos], Details: envDetails[pos], Score: map[string]int{}}
	}

	for _, tag := range nonExecTags {
		re := regexp.MustCompile(`(?is)<` + tag + `>[\s\S]*?` + regexp.QuoteMeta(config.XSSChecker) + `[\s\S]*?</` + tag + `>`)
		nonExec := re.FindAllStringIndex(body, -1)
		for _, match := range nonExec {
			if len(match) < 2 {
				continue
			}
			start, end := match[0], match[1]
			for _, occ := range database {
				if start < occ.Position && occ.Position < end {
					occ.Details.BadTag = tag
				}
			}
		}
	}
	for _, occ := range database {
		if occ.Details.BadTag == "" {
			occ.Details.BadTag = ""
		}
	}

	return database
}

func sortedKeys(in map[int]string) []int {
	keys := make([]int, 0, len(in))
	for key := range in {
		keys = append(keys, key)
	}
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[j] < keys[i] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}
	return keys
}
