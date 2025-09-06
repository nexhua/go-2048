package game

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var CELL_SIZE = 120
var GAP = 10
var CELL_COUNT = 4

type GameStatus int32

const (
	RUNNING GameStatus = iota
	FINISHED
	GAME_OVER
)

type Cell struct {
	x          int
	y          int
	pos_x      int
	pos_y      int
	val        int
	isRendered bool
}

type Background struct {
	x  int
	y  int
	dx int
	dy int
}

type Board struct {
	bg    Background
	cells [][]Cell
}

type Game struct {
	board      Board
	score      int
	fontSource *text.GoTextFaceSource
	fontFace   *text.GoTextFace
	status     GameStatus
}

func FormatCell(cell Cell) string {
	return fmt.Sprintf("Cell{pos_x:%d, pos_y:%d, x:%d, y:%d, val:%d, isRendered:%t}", cell.pos_x, cell.pos_y, cell.x, cell.y, cell.val, cell.isRendered)
}

// TODO add tests for merges
// Merge a horizontal slice
// Merge direction is 0 or SIZE-1, any other value will be rejected
func MergeSlice(slice []Cell, to int) (int, error) {
	mergeScore := 0

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
				mergeScore += lc.val * 2
				lc.val += rc.val
				rc.isRendered = false
				rc.val = 0
			}
		}
	} else {
		for i := len(slice) - 1; i > 0; i-- {
			lc, rc := &slice[i-1], &slice[i]

			if (lc.isRendered && rc.isRendered) && (lc.val == rc.val) {
				mergeScore += lc.val * 2
				rc.val += lc.val
				lc.isRendered = false
				lc.val = 0
			}
		}
	}

	return mergeScore, nil
}

// Merge a vertical (accepts a ref slice)
// Merge direction is 0 or SIZE-1, any other value will be rejected
func MergeSliceRef(slice []*Cell, to int) (int, error) {
	mergeScore := 0

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
				mergeScore += lc.val * 2
				lc.val += rc.val
				rc.isRendered = false
				rc.val = 0
			}
		}
	} else {
		for i := len(slice) - 1; i > 0; i-- {
			lc, rc := slice[i-1], slice[i]

			if (lc.isRendered && rc.isRendered) && (lc.val == rc.val) {
				mergeScore += lc.val * 2
				rc.val += lc.val
				lc.isRendered = false
				lc.val = 0
			}
		}
	}

	return mergeScore, nil
}

func InitGame() *Game {
	screen_x, screen_y := ebiten.WindowSize()
	grid_x, grid_y := CELL_COUNT*CELL_SIZE+(CELL_COUNT+1)*GAP, CELL_COUNT*CELL_SIZE+(CELL_COUNT+1)*GAP
	bg_x_offset, bg_y_offset := (screen_x-grid_x)/2, (screen_y-grid_y)/2

	background := Background{
		x:  bg_x_offset,
		y:  bg_y_offset,
		dx: grid_x,
		dy: grid_y,
	}

	cells := make([][]Cell, CELL_COUNT)
	for i := range cells {
		cells[i] = make([]Cell, CELL_COUNT)

		for j := range cells[i] {
			x, y := CalculateActualCellPosition(background.x, background.y, j, i, CELL_SIZE, GAP)
			cells[i][j] = Cell{pos_x: i, pos_y: j, x: x, y: y, isRendered: false, val: 0}
		}
	}

	b := Board{bg: background, cells: cells}
	g := Game{board: b, status: RUNNING}

	emptyCell, err := GetRandomCell(g.board.cells)

	if err == nil {
		selectedCell := &g.board.cells[emptyCell.pos_x][emptyCell.pos_y]
		selectedCell.isRendered = true
		selectedCell.val = 2
	}

	fontData, err := os.ReadFile("assets/Roboto-Thin.ttf")
	if err != nil {
		panic("font read fail")
	}

	g.fontSource, err = text.NewGoTextFaceSource(bytes.NewReader(fontData))

	if err != nil {
		panic("font fail")
	}

	g.fontFace = &text.GoTextFace{
		Source: g.fontSource,
		Size:   float64(CELL_SIZE) / 4 * ebiten.Monitor().DeviceScaleFactor(),
	}

	return &g
}

func IsGameFinished(cells [][]Cell) bool {
	for i := 0; i < len(cells); i++ {
		for j := 0; j < len(cells[i]); j++ {
			c := cells[i][j]
			if c.isRendered && c.val == 2048 {
				return true
			}
		}
	}

	return false
}

func ResetGame(g *Game) {
	ResetBoard(&g.board)
	g.score = 0
	g.status = RUNNING

	emptyCell, err := GetRandomCell(g.board.cells)

	if err == nil {
		selectedCell := &g.board.cells[emptyCell.pos_x][emptyCell.pos_y]
		selectedCell.isRendered = true
		selectedCell.val = 2
	}
}

func ResetBoard(b *Board) {
	for i, row := range b.cells {
		for j := range row {
			c := &b.cells[i][j]
			c.isRendered = false
			c.val = 0
		}
	}
}
