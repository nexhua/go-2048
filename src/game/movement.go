package game

import (
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Direction int32

const (
	UP Direction = iota
	RIGHT
	DOWN
	LEFT
)

func GetDirection() (Direction, error) {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		return UP, nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		return RIGHT, nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		return DOWN, nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		return LEFT, nil
	}

	return UP, errors.New("arrow keys are not pressed")
}

// Moves cells for a given direction
// Returns the number of movements
func Move(g *Game, d Direction) (int, int) {
	totalNumOfMovements := 0
	totalMergeScore := 0
	switch d {
	case RIGHT:
		for _, row := range g.board.cells {
			totalNumOfMovements += ShiftRight(row)

			mergeScore, err := MergeSlice(row, len(row)-1)

			if err == nil && mergeScore > 0 {
				totalMergeScore += mergeScore
				ShiftRight(row)
			}
		}
	case LEFT:
		for _, row := range g.board.cells {
			totalNumOfMovements += ShiftLeft(row)
			mergeScore, err := MergeSlice(row, 0)

			if err == nil && mergeScore > 0 {
				totalMergeScore += mergeScore
				ShiftLeft(row)
			}
		}
	case UP:
		row_s := len(g.board.cells)
		if row_s < 1 {
			break
		}

		col_s := len(g.board.cells[0])

		for j := range col_s {
			v_slice, err := TakeVerticalSlice(g.board.cells, j)

			if err == nil {
				totalNumOfMovements += ShiftUp(v_slice)
				mergeScore, err := MergeSliceRef(v_slice, 0)

				if err == nil && mergeScore > 0 {
					totalMergeScore += mergeScore
					ShiftUp(v_slice)
				}
			}
		}
	case DOWN:
		row_s := len(g.board.cells)
		if row_s < 1 {
			break
		}

		col_s := len(g.board.cells[0])

		for j := range col_s {
			v_slice, err := TakeVerticalSlice(g.board.cells, j)

			if err == nil {
				totalNumOfMovements += ShiftDown(v_slice)
				mergeScore, err := MergeSliceRef(v_slice, len(v_slice)-1)

				if err == nil && mergeScore > 0 {
					totalMergeScore += mergeScore
					ShiftDown(v_slice)
				}
			}
		}
	}

	return totalNumOfMovements, totalMergeScore

}

func ShiftRight(row []Cell) int {
	numOfMovements := 0
	size := len(row)
	empties := make([]bool, len(row))

	for i := size - 1; i >= 0; i-- {
		cell := &row[i]
		if !cell.isRendered {
			empties[i] = true
		} else {
			right_empty, err := FirstEmptyFrom(empties, size-1)
			if err == nil && right_empty > i {
				cell_to := &row[right_empty]
				ChangeCellState(cell, cell_to)
				empties[right_empty] = false
				empties[i] = true
				numOfMovements++
			}
		}
	}

	return numOfMovements
}

func ShiftLeft(row []Cell) int {
	numOfMovements := 0
	size := len(row)
	empties := make([]bool, len(row))

	for i := 0; i < size; i++ {
		cell := &row[i]
		if !cell.isRendered {
			empties[i] = true
		} else {
			left_empty, err := FirstEmptyFrom(empties, 0)
			if err == nil && left_empty < i {
				cell_to := &row[left_empty]
				ChangeCellState(cell, cell_to)
				empties[left_empty] = false
				empties[i] = true
				numOfMovements++
			}
		}
	}

	return numOfMovements
}

func ShiftUp(col []*Cell) int {
	numOfMovements := 0
	size := len(col)
	empties := make([]bool, size)

	for i, cell := range col {
		if !cell.isRendered {
			empties[i] = true
		} else {
			up_empty, err := FirstEmptyFrom(empties, 0)

			if err == nil && up_empty < i {
				cell_to := col[up_empty]
				ChangeCellState(cell, cell_to)
				empties[up_empty] = false
				empties[i] = true
				numOfMovements++
			}
		}
	}

	return numOfMovements
}

func ShiftDown(col []*Cell) int {
	numOfMovements := 0
	size := len(col)
	empties := make([]bool, size)

	for i := size - 1; i >= 0; i-- {
		cell := col[i]
		if !cell.isRendered {
			empties[i] = true
		} else {
			down_empty, err := FirstEmptyFrom(empties, size-1)

			if err == nil && down_empty > i {
				cell_to := col[down_empty]
				ChangeCellState(cell, cell_to)
				empties[down_empty] = false
				empties[i] = true
				numOfMovements++
			}
		}
	}

	return numOfMovements
}
