package game

import (
	"errors"
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errors.New("SIGKILL")
	}

	switch g.status {
	case RUNNING:
		if IsGameFinished(g.board.cells) {
			g.status = FINISHED
			break
		}

		if IsGameOver(g.board.cells) {
			g.status = GAME_OVER
			break
		}

		// Only accept input if there are no animations running
		if !HasRunningAnimation(g) {
			dir, err := GetDirection()
			if err == nil {
				totalNumOfMovements, totalMergeScore := Move(g, dir)

				if totalNumOfMovements > 0 {
					emptyCell, err := GetRandomCell(g.board.cells)

					if err == nil {
						selectedCell := &g.board.cells[emptyCell.pos_x][emptyCell.pos_y]
						selectedCell.isRendered = false
						selectedCell.val = 0
						selectedCell.animation = CreateCellAnimation(*selectedCell, CREATE_CELL_ANIMATION_DURATION)
					}

					g.score += totalMergeScore
				}
			}
		}

	case FINISHED:
		pressedKeys := inpututil.AppendJustPressedKeys(nil)
		if len(pressedKeys) > 0 {
			ResetGame(g)
		}
	case GAME_OVER:
		pressedKeys := inpututil.AppendJustPressedKeys(nil)
		if len(pressedKeys) > 0 {
			ResetGame(g)
		}
	default:
		return errors.New("unhandled game status reached. exiting")
	}

	// TODO Add ctrl z to go back
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.status {
	case RUNNING:
		drawBackground(g, screen)
		drawBoard(g, screen)
		drawScoreboard(g, screen, g.board.bg.x+g.board.bg.dx, g.board.bg.y)
	case FINISHED:
		drawBackground(g, screen)
		drawBoard(g, screen)
		drawScoreboard(g, screen, g.board.bg.x+g.board.bg.dx, g.board.bg.y)
		drawOverlay(g, screen)
		drawAfterGameText(g, screen, "Congratulations!")
	case GAME_OVER:
		drawBackground(g, screen)
		drawBoard(g, screen)
		drawScoreboard(g, screen, g.board.bg.x+g.board.bg.dx, g.board.bg.y)
		drawOverlay(g, screen)
		drawAfterGameText(g, screen, "Game Over!")
	default:
		panic("TODO")
	}
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

			if cell.animation != nil {
				if cell.animation.GetStatus() == ANIM_FINISHED {
					c := &g.board.cells[cell.pos_x][cell.pos_y]
					// TODO
					// Decide if animation should handle the logic on complete or renderer
					// c.animation.OnFinish(cell.animation, g)
					c.animation = nil
					c.isRendered = true
					c.val = 2
				} else {
					drawAnimation(screen, g, &cell)
					animErr := cell.animation.Step()

					if animErr != nil {
						panic("step called on finished animation. check your logic")
					}
				}
			} else if cell.isRendered {
				txtOp.ColorScale.ScaleWithColor(color.Black)
				DrawCenteredText(screen, g.fontFace, strconv.Itoa(cell.val), cell.x+CELL_SIZE/2, cell.y+CELL_SIZE/2, txtOp)
			}
		}
	}
}

func drawScoreboard(g *Game, screen *ebiten.Image, x int, y int) {
	x_offset := float64(x) + float64(CELL_SIZE)/2
	y_offset := float64(y) + float64(GAP)
	w := CELL_SIZE * 2
	h := CELL_SIZE
	scoreImg := ebiten.NewImage(w, h)
	scoreImg.Fill(DARK_GRAY)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x_offset, y_offset)

	screen.DrawImage(scoreImg, op)

	txtOp := &text.DrawOptions{}
	txtOp.ColorScale.ScaleWithColor(TEXT_DARK)
	DrawCenteredText(screen, g.fontFace, "SCORE", int(x_offset)+w/2, int(y_offset)+h/4, txtOp)
	txtOp = &text.DrawOptions{}
	DrawCenteredText(screen, g.fontFace, strconv.Itoa(g.score), int(x_offset)+w/2, int(y_offset)+3*h/4, txtOp)
}

func drawOverlay(g *Game, screen *ebiten.Image) {
	overlayImage := ebiten.NewImage(g.board.bg.dx, g.board.bg.dy)
	overlayImage.Fill(OVERLAY_BACKGROUND)

	overlayOpt := &ebiten.DrawImageOptions{}
	overlayOpt.GeoM.Translate(float64(g.board.bg.x), float64(g.board.bg.y))

	screen.DrawImage(overlayImage, overlayOpt)
}

func drawAnimation(screen *ebiten.Image, g *Game, cell *Cell) {

	if cell.animation == nil {
		panic("invalid state. drawAnimation should not be called when animation does not exists")
	}

	switch cell.animation.GetType() {
	case ANIM_CREATE_CELL:
		ca := cell.animation.(*CreateAnimation)
		c := ca.position

		cellImg := ebiten.NewImage(ca.currentSize, ca.currentSize)
		cellImg.Fill(GetColor(2))

		cx := c.x + CELL_SIZE/2 // center x
		cy := c.y + CELL_SIZE/2 // center y

		cx = cx - ca.currentSize/2 // center x offsetted by current size
		cy = cy - ca.currentSize/2 // center y offsetted by current size

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(cx), float64(cy))
		screen.DrawImage(cellImg, op)
	default:
		panic("unreachable")
	}
}

func drawAfterGameText(g *Game, screen *ebiten.Image, message string) {
	txtOp := &text.DrawOptions{}
	txtOp.ColorScale.ScaleWithColor(color.Black)
	cx := g.board.bg.x + g.board.bg.dx/2
	cy := g.board.bg.y + g.board.bg.dy/2
	lh := g.fontFace.Metrics().HAscent
	DrawCenteredText(screen, g.fontFace, message, cx, cy, txtOp)

	txtOp = &text.DrawOptions{}
	txtOp.ColorScale.ScaleWithColor(DARK_GRAY)
	// TODO
	// Ugly way to print text at different sizes
	// Refactor this and accept size as parameter for DrawCenteredText
	// Calculate predefined font size based on screen (like big, small, medium etc..)
	tmp := g.fontFace.Size
	g.fontFace.Size = tmp - 20
	DrawCenteredText(screen, g.fontFace, "Press any key to reset", cx, cy+int(lh), txtOp)
	g.fontFace.Size = tmp
}

func DrawCenteredText(screen *ebiten.Image, fontFace *text.GoTextFace, s string, cx int, cy int, txtOp *text.DrawOptions) {
	tx, ty := text.Measure(s, fontFace, 0)
	txtOp.GeoM.Translate(float64(cx)-tx/2, float64(cy)-ty/2)
	text.Draw(screen, s, fontFace, txtOp)
}

func CalculateActualCellPosition(start_x int, start_y int, pos_x int, pos_y int, cell_s int, gap_s int) (int, int) {
	x := GAP + start_x + (pos_x * cell_s) + (pos_x * GAP)
	y := GAP + start_y + (pos_y * cell_s) + (pos_y * GAP)
	return x, y
}
