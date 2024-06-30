package gogpt

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"
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

func saveVocab(name string, vocab map[int][]byte) error {
	filename := name + ".vocab"
	if _, err := os.Stat(filename); err == nil {
		return errors.New("file already exists")
	}

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(vocab)
	if err != nil {
		fmt.Println("Error encoding merges:", err)
		return err
	}

	return nil
}

func saveMerges(name string, merges map[[2]int]int) error {
	filename := name + ".bpe"
	if _, err := os.Stat(filename); err == nil {
		return errors.New("file already exists")
	}

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(merges)
	if err != nil {
		fmt.Println("Error encoding merges:", err)
		return err
	}

	return nil
}

func readMerges(filename string) (map[[2]int]int, error) {
	var merges map[[2]int]int
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return merges, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&merges)
	if err != nil {
		fmt.Println("Error decoding map:", err)
		return merges, err
	}

	return merges, nil
}

func readVocab(filename string) (map[int][]byte, error) {
	var vocab map[int][]byte
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return vocab, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&vocab)
	if err != nil {
		fmt.Println("Error decoding map:", err)
		return vocab, err
	}

	return vocab, nil
}

func Train(text string, vocabSize int, dataName string) (map[[2]int]int, map[int][]byte) {
	tokens := encodeConversion(text)
	newTokens, merges := bytePairEncoding(tokens, vocabSize)
	vocab := generateVocab(merges)
	fmt.Println("Original token length:", len(tokens), "; New token length:", len(newTokens), "; Compression ratio:", float32(len(tokens))/float32(len(newTokens)))

	saveMerges(dataName, merges)
	saveVocab(dataName, vocab)

	return merges, vocab
}

func LoadTokenizer(mergesPath string, vocabPath string) (map[[2]int]int, map[int][]byte, error) {
	merges, err := readMerges(mergesPath)
	if err != nil {
		fmt.Println("Error loading merges:", err)
		return make(map[[2]int]int), make(map[int][]byte), err
	}
	vocab, err := readVocab(vocabPath)
	if err != nil {
		fmt.Println("Error loading vocab:", err)
		return make(map[[2]int]int), make(map[int][]byte), err
	}

	return merges, vocab, nil
}
