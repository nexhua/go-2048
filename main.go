package main

import (
	"bytes"
	"errors"
	"image/color"
	"log"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var CELL_SIZE = 120
var GAP = 10
var CELL_COUNT = 4

type Direction int32

const (
	UP Direction = iota
	RIGHT
	DOWN
	LEFT
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
	fontSource *text.GoTextFaceSource
	fontFace   *text.GoTextFace
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errors.New("SIGKILL")
	}

	dir, err := GetDirection()
	if err == nil {
		Move(g, dir)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	drawBackground(g, screen)
	drawBoard(g, screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// s := ebiten.Monitor().DeviceScaleFactor()
	return int(float64(outsideWidth) * 1), int(float64(outsideHeight))
}

func drawBackground(g *Game, screen *ebiten.Image) {
	bgImg := ebiten.NewImage(g.board.bg.dx, g.board.bg.dy)
	bgImg.Fill(color.NRGBA{0x80, 0x80, 0x80, 0xff})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.board.bg.x), float64(g.board.bg.y))

	screen.DrawImage(bgImg, op)
}

func drawBoard(g *Game, screen *ebiten.Image) {
	cellImg := ebiten.NewImage(CELL_SIZE, CELL_SIZE)
	cellImg.Fill(color.White)

	textImg := ebiten.NewImage(CELL_SIZE, CELL_SIZE)
	textImg.Fill(color.Black)

	for _, row := range g.board.cells {
		for _, cell := range row {
			colour := GetColor(cell.val)
			op := &ebiten.DrawImageOptions{}
			op.ColorScale.ScaleWithColor(colour)
			txtOp := &text.DrawOptions{}
			op.GeoM.Translate(float64(cell.x), float64(cell.y))

			screen.DrawImage(cellImg, op)

			if cell.isRendered {

				txtOp.ColorScale.ScaleWithColor(color.Black)
				DrawCenteredText(screen, g.fontFace, strconv.Itoa(cell.val), cell.x+CELL_SIZE/2, cell.y+CELL_SIZE/2)
			}
		}
	}
}

func initGame() *Game {
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
	g := Game{board: b}

	cell := &g.board.cells[0][0]
	cell.isRendered = true
	cell.val = 2

	cell = &g.board.cells[1][0]
	cell.isRendered = true
	cell.val = 2

	cell = &g.board.cells[0][1]
	cell.isRendered = true
	cell.val = 2

	cell = &g.board.cells[0][2]
	cell.isRendered = true
	cell.val = 4

	cell = &g.board.cells[0][3]
	cell.isRendered = true
	cell.val = 8

	cell = &g.board.cells[1][1]
	cell.isRendered = true
	cell.val = 8

	cell = &g.board.cells[2][2]
	cell.isRendered = true
	cell.val = 16

	cell = &g.board.cells[3][3]
	cell.isRendered = true
	cell.val = 32

	cell = &g.board.cells[3][2]
	cell.isRendered = true
	cell.val = 64

	cell = &g.board.cells[3][3]
	cell.isRendered = true
	cell.val = 128

	fontData, err := os.ReadFile("Roboto-Thin.ttf")
	if err != nil {
		panic("font read fail")
	}

	g.fontSource, err = text.NewGoTextFaceSource(bytes.NewReader(fontData))

	if err != nil {
		panic("font fail")
	}

	g.fontFace = &text.GoTextFace{
		Source: g.fontSource,
		Size:   float64(CELL_SIZE) / 2,
	}

	return &g
}

func main() {
	x, y := ebiten.Monitor().Size()
	ebiten.SetWindowSize(x/2, y/2)
	ebiten.SetWindowTitle("2048!")
	if err := ebiten.RunGame(initGame()); err != nil {
		if err.Error() == "SIGKILL" {
			return
		}
		log.Fatal(err)
	}
}
