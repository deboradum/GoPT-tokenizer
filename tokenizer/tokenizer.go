package tokenizer

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
func GetMostCommonBytePair(stats *map[[2]int]int) ([2]int, int) {
	var mostCommon [2]int
	count := 1
	for k, v := range *stats {
		if v > count {
			count = v
			mostCommon = k
		}
	}

	return mostCommon, count
}
