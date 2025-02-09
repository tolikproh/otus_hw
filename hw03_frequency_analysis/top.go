package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

// Top10 возвращает 10 самых часто встречающихся слов в тексте.
// Если слова имеют одинаковую частоту, они сортируются лексикографически.
func Top10(text string) []string {
	if len(strings.TrimSpace(text)) == 0 {
		return []string{}
	}

	// Разбирается текст на слова.
	words := strings.Fields(text)

	// Создаётся map для подсчета частоты слов.
	frequency := make(map[string]int)
	for _, word := range words {
		word = strings.ToLower(word)
		word = strings.Trim(word, "!\"#$%&'()*+,./:;<=>?@[]\\^{|}~")
		if word == "-" || word == "" {
			continue
		}
		frequency[word]++
	}

	// Создается слайс для хранения уникальных слов
	uniqueWords := make([]string, 0, len(frequency))
	for word := range frequency {
		uniqueWords = append(uniqueWords, word)
	}

	// Сортировка слов по частоте и лексикографически
	sort.Slice(uniqueWords, func(i, j int) bool {
		// Если частота одинаковая, сортируется лексикографически
		if frequency[uniqueWords[i]] == frequency[uniqueWords[j]] {
			return uniqueWords[i] < uniqueWords[j]
		}
		// Иначе сортируется по частоте (по убыванию)
		return frequency[uniqueWords[i]] > frequency[uniqueWords[j]]
	})

	// Возвращается топ-10 слов или меньше, если слов меньше 10
	if len(uniqueWords) > 10 {
		return uniqueWords[:10]
	}

	return uniqueWords
}
