package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type WeaponChoice struct {
	Name        string
	Description string
	Kind        string // "sword", "projectile", "aoe"
}

var DefaultWeaponChoices = []WeaponChoice{
	{Name: "Sword", Description: "Melee sweep, hits all nearby enemies", Kind: "sword"},
	{Name: "Projectile", Description: "Auto-aim ranged shots", Kind: "projectile"},
	{Name: "Pulse", Description: "AoE damage around you", Kind: "aoe"},
	{Name: "Lightning", Description: "Chain lightning jumps between enemies", Kind: "lightning"},
	{Name: "Orbital", Description: "Spinning orbs deal constant damage", Kind: "orbital"},
}

type WeaponSelectScreen struct {
	Choices  []WeaponChoice
	selected int
}

func NewWeaponSelectScreen() *WeaponSelectScreen {
	return &WeaponSelectScreen{Choices: DefaultWeaponChoices}
}

func (ws *WeaponSelectScreen) Up() {
	ws.selected--
	if ws.selected < 0 {
		ws.selected = len(ws.Choices) - 1
	}
}

func (ws *WeaponSelectScreen) Down() {
	ws.selected++
	if ws.selected >= len(ws.Choices) {
		ws.selected = 0
	}
}

func (ws *WeaponSelectScreen) Selected() WeaponChoice {
	return ws.Choices[ws.selected]
}

// Draw renders the weapon selection screen directly onto an Ebitengine image.
func (ws *WeaponSelectScreen) Draw(dst *ebiten.Image) {
	w := dst.Bounds().Dx()
	h := dst.Bounds().Dy()

	// Fill background
	dst.Fill(color.RGBA{R: 10, G: 5, B: 15, A: 255})

	// Box
	boxW := 280.0
	boxH := 240.0
	boxX := (float64(w) - boxW) / 2
	boxY := (float64(h) - boxH) / 2
	DrawFilledRect(dst, boxX, boxY, boxW, boxH, color.RGBA{R: 15, G: 10, B: 25, A: 240})

	title := "=== CHOOSE YOUR WEAPON ==="
	titleW := len(title) * 7
	DrawText(dst, (w-titleW)/2, int(boxY)+10, title, ColorToRGBA(ColorYellow))

	hint := "Up/Down to select, Enter to confirm"
	hintW := len(hint) * 7
	DrawText(dst, (w-hintW)/2, int(boxY)+24, hint, ColorToRGBA(ColorDim))

	startY := int(boxY) + 40
	for i, choice := range ws.Choices {
		nameClr := ColorToRGBA(ColorWhite)
		descClr := ColorToRGBA(ColorDim)
		prefix := "  "
		if i == ws.selected {
			nameClr = ColorToRGBA(ColorGreen)
			descClr = ColorToRGBA(ColorCyan)
			prefix = "> "
			DrawFilledRect(dst, boxX+10, float64(startY+i*36-2), boxW-20, 32, color.RGBA{R: 20, G: 40, B: 20, A: 200})
		}
		name := prefix + choice.Name
		nameW := len(name) * 7
		DrawText(dst, (w-nameW)/2, startY+i*36, name, nameClr)
		descW := len(choice.Description) * 7
		DrawText(dst, (w-descW)/2, startY+i*36+16, choice.Description, descClr)
	}
}
