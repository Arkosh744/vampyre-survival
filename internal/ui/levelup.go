package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/arkosh/vampyre-survival/internal/skill"
)

type LevelUpScreen struct {
	Choices  []skill.Upgrade
	selected int
}

func NewLevelUpScreen(choices []skill.Upgrade) *LevelUpScreen {
	return &LevelUpScreen{Choices: choices}
}

func (l *LevelUpScreen) Up() {
	l.selected--
	if l.selected < 0 {
		l.selected = len(l.Choices) - 1
	}
}

func (l *LevelUpScreen) Down() {
	l.selected++
	if l.selected >= len(l.Choices) {
		l.selected = 0
	}
}

func (l *LevelUpScreen) Selected() int {
	return l.selected
}

func (l *LevelUpScreen) SelectedUpgrade() skill.Upgrade {
	if l.selected >= 0 && l.selected < len(l.Choices) {
		return l.Choices[l.selected]
	}
	return skill.Upgrade{}
}

// Draw renders the level-up screen directly onto an Ebitengine image.
func (l *LevelUpScreen) Draw(dst *ebiten.Image) {
	w := dst.Bounds().Dx()
	h := dst.Bounds().Dy()

	// Semi-transparent overlay box in center
	boxW := 250.0
	boxH := 180.0
	boxX := (float64(w) - boxW) / 2
	boxY := (float64(h) - boxH) / 2
	DrawFilledRect(dst, boxX, boxY, boxW, boxH, color.RGBA{R: 10, G: 10, B: 20, A: 230})

	// Title
	title := "=== LEVEL UP! ==="
	titleW := len(title) * 7
	DrawText(dst, (w-titleW)/2, int(boxY)+10, title, ColorToRGBA(ColorYellow))

	subtitle := "Choose an upgrade:"
	subW := len(subtitle) * 7
	DrawText(dst, (w-subW)/2, int(boxY)+24, subtitle, ColorToRGBA(ColorWhite))

	startY := int(boxY) + 40
	for i, choice := range l.Choices {
		nameClr := ColorToRGBA(ColorWhite)
		descClr := ColorToRGBA(ColorDim)
		prefix := "  "
		if i == l.selected {
			nameClr = ColorToRGBA(ColorGreen)
			descClr = ColorToRGBA(ColorCyan)
			prefix = "> "
			// Highlight bar
			DrawFilledRect(dst, boxX+10, float64(startY+i*35-2), boxW-20, 30, color.RGBA{R: 20, G: 40, B: 20, A: 200})
		}
		name := prefix + choice.Name
		nameW := len(name) * 7
		DrawText(dst, (w-nameW)/2, startY+i*35, name, nameClr)
		descW := len(choice.Description) * 7
		DrawText(dst, (w-descW)/2, startY+i*35+16, choice.Description, descClr)
	}
}
