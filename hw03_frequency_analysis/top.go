package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type freqWord struct {
	Word string
	Freq int
}

var rg = regexp.MustCompile(`(!|\.|,|\s-)*\s+`)

func Top10(text string) []string {
	slicedText := rg.Split(text, -1)
	if len(slicedText) == 1 {
		slicedText = []string{}
		return slicedText
	}
	freqStr := []freqWord{}
	sortedFreq := []string{}
	keys := make(map[string]int)
	for _, value := range slicedText {
		word := strings.ToLower(value)
		if _, value := keys[word]; !value {
			keys[word] = 1
		} else {
			keys[word]++
		}
	}

	for k, v := range keys {
		freqStr = append(freqStr, freqWord{Word: k, Freq: v})
	}

	sort.Slice(freqStr, func(i, j int) bool {
		if freqStr[i].Freq == freqStr[j].Freq {
			return freqStr[i].Word < freqStr[j].Word
		}

		return freqStr[i].Freq > freqStr[j].Freq
	})
	for _, item := range freqStr {
		sortedFreq = append(sortedFreq, item.Word)
	}
	var top10Freq []string

	if len(sortedFreq) > 10 {
		top10Freq = sortedFreq[:10]
	} else {
		top10Freq = sortedFreq
	}
	return top10Freq
}
