package main

import (
	"log"
	"runtime"

	"mkoca/2048/src/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	x, y := ebiten.Monitor().Size()

	// TODO
	// WASM should also scale properly to users screen
	if runtime.GOOS == "js" && runtime.GOARCH == "wasm" {
		ebiten.SetWindowTitle("2048!")
		ebiten.SetWindowSize(800, 600)
	} else {
		ebiten.SetWindowSize(x/2, y/2)
		ebiten.SetWindowTitle("2048!")
		ebiten.SetTPS(game.TARGET_TPS)
	}

	if err := ebiten.RunGame(game.InitGame()); err != nil {
		if err.Error() == "SIGKILL" {
			return
		}
		log.Fatal(err)
	}
}
