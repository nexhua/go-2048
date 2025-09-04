package game

import (
	"errors"
	"math/rand/v2"
)

func GetRandomCell(cells [][]Cell) (Cell, error) {
	if len(cells) < 1 {
		return Cell{}, errors.New("could not determine grid size")
	}

	emptyCells := make([]Cell, 0, len(cells)*len(cells[0]))

	for i, row := range cells {
		for j := range row {
			cell := &cells[i][j]

			if !cell.isRendered {
				newCell := Cell{pos_x: i, pos_y: j}
				emptyCells = append(emptyCells, newCell)
			}
		}
	}

	if len(emptyCells) <= 0 {
		return Cell{}, errors.New("no empty cells found")
	}

	return emptyCells[rand.IntN(len(emptyCells))], nil
}

// From a boolean array find first empty from given index(START/END)
// Returns error if there are no empty indexes
func FirstEmptyFrom(empties []bool, from int) (int, error) {
	size := len(empties)

	if from < 0 || from >= size {
		return -1, errors.New("invalid from index")
	}

	if from > 0 {
		for i := from; i >= 0; i-- {
			if empties[i] {
				return i, nil
			}
		}
	} else {
		for i, isEmpty := range empties {
			if isEmpty {
				return i, nil
			}
		}
	}

	return -1, errors.New("no empty index found")
}

// Returns a 1D Cell array from 2D grid
// Returned array starts from 0 (meaning UP)
func TakeVerticalSlice(cells [][]Cell, col int) ([]*Cell, error) {
	col_size := len(cells[0])

	if col >= col_size || col < 0 {
		return nil, errors.New("invalid column index")
	}

	v_slice := make([]*Cell, col_size)

	for i, row := range cells {
		v_slice[i] = &row[col]
	}

	return v_slice, nil
}

func ChangeCellState(src *Cell, dst *Cell) {
	tmpRendered := dst.isRendered
	tmpVal := dst.val

	dst.isRendered = src.isRendered
	dst.val = src.val

	src.val = tmpVal
	src.isRendered = tmpRendered
}
