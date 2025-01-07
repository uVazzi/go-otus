package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type WordCounter struct {
	word  string
	count int
}

func Top10(inputText string) []string {
	if inputText == "" {
		return nil
	}

	countWordsByKey := make(map[string]int)
	for _, word := range strings.Fields(inputText) {
		countWordsByKey[word]++
	}

	structWords := make([]WordCounter, 0)
	for word, count := range countWordsByKey {
		structWords = append(structWords, WordCounter{word, count})
	}

	sort.Slice(structWords, func(i, j int) bool {
		if structWords[i].count == structWords[j].count {
			return structWords[i].word < structWords[j].word
		}
		return structWords[i].count > structWords[j].count
	})

	result := make([]string, 0, 10)
	for _, item := range structWords {
		result = append(result, item.word)
		if len(result) == 10 {
			break
		}
	}

	return result
}
