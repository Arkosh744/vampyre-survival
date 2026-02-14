package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/basicfont"
)

var defaultFace *text.GoXFace = text.NewGoXFace(basicfont.Face7x13)

// DrawText draws a string onto an Ebitengine image at pixel coordinates.
func DrawText(dst *ebiten.Image, x, y int, s string, clr color.RGBA) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(clr)
	text.Draw(dst, s, defaultFace, op)
}
