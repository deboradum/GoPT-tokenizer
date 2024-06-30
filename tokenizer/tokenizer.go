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

func getTopBytePair(stats *map[[2]int]int) ([2]int, int) {
	var topPair [2]int
	count := 1
	for pair, pairCount := range *stats {
		if pairCount >= count {
			count = pairCount
			topPair = pair
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
		merges[topPair] = newToken

		i += 1
	}

	return tokens, merges
}

func Train(text string, vocabSize int) (map[[2]int]int, map[int][]byte) {
	tokens := EncodeConversion(text)
	newTokens, merges := BytePairEncoding(tokens, vocabSize)
	vocab := GenerateVocab(merges)
	fmt.Println("Original token length:", len(tokens), "; New token length:", len(newTokens), "; Compression ratio:", float32(len(tokens))/float32(len(newTokens)))

	return merges, vocab
}
