package utils

import (
	"math"
	"regexp"
	"strings"
)

var reCounter = regexp.MustCompile(`\s|\w`)

// Counter returns number of non-space and non-word characters.
func Counter(input string) int {
	stripped := reCounter.ReplaceAllString(input, "")
	return len(stripped)
}

// Closest returns the map entry whose value is closest to number.
func Closest(number int, numbers map[int]int) map[int]int {
	result := map[int]int{}
	bestDiff := math.MaxInt
	for idx, value := range numbers {
		diff := int(math.Abs(float64(number - value)))
		if diff < bestDiff {
			bestDiff = diff
			result = map[int]int{idx: value}
		}
	}
	return result
}

// FillHoles mirrors XSStrike fillHoles behavior.
func FillHoles(original []int, newer []int) []int {
	filler := 0
	filled := []int{}
	length := len(original)
	if len(newer) < length {
		length = len(newer)
	}
	for i := 0; i < length; i++ {
		x := original[i]
		y := newer[i]
		if x == (y + filler) {
			filled = append(filled, y)
		} else {
			filled = append(filled, 0, y)
			filler += (x - y)
		}
	}
	return filled
}

// Stripper removes first occurrence of substring from left/right side traversal.
func Stripper(input, substring, direction string) string {
	if len(substring) != 1 {
		return input
	}
	target := substring
	working := input
	if direction == "right" {
		working = reverse(working)
	}
	removed := false
	builder := strings.Builder{}
	for _, ch := range working {
		current := string(ch)
		if current == target && !removed {
			removed = true
			continue
		}
		builder.WriteRune(ch)
	}
	result := builder.String()
	if direction == "right" {
		return reverse(result)
	}
	return result
}

// DeJSON converts escaped backslashes to plain backslashes.
func DeJSON(data string) string {
	return strings.ReplaceAll(data, `\\`, `\`)
}

type ContextRange struct {
	Start int
	End   int
	Name  string
}

// IsBadContext returns context name when position is in any non-executable range.
func IsBadContext(position int, contexts []ContextRange) string {
	for _, current := range contexts {
		if current.Start < position && position < current.End {
			return current.Name
		}
	}
	return ""
}

// Equalize appends one empty string when len(values) is below number.
func Equalize(values []string, number int) []string {
	if len(values) < number {
		values = append(values, "")
	}
	return values
}

// Escaped checks if character at position is escaped by odd count of preceding backslashes.
func Escaped(position int, input string) bool {
	if position <= 0 || position > len(input) {
		return false
	}
	count := 0
	for i := position - 1; i >= 0; i-- {
		if input[i] == '\\' {
			count++
			continue
		}
		break
	}
	if count == 0 {
		return false
	}
	return count%2 == 1
}

func reverse(input string) string {
	runes := []rune(input)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
