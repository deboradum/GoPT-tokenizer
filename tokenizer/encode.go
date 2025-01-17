package gogpt

import (
	"errors"
	"math"
	"regexp"
	"slices"
)

// Receives a pointer to a stats map and a pointer to a merges map. Returns the
// pair in the stats map with the lowest corresponding merge value in the merge
// map. (We want to tokenize ealier tokens first)
func getPairToTokenize(stats *map[[2]int]int, mergeMap map[[2]int]int) ([2]int, int, error) {
	lowestMergeValue := int(math.Inf(1))
	pair := [2]int{-1, -1}
	for k := range *stats {
		if mergeMap[k] > 0 && mergeMap[k] < lowestMergeValue {
			lowestMergeValue = mergeMap[k]
			pair = k
		}
	}
	if lowestMergeValue == int(math.Inf(1)) {
		return pair, lowestMergeValue, errors.New("no mergable pair present")
	}

	return pair, lowestMergeValue, nil
}

func preProcess(text string) []byte {
	pattern := `'(?i:[sdmt]|ll|ve|re)|[^\r\n\p{L}\p{N}]?\p{L}+|\p{N}{1,3}| ?[^\s\p{L}\p{N}]+[\r\n]*|\s*[\r\n]|\s+`
	re := regexp.MustCompile(pattern)
	splitted := re.FindAll([]byte(text), -1)

	return slices.Concat(splitted...)
}

func Encode(text string, mergeMap map[[2]int]int) []int {
	// Pre processing step, runs some regex to split certain character groups.
	tokens := bytesToInts(preProcess(text))

	for {
		stats := getStats(tokens)
		pair, newToken, err := getPairToTokenize(&stats, mergeMap)
		// When nothing can be merged
		if err != nil {
			break
		}
		tokens = merge(tokens, pair, newToken)
	}

	return tokens
}
