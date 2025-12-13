package ui

import "fmt"

type WeaponChoice struct {
	Name        string
	Description string
	Kind        string // "sword", "projectile", "aoe"
}

var DefaultWeaponChoices = []WeaponChoice{
	{Name: "Sword", Description: "Melee sweep, hits all nearby enemies", Kind: "sword"},
	{Name: "Projectile", Description: "Auto-aim ranged shots", Kind: "projectile"},
	{Name: "Pulse", Description: "AoE damage around you", Kind: "aoe"},
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

func (ws *WeaponSelectScreen) Draw(t *Terminal) {
	w := t.Width()
	h := t.Height()

	title := "=== CHOOSE YOUR WEAPON ==="
	t.WriteStr((w-len(title))/2, h/2-5, title, ColorYellow)

	subtitle := "Up/Down to select, Enter to confirm"
	t.WriteStr((w-len(subtitle))/2, h/2-3, subtitle, ColorDim)

	for i, choice := range ws.Choices {
		color := ColorWhite
		descColor := ColorDim
		prefix := "  "
		if i == ws.selected {
			color = ColorGreen
			descColor = ColorCyan
			prefix = "> "
		}
		name := fmt.Sprintf("%s%s", prefix, choice.Name)
		t.WriteStr((w-len(name))/2, h/2-1+i*2, name, color)
		t.WriteStr((w-len(choice.Description))/2, h/2+i*2, choice.Description, descColor)
	}
}
