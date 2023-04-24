package crossword

import (
	"fmt"

	"github.com/nitzanhen/crossword/src/structure"
)

type Crossword struct {
	CutMatrix

	Embeddings []CutWithWord
}

func NewCrossword(width, height int) Crossword {
	var cw Crossword

	cw.CutMatrix = NewCutMatrix(width, height, ".", "1")
	cw.Embeddings = []CutWithWord{}

	return cw
}

func (cw *Crossword) IsWordEmbedded(word Word) bool {
	for _, cutword := range cw.Embeddings {
		if word == cutword.Word {
			return true
		}
	}

	return false
}

func (cw *Crossword) IsCutEmbedded(cut Cut) bool {
	for _, cutword := range cw.Embeddings {
		if cut == cutword.Cut {
			return true
		}
	}

	return false
}

func (cw *Crossword) Embed(cut Cut, word Word) error {
	if cw.IsWordEmbedded(word) {
		return fmt.Errorf("Word %v is already embedded", word)
	}

	chars := Chars(string(word))

	err := cw.FillIn(chars, cut)
	if err != nil {
		return err
	}

	row, col, o, len := cut.Row, cut.Col, cut.Orientation, cut.Len

	if preRow, preCol := Move(row, col, o, -1); cw.IsValid(preRow, preCol) {
		cw.Set(preRow, preCol, cw.Stop)
	}
	if postRow, postCol := Move(row, col, o, len); cw.IsValid(postRow, postCol) {
		cw.Set(postRow, postCol, cw.Stop)
	}

	cw.Embeddings = append(cw.Embeddings, CutWithWord{cut, word})

	return nil
}

// Returns the non-embedded cuts of the crossword as a graph,
// With edges representing intersection
func (cw *Crossword) GetCutGraph() structure.Graph[Cut] {
	cuts := Filter(
		cw.GetCuts(),
		func(cut Cut) bool { return !cw.IsCutEmbedded(cut) },
	)

	graph := structure.NewGraph(
		structure.SetFromSlice(cuts),
	)

	for i, cut1 := range cuts {
		for _, cut2 := range cuts[i:] {
			if cw.DoCutsMeet(cut1, cut2) {
				graph.Connect(cut1, cut2)
			}
		}
	}

	return graph
}

func (cw *Crossword) Copy() Crossword {
	var copy Crossword

	copy.CutMatrix = cw.CutMatrix.Copy()
	copy.Embeddings = Map(
		cw.Embeddings,
		func(cutword CutWithWord) CutWithWord { return cutword.Copy() },
	)

	return copy
}

type CutWithWord struct {
	Cut  Cut
	Word Word
}

func (cutword *CutWithWord) Copy() CutWithWord {
	return CutWithWord{cutword.Cut.Copy(), cutword.Word}
}
