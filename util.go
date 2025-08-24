package main

import (
	"errors"
	"fmt"
	"image/color"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	BEIGE       = color.NRGBA{0xfa, 0xf8, 0xef, 0xff}
	LIGHT_BROWN = color.NRGBA{0xcd, 0xc1, 0xb4, 0xff}

	// Tile colors
	LIGHT_TAN    = color.NRGBA{0xee, 0xe4, 0xda, 0xff} // 2
	TAN          = color.NRGBA{0xed, 0xe0, 0xc8, 0xff} // 4
	ORANGE_LIGHT = color.NRGBA{0xf2, 0xb1, 0x79, 0xff} // 8
	ORANGE       = color.NRGBA{0xf5, 0x95, 0x63, 0xff} // 16
	RED_ORANGE   = color.NRGBA{0xf6, 0x7c, 0x5f, 0xff} // 32
	RED          = color.NRGBA{0xf6, 0x5e, 0x3b, 0xff} // 64
	YELLOW_LIGHT = color.NRGBA{0xed, 0xcf, 0x72, 0xff} // 128
	YELLOW       = color.NRGBA{0xed, 0xcc, 0x61, 0xff} // 256
	GOLD_LIGHT   = color.NRGBA{0xed, 0xc8, 0x50, 0xff} // 512
	GOLD         = color.NRGBA{0xed, 0xc5, 0x3f, 0xff} // 1024
	GOLD_DEEP    = color.NRGBA{0xed, 0xc2, 0x2e, 0xff} // 2048

	// For larger values
	DARK_GRAY = color.NRGBA{0x3c, 0x3a, 0x32, 0xff}

	// Text colors
	TEXT_DARK  = color.NRGBA{0x77, 0x6e, 0x65, 0xff} // for 2, 4
	TEXT_LIGHT = color.NRGBA{0xf9, 0xf6, 0xf2, 0xff} // for others
)

// TileColors maps tile values to colors
var TileColors = map[int]color.NRGBA{
	2:    LIGHT_TAN,
	4:    TAN,
	8:    ORANGE_LIGHT,
	16:   ORANGE,
	32:   RED_ORANGE,
	64:   RED,
	128:  YELLOW_LIGHT,
	256:  YELLOW,
	512:  GOLD_LIGHT,
	1024: GOLD,
	2048: GOLD_DEEP,
}

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

			mergeCount, err := MergeSlice(row, len(row)-1)

			if err == nil && mergeCount > 0 {
				ShiftRight(row)
			}
		}
	case LEFT:
		for _, row := range g.board.cells {
			ShiftLeft(row)
			mergeCount, err := MergeSlice(row, 0)

			if err == nil && mergeCount > 0 {
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
				ShiftUp(v_slice)
				mergeCount, err := MergeSliceRef(v_slice, 0)

				if err == nil && mergeCount > 0 {
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
				ShiftDown(v_slice)
				mergeCount, err := MergeSliceRef(v_slice, len(v_slice)-1)

				if err == nil && mergeCount > 0 {
					ShiftDown(v_slice)
				}
			}
		}
	}

}

func GetRandomCell(cells [][]Cell) (Cell, error) {
	emptyCells := make([]Cell, 0, CELL_COUNT*CELL_COUNT)

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

// Merge a horizontal slice
// Merge direction is 0 or SIZE-1, any other value will be rejected
func MergeSlice(slice []Cell, to int) (int, error) {
	mergeCount := 0

	if !(to == 0 || to == CELL_COUNT-1) {
		return 0, errors.New("invalid to argument")
	}

	if len(slice) < 2 {
		return 0, errors.New("slice is too small to merge")
	}

	// merge to start of slice
	if to == 0 {
		for i := 0; i < len(slice)-1; i++ {
			lc, rc := &slice[i], &slice[i+1]

			if (lc.isRendered && rc.isRendered) && (lc.val == rc.val) {
				lc.val += rc.val
				rc.isRendered = false
				rc.val = 0
				mergeCount++
			}
		}
	} else {
		for i := len(slice) - 1; i > 0; i-- {
			lc, rc := &slice[i-1], &slice[i]

			if (lc.isRendered && rc.isRendered) && (lc.val == rc.val) {
				rc.val += lc.val
				lc.isRendered = false
				lc.val = 0
				mergeCount++
			}
		}
	}

	return mergeCount, nil
}

// Merge a vertical (accepts a ref slice)
// Merge direction is 0 or SIZE-1, any other value will be rejected
func MergeSliceRef(slice []*Cell, to int) (int, error) {
	mergeCount := 0

	if !(to == 0 || to == CELL_COUNT-1) {
		fmt.Println("invalid to argument")
		return 0, errors.New("invalid to argument")
	}

	if len(slice) < 2 {
		fmt.Println("slice is too small to merge")
		return 0, errors.New("slice is too small to merge")
	}

	// merge to start of slice
	if to == 0 {
		for i := 0; i < len(slice)-1; i++ {
			lc, rc := slice[i], slice[i+1]

			if (lc.isRendered && rc.isRendered) && (lc.val == rc.val) {
				lc.val += rc.val
				rc.isRendered = false
				rc.val = 0
				mergeCount++
			}
		}
	} else {
		for i := len(slice) - 1; i > 0; i-- {
			lc, rc := slice[i-1], slice[i]

			if (lc.isRendered && rc.isRendered) && (lc.val == rc.val) {
				rc.val += lc.val
				lc.isRendered = false
				lc.val = 0
				mergeCount++
			}
		}
	}

	return mergeCount, nil
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

// Returns a 1D Cell array from 2D grid
// Returned array starts from 0 (meaning UP)
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

func GetColor(val int) color.NRGBA {
	if col, ok := TileColors[val]; ok {
		return col
	}

	return DARK_GRAY
}
