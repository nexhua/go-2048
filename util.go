package main

import (
	"errors"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
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

func Move(g *Game, d Direction) {
	switch d {
	case RIGHT:
		for _, row := range g.board.cells {
			ShiftRight(row)
		}
	case LEFT:
		for _, row := range g.board.cells {
			ShiftLeft(row)
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
				ShiftUp(v_slice)
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
				ShiftDown(v_slice)
			}
		}
	}

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

func ShiftRight(row []Cell) {
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
				// MoveCellValues(cells, cell.pos_x, cell.pos_y, cell_to.pos_x, cell_to.pos_y)
				ChangeCellState(cell, cell_to)
				empties[right_empty] = false
				empties[i] = true
			}
		}
	}
}

func ShiftLeft(row []Cell) {
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
			}
		}
	}
}

func ShiftUp(col []*Cell) {
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
			}
		}
	}
}

func ShiftDown(col []*Cell) {
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
			}
		}
	}
}

func TakeVerticalSlice(cells [][]Cell, col int) ([]*Cell, error) {
	col_size := len(cells[0])

	if col >= col_size {
		return nil, errors.New("invalid column index")
	}

	v_slice := make([]*Cell, col_size)

	for i, row := range cells {
		v_slice[i] = &row[col]
	}

	return v_slice, nil
}

func MoveCellValues(cells [][]Cell, src_x int, src_y int, dst_x int, dst_y int) {
	from := &cells[src_x][src_y]
	to := &cells[dst_x][dst_y]

	tmpIsRendered := to.isRendered
	tmpVal := to.val

	to.isRendered = from.isRendered
	to.val = from.val

	from.isRendered = tmpIsRendered
	from.val = tmpVal
}

func ChangeCellState(src *Cell, dst *Cell) {
	tmpRendered := dst.isRendered
	tmpVal := dst.val

	dst.isRendered = src.isRendered
	dst.val = src.val

	src.val = tmpVal
	src.isRendered = tmpRendered
}

func DrawCenteredText(screen *ebiten.Image, fontFace *text.GoTextFace, s string, cx int, cy int) {

	// bounds, _ := font.BoundString(fontFace, s)
	// dx := (bounds.Max.X - bounds.Min.X).Round()
	// dy := (bounds.Max.Y - bounds.Min.Y).Round()
	// x, y := cx-bounds.Min.X.Round()-dx/2, cy-bounds.Min.Y.Round()-dy/2
	tx, ty := text.Measure(s, fontFace, 0)
	txtOp := &text.DrawOptions{}
	txtOp.GeoM.Translate(float64(cx)-tx/2, float64(cy)-ty/2)
	txtOp.ColorScale.ScaleWithColor(color.Black)
	text.Draw(screen, s, fontFace, txtOp)
}

func CalculateActualCellPosition(start_x int, start_y int, pos_x int, pos_y int, cell_s int, gap_s int) (int, int) {
	x := GAP + start_x + (pos_x * cell_s) + (pos_x * GAP)
	y := GAP + start_y + (pos_y * cell_s) + (pos_y * GAP)
	return x, y
}

func FormatCell(cell Cell) string {
	return fmt.Sprintf("Cell{pos_x:%d, pos_y:%d, x:%d, y:%d, val:%d, isRendered:%t}", cell.pos_x, cell.pos_y, cell.x, cell.y, cell.val, cell.isRendered)
}
