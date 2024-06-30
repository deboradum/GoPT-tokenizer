package gogpt

import "fmt"

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

func encodeConversion(str string) []int {
	return bytesToInts(encodeUTF8Conversion(str))
}

func getStats(tokens []int) map[[2]int]int {
	stats := make(map[[2]int]int)
	for i := range tokens[:len(tokens)-1] {
		pair := [2]int{tokens[i], tokens[i+1]}
		stats[pair] = stats[pair] + 1
	}

	return stats
}

func comparePairs(a, b [2]int) int {
	for i := 0; i < 2; i++ {
		if a[i] < b[i] {
			return -1
		} else if a[i] > b[i] {
			return 1
		}
	}
	return 0
}

func getTopBytePair(stats *map[[2]int]int) ([2]int, int) {
	var topPair [2]int
	topCount := 1
	first := true

	for pair, pairCount := range *stats {
		// Gets top pair and applies some logic in order to remain deterministic
		// in case two pairs are present an equal amount of time.
		if pairCount > topCount || (pairCount == topCount && first) {
			topCount = pairCount
			topPair = pair
			first = false
		} else if pairCount == topCount {
			if comparePairs(pair, topPair) < 0 {
				topPair = pair
			}
		}
	}

	return topPair, topCount
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

func bytePairEncoding(tokens []int, vocabSize int) ([]int, map[[2]int]int) {
	merges := make(map[[2]int]int)

	numMerges := vocabSize - 256
	i := 0
	for i < numMerges {
		newToken := 256 + i
		stats := getStats(tokens)
		topPair, _ := getTopBytePair(&stats)
		tokens = merge(tokens, topPair, newToken)
		merges[topPair] = newToken

		i += 1
	}

	return tokens, merges
}

func Train(text string, vocabSize int) (map[[2]int]int, map[int][]byte) {
	tokens := encodeConversion(text)
	newTokens, merges := bytePairEncoding(tokens, vocabSize)
	vocab := generateVocab(merges)
	fmt.Println("Original token length:", len(tokens), "; New token length:", len(newTokens), "; Compression ratio:", float32(len(tokens))/float32(len(newTokens)))

	return merges, vocab
}
