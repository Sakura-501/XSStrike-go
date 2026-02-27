package reflection

import (
	"strings"

	"github.com/Sakura-501/XSStrike-go/internal/config"
	"github.com/Sakura-501/XSStrike-go/internal/encoder"
	"github.com/Sakura-501/XSStrike-go/internal/requester"
	"github.com/Sakura-501/XSStrike-go/internal/utils"
)

func Check(client *requester.Client, url string, params map[string]string, headers map[string]string, isGET bool, jsonData bool, payload string, positions []int, encodeMode string) []int {
	if client == nil {
		return make([]int, len(positions))
	}

	checkString := "st4r7s" + payload + "3nd"
	sentCheck := checkString
	if encodeMode != "" {
		sentCheck = encoder.Apply(encodeMode, checkString)
	}

	requestParams := replaceValueMap(params, config.XSSChecker, sentCheck)
	var (
		resp *requester.Response
		err  error
	)
	if isGET {
		resp, err = client.DoGet(url, requestParams, headers)
	} else {
		resp, err = client.DoPost(url, requestParams, headers, jsonData)
	}
	if err != nil || resp == nil {
		return make([]int, len(positions))
	}

	response := strings.ToLower(resp.Body)
	reflectedPositions := findAll(response, "st4r7s")
	if len(reflectedPositions) == 0 && encodeMode != "" {
		encodedToken := strings.ToLower(encoder.Apply(encodeMode, "st4r7s"))
		reflectedPositions = findAll(response, encodedToken)
	}

	filled := utils.FillHoles(positions, reflectedPositions)
	efficiencies := make([]int, 0, len(filled))
	lowerCheck := strings.ToLower(checkString)
	if encodeMode != "" {
		lowerCheck = strings.ToLower(encoder.Apply(encodeMode, lowerCheck))
	}

	num := 0
	for _, position := range filled {
		best := 0
		if num < len(reflectedPositions) {
			fragment := sliceSafe(response, reflectedPositions[num], len(strings.ToLower(sentCheck)))
			eff := partialRatio(fragment, strings.ToLower(sentCheck))
			if eff > best {
				best = eff
			}
		}
		if position > 0 {
			fragment := sliceSafe(response, position, len(lowerCheck))
			eff := partialRatio(fragment, lowerCheck)
			core := strings.ReplaceAll(strings.ReplaceAll(lowerCheck, "st4r7s", ""), "3nd", "")
			if len(fragment) >= 2 && fragment[:len(fragment)-2] == "\\"+core {
				eff = 90
			}
			if eff > best {
				best = eff
			}
		}
		efficiencies = append(efficiencies, best)
		num++
	}
	return efficiencies
}

func replaceValueMap(in map[string]string, oldValue string, newValue string) map[string]string {
	out := map[string]string{}
	for k, v := range in {
		if v == oldValue {
			out[k] = newValue
		} else {
			out[k] = v
		}
	}
	return out
}

func findAll(text string, needle string) []int {
	positions := []int{}
	if needle == "" {
		return positions
	}
	offset := 0
	for {
		idx := strings.Index(text[offset:], needle)
		if idx == -1 {
			break
		}
		positions = append(positions, offset+idx)
		offset += idx + len(needle)
	}
	return positions
}

func sliceSafe(text string, start int, size int) string {
	if start < 0 || start >= len(text) || size <= 0 {
		return ""
	}
	end := start + size
	if end > len(text) {
		end = len(text)
	}
	return text[start:end]
}

func partialRatio(a string, b string) int {
	if a == "" || b == "" {
		return 0
	}
	if a == b {
		return 100
	}
	shorter, longer := a, b
	if len(shorter) > len(longer) {
		shorter, longer = longer, shorter
	}
	if len(shorter) == 0 {
		return 0
	}

	best := 0
	window := len(shorter)
	for i := 0; i <= len(longer)-window; i++ {
		score := ratio(shorter, longer[i:i+window])
		if score > best {
			best = score
		}
		if best == 100 {
			return 100
		}
	}
	return best
}

func ratio(a string, b string) int {
	if len(a) == 0 {
		return 0
	}
	dist := levenshtein(a, b)
	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}
	score := (maxLen - dist) * 100 / maxLen
	if score < 0 {
		return 0
	}
	if score > 100 {
		return 100
	}
	return score
}

func levenshtein(a string, b string) int {
	ra := []rune(a)
	rb := []rune(b)
	if len(ra) == 0 {
		return len(rb)
	}
	if len(rb) == 0 {
		return len(ra)
	}

	dp := make([][]int, len(ra)+1)
	for i := range dp {
		dp[i] = make([]int, len(rb)+1)
	}
	for i := 0; i <= len(ra); i++ {
		dp[i][0] = i
	}
	for j := 0; j <= len(rb); j++ {
		dp[0][j] = j
	}

	for i := 1; i <= len(ra); i++ {
		for j := 1; j <= len(rb); j++ {
			cost := 0
			if ra[i-1] != rb[j-1] {
				cost = 1
			}
			del := dp[i-1][j] + 1
			ins := dp[i][j-1] + 1
			sub := dp[i-1][j-1] + cost
			dp[i][j] = min(del, min(ins, sub))
		}
	}
	return dp[len(ra)][len(rb)]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
