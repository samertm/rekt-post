package main

// used for data gathering
// now here for dust gathering

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