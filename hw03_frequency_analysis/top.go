package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type Words struct {
	Word  string
	Count int
}

func Top10(s string) []string {
	if len(s) == 0 {
		return []string{}
	}

	wordsMap := make(map[string]int)

	for _, word := range strings.Fields(s) {
		wordsMap[word]++
	}

	w := make([]Words, 0, len(wordsMap))

	for key, val := range wordsMap {
		w = append(w, Words{key, val})
	}

	sort.Slice(w, func(i, j int) bool {
		return w[i].Count > w[j].Count || (w[i].Count == w[j].Count && (strings.Compare(w[i].Word, w[j].Word) < 0))
	})

	res := make([]string, 0, 10)

	for i := 0; i < len(w) && i < 11; i++ {
		res = append(res, w[i].Word)
	}

	if len(res) < 11 {
		return res
	}

	return res[:10]
}
