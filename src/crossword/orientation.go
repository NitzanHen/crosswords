package crossword

import "fmt"

type Orientation int

const (
	HORIZONTAL Orientation = iota
	VERTICAL   Orientation = iota
)

func (o Orientation) String() string {
	switch o {
	case HORIZONTAL:
		return "horizontal"
	case VERTICAL:
		return "vertical"
	}

	panic(fmt.Sprintf("Invalid Orientation %d", int(o)))
}
