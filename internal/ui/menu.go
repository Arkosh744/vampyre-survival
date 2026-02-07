package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type MenuOption int

const (
	MenuStart MenuOption = iota
	MenuQuit
)

type Menu struct {
	selected int
	options  []string
}

func NewMenu() *Menu {
	return &Menu{
		options: []string{"Start Game", "Quit"},
	}
}

func (m *Menu) Up() {
	m.selected--
	if m.selected < 0 {
		m.selected = len(m.options) - 1
	}
}

func (m *Menu) Down() {
	m.selected++
	if m.selected >= len(m.options) {
		m.selected = 0
	}
}

func (m *Menu) Selected() MenuOption {
	return MenuOption(m.selected)
}

// Draw renders the main menu directly onto an Ebitengine image.
func (m *Menu) Draw(dst *ebiten.Image) {
	w := dst.Bounds().Dx()
	h := dst.Bounds().Dy()

	// Fill background
	dst.Fill(color.RGBA{R: 10, G: 5, B: 15, A: 255})

	// Title centered
	title := "VAMPYRE SURVIVAL"
	titleW := len(title) * 7
	DrawText(dst, (w-titleW)/2, h/2-40, title, ColorToRGBA(ColorBoldRed))

	// Subtitle
	sub := "Arrow Keys to move, ESC to pause"
	subW := len(sub) * 7
	DrawText(dst, (w-subW)/2, h/2-22, sub, ColorToRGBA(ColorDim))

	// Menu options
	for i, opt := range m.options {
		clr := ColorToRGBA(ColorWhite)
		prefix := "  "
		if i == m.selected {
			clr = ColorToRGBA(ColorYellow)
			prefix = "> "
			// Selection highlight bar
			text := prefix + opt
			textW := len(text) * 7
			barX := float64((w-textW)/2) - 4
			barY := float64(h/2+i*18) - 2
			DrawFilledRect(dst, barX, barY, float64(textW)+8, 16, color.RGBA{R: 40, G: 40, B: 60, A: 255})
		}
		text := prefix + opt
		textW := len(text) * 7
		DrawText(dst, (w-textW)/2, h/2+i*18, text, clr)
	}
}
