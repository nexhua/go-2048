package game

import (
	"errors"
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errors.New("SIGKILL")
	}

	dir, err := GetDirection()
	if err == nil {
		totalNumOfMovements := Move(g, dir)

		if totalNumOfMovements > 0 {
			emptyCell, err := GetRandomCell(g.board.cells)

			if err == nil {
				selectedCell := &g.board.cells[emptyCell.pos_x][emptyCell.pos_y]
				selectedCell.isRendered = true
				selectedCell.val = 2
			}
		}

	}

	// TODO Detect if user can do any movements, if not game over

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
