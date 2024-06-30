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

func Decode(tokens []int, vocab map[int][]byte) string {
	newTokens := []int{}
	for _, t := range tokens {
		// fmt.Println("token", t, "maps to", vocab[t])
		newTokens = append(newTokens, bytesToInts(vocab[t])...)
	}

	return string(intSliceToByteSlice(newTokens))
}
