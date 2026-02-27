package reflection

import (
	"sort"

	"github.com/Sakura-501/XSStrike-go/internal/requester"
)

func FilterCheck(client *requester.Client, url string, params map[string]string, headers map[string]string, isGET bool, jsonData bool, occurrences Occurrences, encodeMode string) Occurrences {
	positions := occurrences.Positions()
	environments := map[string]struct{}{"<": {}, ">": {}}

	for _, pos := range positions {
		occ := occurrences[pos]
		if occ.Score == nil {
			occ.Score = map[string]int{}
		}
		switch occ.Context {
		case "comment":
			environments["-->"] = struct{}{}
		case "script":
			if occ.Details.Quote != "" {
				environments[occ.Details.Quote] = struct{}{}
			}
			environments["</scRipT/>"] = struct{}{}
		case "attribute":
			if occ.Details.Type == "value" && occ.Details.Name == "srcdoc" {
				environments["&lt;"] = struct{}{}
				environments["&gt;"] = struct{}{}
			}
			if occ.Details.Quote != "" {
				environments[occ.Details.Quote] = struct{}{}
			}
		}
	}

	envKeys := make([]string, 0, len(environments))
	for env := range environments {
		envKeys = append(envKeys, env)
	}
	sort.Strings(envKeys)

	for _, env := range envKeys {
		if env == "" {
			continue
		}
		efficiencies := Check(client, url, params, headers, isGET, jsonData, env, positions, encodeMode)
		if len(efficiencies) < len(positions) {
			pad := make([]int, len(positions)-len(efficiencies))
			efficiencies = append(efficiencies, pad...)
		}
		for i, pos := range positions {
			occurrences[pos].Score[env] = efficiencies[i]
		}
	}
	return occurrences
}
