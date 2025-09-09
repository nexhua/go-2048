package game

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

var CELL_SIZE = 120
var GAP = 10
var CELL_COUNT = 4
var TARGET_TPS = 60
var CREATE_CELL_ANIMATION_DURATION = TARGET_TPS / 8

type GameStatus int32

const (
	RUNNING GameStatus = iota
	FINISHED
	GAME_OVER
	ANIMATING
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
	animations []Animation
}

func FormatCell(cell Cell) string {
	return fmt.Sprintf("Cell{pos_x:%d, pos_y:%d, x:%d, y:%d, val:%d, isRendered:%t}", cell.pos_x, cell.pos_y, cell.x, cell.y, cell.val, cell.isRendered)
}

func InitGame() *Game {
	screen_x, screen_y := ebiten.WindowSize()
	grid_x, grid_y := CELL_COUNT*CELL_SIZE+(CELL_COUNT+1)*GAP, CELL_COUNT*CELL_SIZE+(CELL_COUNT+1)*GAP
	bg_x_offset, bg_y_offset := (screen_x-grid_x)/2, (screen_y-grid_y)/2

	// TODO
	// This section could be better
	// For now good enough for basic wasm support
	if runtime.GOOS == "js" && runtime.GOARCH == "wasm" {
		bg_x_offset = 0
		bg_y_offset = 0
	}

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

	anims := make([]Animation, 0)

	b := Board{bg: background, cells: cells}
	g := Game{board: b, status: RUNNING, animations: anims}

	emptyCell, err := GetRandomCell(g.board.cells)

	if err == nil {
		selectedCell := &g.board.cells[emptyCell.pos_x][emptyCell.pos_y]
		selectedCell.isRendered = true
		selectedCell.val = 2
	}

	// TODO
	// Current font fetching is a bit of a mess
	// We need different types of fonts for different purposes and also scale it properly based on screen size
	// Also load a better font for WASM
	if runtime.GOARCH == "wasm" && runtime.GOOS == "js" {
		fontData, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
		if err != nil {
			log.Fatal("failed to load default WASM font:", err)
		}

		g.fontSource = fontData
		g.fontFace = &text.GoTextFace{
			Source: g.fontSource,
			Size:   float64(CELL_SIZE) / 4 * ebiten.Monitor().DeviceScaleFactor(),
		}
	} else {
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

func IsGameOver(cells [][]Cell) bool {
	hasEmpty := HasEmptyCell(cells)

	if hasEmpty {
		return !hasEmpty
	} else {
		return !HasPossibleMerge(cells)
	}
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
