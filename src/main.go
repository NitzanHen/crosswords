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
	Width         int
	Height        int
	Success       bool
	Time          float64
	Calls         int
	Failures      int
}

const (
	BATCH_SIZE  = 50
	NUM_BATCHES = 50
	WIDTH       = 5
	HEIGHT      = 5
	TIMEOUT     = 10 * time.Second
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
			builder := crossword.NewBuilder(WIDTH, HEIGHT, shuffled, false)
			// builder.SetListener(func(cw *crossword.Crossword) {
			// 	fmt.Printf("\033[2;0H")
			// 	fmt.Printf("\n%s\n\n", cw.PrintData())
			// })

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
				results.Add(BuildResult{res, startingWords, WIDTH, HEIGHT, true, elapsed, builder.Calls, builder.Failures})

				fmt.Printf("Success in %f seconds:\n%s\n", elapsed, res.PrintData())

			case <-time.After(TIMEOUT):
				elapsed := time.Since(start).Seconds()
				results.Add(BuildResult{nil, startingWords, WIDTH, HEIGHT, false, elapsed, builder.Calls, builder.Failures})

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

	os.Mkdir("./output", 0755)
	filename := fmt.Sprintf("./output/result-%d-%d.json", runId, batchId)

	if err := os.WriteFile(filename, data, 0644); err != nil {
		log.Fatalf("%v", err)
	}
}
