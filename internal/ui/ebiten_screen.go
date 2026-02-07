package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	CellW = 10
	CellH = 18
)

// EbitenScreen wraps an *ebiten.Image and provides cell-dimension info
// used by the camera system. All actual rendering goes through pixel-based
// draw functions (DrawText, DrawFilledRect, etc.).
type EbitenScreen struct {
	img    *ebiten.Image
	width  int // cell count horizontal
	height int // cell count vertical
}

// NewEbitenScreen creates a screen backed by an off-screen ebiten.Image.
func NewEbitenScreen(pixelW, pixelH int) *EbitenScreen {
	return &EbitenScreen{
		img:    ebiten.NewImage(pixelW, pixelH),
		width:  pixelW / CellW,
		height: pixelH / CellH,
	}
}

func (s *EbitenScreen) Width() int  { return s.width }
func (s *EbitenScreen) Height() int { return s.height }

// SetTarget replaces the underlying ebiten.Image (called each frame from Draw).
func (s *EbitenScreen) SetTarget(img *ebiten.Image) {
	s.img = img
	bounds := img.Bounds()
	s.width = bounds.Dx() / CellW
	s.height = bounds.Dy() / CellH
}

// Image returns the underlying ebiten.Image.
func (s *EbitenScreen) Image() *ebiten.Image {
	return s.img
}
