package reflection

import "sort"

type Details struct {
	Tag    string `json:"tag,omitempty"`
	Type   string `json:"type,omitempty"`
	Quote  string `json:"quote,omitempty"`
	Value  string `json:"value,omitempty"`
	Name   string `json:"name,omitempty"`
	BadTag string `json:"bad_tag,omitempty"`
}

type Occurrence struct {
	Position int            `json:"position"`
	Context  string         `json:"context"`
	Details  Details        `json:"details"`
	Score    map[string]int `json:"score,omitempty"`
}

type Occurrences map[int]*Occurrence

func (o Occurrences) Positions() []int {
	keys := make([]int, 0, len(o))
	for k := range o {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

func (o Occurrences) Count() int {
	return len(o)
}
