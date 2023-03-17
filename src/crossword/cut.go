package crossword

import (
	"fmt"
	"nitzanhen/crossword/src/structure"
)

type Cut struct {
	Row         int
	Col         int
	Orientation Orientation
	Len         int
}

func (cut *Cut) String() string {
	return fmt.Sprintf("Cut{%d, %d, %s, %d}",
		cut.Row,
		cut.Col,
		cut.Orientation.String(),
		cut.Len,
	)
}

func (cut *Cut) Copy() Cut {
	return Cut{cut.Row, cut.Col, cut.Orientation, cut.Len}
}

type CutStatus int

const (
	EMPTY    CutStatus = iota
	PARTIAL  CutStatus = iota
	EMBEDDED CutStatus = iota
)

func (status CutStatus) String() string {
	switch status {
	case EMPTY:
		return "empty"
	case PARTIAL:
		return "partial"
	case EMBEDDED:
		return "embedded"
	}

	return "INVALID CUT STATUS"
}

// Returns the non-embedded cuts of the crossword as a graph,
// With edges representing intersection
func GetCutGraph(cw *Crossword, cuts *structure.Set[Cut]) *structure.Graph[Cut] {
	graph := structure.NewGraph(*cuts)

	cutsSlice := cuts.ToSlice()

	for i, cut1 := range cutsSlice {
		for _, cut2 := range cutsSlice[i:] {
			if cw.DoCutsMeet(cut1, cut2) {
				graph.Connect(cut1, cut2)
			}
		}
	}

	return &graph
}
