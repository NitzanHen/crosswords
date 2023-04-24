package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/nitzanhen/crossword/src/crossword"
	"github.com/nitzanhen/crossword/src/structure"
)

type BuildResult struct {
	Result        *crossword.Crossword
	StartingWords []crossword.Word
	Success       bool
	Time          float64
}

const (
	BATCH_SIZE  = 200
	NUM_BATCHES = 15
)

func main() {
	words := getWords()

	// Clears the screen
	fmt.Print("\033[H\033[2J")

	runId := rand.Intn(100_000)

	for i := 0; i < NUM_BATCHES; i++ {
		results := structure.List[BuildResult]{}

		for j := 0; j < BATCH_SIZE; j++ {
			resultChan := make(chan *crossword.Crossword, 1)

			shuffled := shuffle(words)
			builder := crossword.NewBuilder(6, 6, shuffled, false)

			startingWords := shuffled[:10]

			fmt.Printf("Iteration %d-%d: \nFirst words: %v\n", i, j, startingWords)

			var start time.Time

			go func() {
				start = time.Now()
				resultChan <- builder.Build()
			}()

			select {

			case res := <-resultChan:
				elapsed := time.Since(start).Seconds()
				results.Add(BuildResult{res, startingWords, true, elapsed})

				fmt.Printf("Success in %f seconds:\n%s\n", elapsed, res.PrintData())

			case <-time.After(10 * time.Second):
				elapsed := time.Since(start)
				results.Add(BuildResult{nil, startingWords, false, elapsed.Seconds()})

				fmt.Println("Timed out.")
			}
		}

		writeResult(results.ToSlice(), runId, i)

	}

}

func getWords() []crossword.Word {
	raw, err := os.ReadFile("./hebrew.json")
	if err != nil {
		log.Fatalf("Unable to read file: %v", err)
	}

	var words []crossword.Word
	err = json.Unmarshal(raw, &words)
	if err != nil {
		log.Fatalf("Unable to unmarshal: %v", err)
	}

	for i, w := range words {
		s := string(w)
		s = strings.ToLower(s)
		s = strings.ReplaceAll(s, "ף", "פ")
		s = strings.ReplaceAll(s, "ץ", "צ")
		s = strings.ReplaceAll(s, "ך", "כ")
		s = strings.ReplaceAll(s, "ן", "נ")
		s = strings.ReplaceAll(s, "ם", "מ")
		words[i] = crossword.Word(s)
	}

	fmt.Printf("Read %d words\n", len(words))

	return words
}

func shuffle[T any](items []T) []T {
	perm := rand.Perm(len(items))
	shuffled := make([]T, len(items))

	for i, j := range perm {
		shuffled[j] = items[i]
	}

	return shuffled
}

func writeResult(results []BuildResult, runId, batchId int) {
	data, _ := json.MarshalIndent(results, "", "  ")

	filename := fmt.Sprintf("./output/result-%d-%d.json", runId, batchId)

	os.WriteFile(filename, data, 0644)
}
