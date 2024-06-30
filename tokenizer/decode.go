package gogpt

import (
	"sort"
)

func intSliceToByteSlice(ints []int) []byte {
	bytes := make([]byte, len(ints))
	for i, v := range ints {
		bytes[i] = byte(v)
	}

	return bytes
}

// Orders the keys of the merge map so most recent byte pair merges come first.
func getOrderedMerges(mergeMap map[[2]int]int) [][2]int {
	orderedKeys := make([][2]int, 0, len(mergeMap))
	for key := range mergeMap {
		orderedKeys = append(orderedKeys, key)
	}
	sort.SliceStable(orderedKeys, func(i, j int) bool {
		return mergeMap[orderedKeys[i]] < mergeMap[orderedKeys[j]]
	})

	return orderedKeys
}

func generateVocab(mergeMap map[[2]int]int) map[int][]byte {
	vocab := make(map[int][]byte)
	i := 0
	for i < 256 {
		vocab[i] = []byte{byte(i)}
		i += 1
	}
	for _, k := range getOrderedMerges(mergeMap) {
		vocab[mergeMap[k]] = append(vocab[k[0]], vocab[k[1]]...)
	}

	return vocab
}

// TODO do not give mergemap and make vocab, but just take vocab as parameter
func Decode(tokens []int, vocab map[int][]byte) string {
	newTokens := []int{}
	for _, t := range tokens {
		// fmt.Println("token", t, "maps to", vocab[t])
		newTokens = append(newTokens, bytesToInts(vocab[t])...)
	}

	return string(intSliceToByteSlice(newTokens))
}
