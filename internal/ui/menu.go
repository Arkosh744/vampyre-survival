package ui

import "fmt"

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

func (m *Menu) Draw(t *Terminal) {
	w := t.Width()
	h := t.Height()

	title := "VAMPYRE SURVIVAL"
	subtitle := "Arrow Keys to move, ESC to pause"

	t.WriteStr((w-len(title))/2, h/2-4, title, ColorBoldRed)
	t.WriteStr((w-len(subtitle))/2, h/2-2, subtitle, ColorDim)

	for i, opt := range m.options {
		color := ColorWhite
		prefix := "  "
		if i == m.selected {
			color = ColorYellow
			prefix = "> "
		}
		text := fmt.Sprintf("%s%s", prefix, opt)
		t.WriteStr((w-len(text))/2, h/2+i, text, color)
	}
}
