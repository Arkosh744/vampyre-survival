package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/arkosh/vampyre-survival/internal/game"
)

func main() {
	g := game.NewGame()
	ebiten.SetWindowSize(1920, 1080)
	ebiten.SetWindowTitle("Vampyre Survival")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
