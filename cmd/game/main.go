package main

import (
	"fmt"
	"os"

	"github.com/arkosh/vampyre-survival/internal/game"
)

func main() {
	g, err := game.NewGame()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to start: %v\n", err)
		os.Exit(1)
	}
	g.Run()
}
