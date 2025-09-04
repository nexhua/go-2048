package main

import (
	"log"

	"mkoca/2048/src/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	x, y := ebiten.Monitor().Size()
	ebiten.SetWindowSize(x/2, y/2)
	ebiten.SetWindowTitle("2048!")
	if err := ebiten.RunGame(game.InitGame()); err != nil {
		if err.Error() == "SIGKILL" {
			return
		}
		log.Fatal(err)
	}
}
