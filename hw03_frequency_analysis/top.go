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
	w := []Words{}
	wordMap := make(map[string]int)

	for _, word := range strings.Fields(s) {
		wordMap[word]++
	}

	for key, val := range wordMap {
		var temp Words
		temp.Word = key
		temp.Count = val
		w = append(w, temp)
	}

	sort.Slice(w, func(i, j int) bool {
		return w[i].Count > w[j].Count || (w[i].Count == w[j].Count && (strings.Compare(w[i].Word, w[j].Word) < 0))
	})

	res := []string{}

	for _, val := range w {
		res = append(res, val.Word)
	}

	if len(res) < 11 {
		return res
	}

	return res[:10]
}
