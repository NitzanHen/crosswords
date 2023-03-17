package crossword

import (
	"fmt"
	"nitzanhen/crossword/src/structure"
	"strings"
)

type CutMatrix struct {
	Width  int
	Height int

	// Indicates an Empty cell.
	Empty string
	// Indicates a cell that cannot have a value, i.e. cannot be in any cut.
	Stop string

	Data [][]string
}

func NewCutMatrix(width, height int, empty, stop string) CutMatrix {
	data := MakeMatrix(
		height,
		width,
		func(i, j int) string { return empty },
	)

	return CutMatrix{width, height, empty, stop, data}
}

func Move(row, col int, o Orientation, step int) (i, j int) {
	switch o {
	case HORIZONTAL:
		return row, col + step
	case VERTICAL:
		return row + step, col
	}

	return -1, -1
}

func IsInCut(i, j int, cut Cut) bool {
	switch cut.Orientation {
	case HORIZONTAL:
		return i == cut.Row && (cut.Col <= j && j < cut.Col+cut.Len)
	case VERTICAL:
		return j == cut.Col && (cut.Row <= i && i < cut.Row+cut.Len)
	}

	return false
}

type CellData struct {
	i, j  int
	value string
}

func (mat *CutMatrix) IterateCut(cut Cut) []CellData {
	row, col, o, len := cut.Row, cut.Col, cut.Orientation, cut.Len

	data := make([]CellData, len)
	for k := 0; k < len; k++ {
		i, j := Move(row, col, o, k)
		data[k] = CellData{i, j, mat.Data[i][j]}
	}

	return data
}

func (mat *CutMatrix) GetCutData(cut Cut) []string {
	return Map(
		mat.IterateCut(cut),
		func(data CellData) string { return data.value },
	)
}

func (mat *CutMatrix) getSubcuts(cut Cut) []Cut {
	startRow, startCol, o := cut.Row, cut.Col, cut.Orientation
	len := 0

	subcuts := structure.List[Cut]{}

	for _, data := range mat.IterateCut(cut) {
		i, j, value := data.i, data.j, data.value

		if value == mat.Stop {
			if len > 0 {
				subcuts.Add(Cut{startRow, startCol, o, len})
			}
			startRow, startCol = Move(i, j, o, 1)
			len = 0
		} else {
			len++
		}
	}

	if len > 0 {
		subcuts.Add(Cut{startRow, startCol, o, len})
	}

	return subcuts.ToSlice()
}

func (mat *CutMatrix) SubcutsOf(cuts []Cut) []Cut {
	subcuts := structure.List[Cut]{}
	for _, cut := range cuts {
		for _, subcut := range mat.getSubcuts(cut) {
			if subcut.Len > 1 {
				subcuts.Add(subcut)
			}
		}
	}

	return subcuts.ToSlice()
}

func (mat *CutMatrix) GetCuts() []Cut {
	width, height := mat.Width, mat.Height

	initialCuts := make([]Cut, width+height)
	for row := 0; row < height; row++ {
		initialCuts[row] = Cut{row, 0, HORIZONTAL, width}
	}
	for col := 0; col < width; col++ {
		initialCuts[height+col] = Cut{0, col, VERTICAL, height}
	}

	return mat.SubcutsOf(initialCuts)
}

func (mat *CutMatrix) IsValid(row, col int) bool {
	return (0 <= row && row < mat.Height) &&
		(0 <= col && col < mat.Width)
}

func (mat *CutMatrix) Set(i, j int, value string) error {
	if !mat.IsValid(i, j) {
		return fmt.Errorf(
			"invalid placement: Out of bounds. Got coords (%d, %d), dimensions are %dx%d",
			i, j, mat.Width, mat.Height,
		)
	}

	currentValue := mat.Data[i][j]
	if currentValue != mat.Empty && currentValue != value {
		return fmt.Errorf(
			"invalid placement: cannot write value %v at coords (%d, %d), position already populated by %v",
			value, i, j, currentValue,
		)
	}

	mat.Data[i][j] = value

	return nil
}

func (mat *CutMatrix) Subcut(cut Cut, start, end int) Cut {
	row, col := Move(cut.Row, cut.Col, cut.Orientation, start)

	return Cut{row, col, cut.Orientation, end - start}
}

func (mat *CutMatrix) FillIn(data []string, cut Cut) error {
	n := len(data)
	if n > cut.Len {
		return fmt.Errorf(
			"invalid embedding: data of length %d cannot be embedded in cut of length %d",
			n, cut.Len,
		)
	}

	cutData := mat.IterateCut(cut)
	for k, cell := range cutData[:n] {
		err := mat.Set(cell.i, cell.j, data[k])

		if err != nil {
			return err
		}
	}

	return nil
}

func (mat *CutMatrix) DoCutsMeet(cut1, cut2 Cut) bool {
	// Two cuts of different orientation meet iff they are equal.
	if cut1.Orientation == cut2.Orientation {
		return cut1 == cut2
	}

	// Different orientations
	// Intersection point is necessarily in the row of the horizontal cut and the col of the vertical
	var intI, intJ int
	if cut1.Orientation == HORIZONTAL {
		intI, intJ = cut1.Row, cut2.Col
	} else { // s2.o == HORIZONTAL
		intI, intJ = cut2.Row, cut1.Col
	}

	return IsInCut(intI, intJ, cut1) && IsInCut(intI, intJ, cut2)
}

func (mat *CutMatrix) PrintData() string {
	rowStrings := Map(mat.Data, func(row []string) string {
		return "| " + strings.Join(row, " | ") + " |"
	})

	return strings.Join(rowStrings, "\n")
}

func (mat *CutMatrix) Copy() CutMatrix {
	var copy CutMatrix
	copy.Width = mat.Width
	copy.Height = mat.Height
	copy.Empty = mat.Empty
	copy.Stop = mat.Stop

	copy.Data = MakeMatrix(
		mat.Height, mat.Width,
		func(i, j int) string { return mat.Data[i][j] },
	)

	return copy
}
