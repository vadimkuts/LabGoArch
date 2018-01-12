package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/pkg/profile"
	"io"
	"log"
	"os"
	"sort"
)

// A little bit optimized Levenshtein algorithm, so it uses O(min(m,n))
// space instead of O(mn), where m and n - lengths of compared strings.
// The key observation is that we only need to access the contents
// of the previous column when filling the matrix column-by-column.
// Hence, we can re-use a single column over and over, overwriting its contents as we proceed.
func levenshteinDistance(s1, s2 string) int {
	var prevDiagonalValue, buffer int
	s1len := len(s1)
	s2len := len(s2)

	// Initialize column
	var curColumn = make([]int, s1len+1)
	for i := 0; i < s1len; i++ {
		curColumn[i] = i
	}

	// Fill matrix column by column
	for i := 1; i <= s2len; i++ {
		curColumn[0] = i
		prevDiagonalValue = i - 1
		for j := 1; j <= s1len; j++ {
			// Set operation cost (all operations except match(M) has value 1)
			operationCost := 1
			if s1[j-1] == s2[i-1] {
				operationCost = 0
			}
			buffer = curColumn[j]
			curColumn[j] = minOfThree(curColumn[j]+1, curColumn[j-1]+1, prevDiagonalValue+operationCost)
			prevDiagonalValue = buffer
		}
	}
	return curColumn[s1len]
}

func run(startWord string, file io.Reader) Words {
	var currentWord Word
	var words Words

	// Read file word by word
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentWord.Text = scanner.Text()
		currentWord.Distance = levenshteinDistance(startWord, currentWord.Text)
		words = append(words, currentWord)
	}
	// Handle scanner errors
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	sort.Sort(words)
	return words
}

func main() {
	wordsQuantity := flag.Int("n", 40000000, "Number of random words in test set.")
	startWord := flag.String("word", "test", "Program will calculate Levenshtein distance from this word.")
	flag.Parse()

	// Create test set, if it was not previously created
	fileName := fmt.Sprintf("test%v.txt", *wordsQuantity)
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		generateTestFileWithLength(*wordsQuantity)
	}

	// Run profiler to measure time and memory consumption
	defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	run(*startWord, file)
	fmt.Println(*wordsQuantity, "words has been sorted succesfully.")
}
