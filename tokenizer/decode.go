package gogpt

import "sort"

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
		return mergeMap[orderedKeys[i]] > mergeMap[orderedKeys[j]]
	})

	return orderedKeys
}

func Decode(tokens []int, mergeMap map[[2]int]int) string {
	// It is necessary to replace more recent merges (higher tokens) first.
	for _, k := range getOrderedMerges(mergeMap) {
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
