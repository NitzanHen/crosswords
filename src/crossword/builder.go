package crossword

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/nitzanhen/crossword/src/structure"
)

type Builder struct {
	width, height int
	corpus        Corpus

	debug bool
	calls int
	start time.Time

	listener *func(cw *Crossword)
}

func NewBuilder(width, height int, words []Word, debug bool) Builder {
	return Builder{width, height, NewCorpus(words), debug, 0, time.Time{}, nil}
}

func (builder *Builder) getExactCutRegex(cutData []string) regexp.Regexp {
	return *regexp.MustCompile(
		strings.Join(cutData, ""),
	)
}

func (builder *Builder) getGracefulCutRegex(cutData []string) regexp.Regexp {
	len := len(cutData)

	openingIndex := FirstIndex(
		cutData,
		func(c string) bool { return c != "." },
	)

	trailingIndex := LastIndex(
		cutData,
		func(c string) bool { return c != "." },
	)

	if openingIndex == -1 {
		// Entire cut is periods
		return *regexp.MustCompile(
			fmt.Sprintf("^.{1,%d}$", len),
		)
	}

	var opening, trailing string
	if openingIndex > 0 {
		opening = fmt.Sprintf(".{0,%d}", openingIndex)
	} else {
		opening = ""
	}

	if trailingIndex < len-1 {
		trailing = fmt.Sprintf(".{0,%d}", (len-1)-trailingIndex)
	} else {
		trailing = ""
	}

	inner := strings.Join(cutData[openingIndex:trailingIndex+1], "")

	return *regexp.MustCompile(
		fmt.Sprintf("^%s%s%s$", opening, inner, trailing),
	)
}

func (builder *Builder) getMatchingWords(cw *Crossword, cut Cut) []Word {
	cutData := cw.GetCutData(cut)
	regex := builder.getGracefulCutRegex(cutData)

	matches := builder.corpus.Filter(regex)

	return Filter(
		matches,
		func(w Word) bool { return !cw.IsWordEmbedded(w) },
	)
}

func (builder *Builder) isValidOffset(
	cw *Crossword,
	cut Cut,
	word Word,
	offset int,
) bool {
	subcut := cw.Subcut(cut, offset, offset+len([]rune(word)))
	row, col, o, len := subcut.Row, subcut.Col, subcut.Orientation, subcut.Len

	// Make sure the cells before and after the subcut are not already filled

	if preI, preJ := Move(row, col, o, -1); cw.IsValid(preI, preJ) {
		value := cw.Data[preI][preJ]
		if value != cw.Empty && value != cw.Stop {
			return false
		}
	}
	if postI, postJ := Move(row, col, o, len); cw.IsValid(postI, postJ) {
		value := cw.Data[postI][postJ]
		if value != cw.Empty && value != cw.Stop {
			return false
		}
	}

	// Test Regex for string

	data := cw.GetCutData(subcut)
	regex := builder.getExactCutRegex(data)

	return regex.MatchString(string(word))
}

// func (builder *Builder) debugBuild(word Word, cw *Crossword) {
// 	fmt.Printf("\033[0;1H")
// 	fmt.Printf(" %f Seconds, %d calls, %d cached \n%s\n%s\n\n",
// 		time.Since(builder.start).Seconds(),
// 		builder.calls,
// 		len(builder.corpus.cache),
// 		word,
// 		cw.PrintData(),
// 		//Map(components, func(cuts structure.Set[Cut]) []Cut { return cuts.ToSlice() }),
// 	)
// }

// Attempts to find a suitable crossword with the given cuts embedded
func (builder *Builder) build(cw *Crossword, cuts structure.Set[Cut]) *Crossword {
	builder.calls++

	if cuts.Size() == 0 {
		// Crossword is complete
		return cw
	}

	// Check for inconsistencies

	cutMatchMap := structure.NewOrderedMap[Cut, []Word](cuts.Size())
	for _, cut := range cuts.ToSlice() {
		matches := builder.getMatchingWords(cw, cut)

		if len(matches) == 0 {
			// We have a cut with no matches.
			// It's an inconsistency, stop here.

			//fmt.Printf("INCONSISTENCY: %s %s\n\n", strings.Join(cw.GetCutData(cut), ""), cut.String())
			return nil
		}

		cutMatchMap.Set(cut, matches)
	}

	cutMatches := cutMatchMap.Entries()
	sort.Slice(cutMatches, func(i, j int) bool {
		return len(cutMatches[i].Value) < len(cutMatches[j].Value)
	})

	// Find a suitable next embedding
	for _, entry := range cutMatches {
		cut, matches := entry.Key, entry.Value
		if len(matches) > 100 {
			matches = matches[:100]
		}
		for _, word := range matches {

			// Get valid offsets
			maxOffset := cut.Len - len([]rune(word)) + 1
			validOffsets := Filter(IndexArray(maxOffset), func(offset int) bool {
				return builder.isValidOffset(cw, cut, word, offset)
			})

		EmbeddingLoop:
			for _, offset := range validOffsets {
				next := cw.Copy()

				subcut := cw.Subcut(cut, offset, offset+len([]rune(word)))

				next.Embed(subcut, word)

				subcuts := next.SubcutsOf(
					cuts.ToSlice(),
				)
				nextCuts := structure.SetFromSlice(subcuts)
				nextCuts.Delete(subcut)

				components := GetCutGraph(&next, &nextCuts).Components()
				if len(components) > 1 {
					sort.Slice(components, func(i, j int) bool {
						return components[i].Size() < components[j].Size()
					})
				}

				if builder.listener != nil && builder.calls%2_000 == 0 {
					//builder.debugBuild(word, &next)
					(*builder.listener)(&next)
				}

				// Try filling in each of the components
				result := &next
				for _, component := range components {
					result = builder.build(result, component)

					// if builder.calls > 1000 {
					// 	return result
					// }

					if result == nil {
						// One of the components cant be completed - try the next embedding
						continue EmbeddingLoop
					}

				}

				// Result is non nil, we've completed the embedding
				return result
			}
		}
	}

	return nil
}

func (builder *Builder) Build() *Crossword {
	cw := NewCrossword(builder.width, builder.height)
	cuts := structure.SetFromSlice(cw.GetCuts())

	builder.start = time.Now()

	return builder.build(&cw, cuts)
}

func (builder *Builder) SetListener(listener func(cw *Crossword)) {
	builder.listener = &listener
}
