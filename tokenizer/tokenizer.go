package tokenizer

import (
	"sort"
)

func encodeUTF8Conversion(str string) []byte {
	return []byte(str)
}

// More efficient conversion for longer text
// https://stackoverflow.com/a/41460993
// func EncodeUTF8Reader(str string) {

// }

func bytesToInts(bytes []byte) []int {
	intSlice := make([]int, len(bytes))
	for i, b := range bytes {
		intSlice[i] = int(b)
	}

	return intSlice
}

func EncodeConversion(str string) []int {
	return bytesToInts(encodeUTF8Conversion(str))
}

func GetStats(tokens []int) map[[2]int]int {
	stats := make(map[[2]int]int)
	for i := range tokens[:len(tokens)-1] {
		pair := [2]int{tokens[i], tokens[i+1]}
		stats[pair] = stats[pair] + 1
	}

	return stats
}

// TODO: Check if all pairs are as common as each other
func getTopBytePair(stats *map[[2]int]int) ([2]int, int) {
	var topPair [2]int
	count := 1
	for k, v := range *stats {
		if v > count {
			count = v
			topPair = k
		}
	}

	return topPair, count
}

func merge(tokens []int, pair [2]int, newToken int) []int {
	newTokens := []int{}

	i := 0
	for i < len(tokens) {
		if i < len(tokens)-1 && pair[0] == tokens[i] && pair[1] == tokens[i+1] {
			newTokens = append(newTokens, newToken)
			i += 2
		} else {
			newTokens = append(newTokens, tokens[i])
			i += 1
		}
	}

	return newTokens
}

func BytePairStep(tokens []int, newToken int) []int {
	stats := GetStats(tokens)
	topPair, _ := getTopBytePair(&stats)
	newTokens := merge(tokens, topPair, newToken)

	return newTokens
}

func BytePairEncoding(tokens []int, vocabSize int) ([]int, map[[2]int]int) {
	merges := make(map[[2]int]int)

	numMerges := vocabSize - 256
	i := 0
	for i < numMerges {
		newToken := 256 + i
		stats := GetStats(tokens)
		topPair, _ := getTopBytePair(&stats)
		tokens = merge(tokens, topPair, newToken)
		// fmt.Println("Merging", topPair, "into a new token", newToken)
		merges[topPair] = newToken

		i += 1
	}

	return tokens, merges
}

func intSliceToByteSlice(ints []int) []byte {
	bytes := make([]byte, len(ints))
	for i, v := range ints {
		bytes[i] = byte(v)
	}

	return bytes
}

func Decode(tokens []int, mergeMap map[[2]int]int) string {
	// It is necessary to replace more recent merges (higher tokens) first.
	orderedKeys := make([][2]int, 0, len(mergeMap))
	for key := range mergeMap {
		orderedKeys = append(orderedKeys, key)
	}
	sort.SliceStable(orderedKeys, func(i, j int) bool {
		return mergeMap[orderedKeys[i]] > mergeMap[orderedKeys[j]]
	})

	// Go over merge map from more recent to less recent
	for _, k := range orderedKeys {
		currentMapToken := mergeMap[k]
		newTokens := []int{}
		// Swap generated tokens for their original ones.
		for _, token := range tokens {
			if token == currentMapToken {
				newTokens = append(newTokens, k[:]...)
			} else {
				newTokens = append(newTokens, token)
			}
		}

		tokens = newTokens
	}

	return string(intSliceToByteSlice(tokens))
}
