package main

import (
	"fmt"
)

// used for data gathering
// now here for dust gathering

// the only function you should use
func orderedFreqs(posts []post) string {
	cats := make(map[string]int)
	for _, p := range posts {
		cats = concatFreqs(cats, p.freqs)
	}
	fs := sortFreqs(cats)
	var s string
	for _, fp := range fs {
		s += fmt.Sprint("\"", fp.word, "\", ", fp.freq, "\n")
	}
	return s
}

func concatFreqs(freq1, freq2 map[string]int) map[string]int {
	catted := make(map[string]int)
	for k, v := range freq1 {
		catted[k] = v
	}
	for k, v := range freq2 {
		catted[k] += v
	}
	return catted
}

type freqPair struct {
	word string
	freq int
}

func sortFreqs(freqs map[string]int) []freqPair {
	var sortedAdd func([]freqPair, freqPair) []freqPair
	sortedAdd = func(fps []freqPair, fp freqPair) []freqPair {
		for i := range fps {
			if fp.freq > fps[i].freq {
				// add where i is
				return append(fps[:i], append([]freqPair{fp}, fps[i:]...)...)
			}
		}
		// add to end
		return append(fps, fp)
	}
	storeFreqs := make([]freqPair, 0, len(freqs))
	for k, v := range freqs {
		storeFreqs = sortedAdd(storeFreqs, freqPair{word: k, freq: v})
	}
	return storeFreqs
}

