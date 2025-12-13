package ui

import (
	"fmt"

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

func (l *LevelUpScreen) Draw(t *Terminal) {
	w := t.Width()
	h := t.Height()

	title := "=== LEVEL UP! ==="
	t.WriteStr((w-len(title))/2, h/2-5, title, ColorYellow)

	subtitle := "Choose an upgrade:"
	t.WriteStr((w-len(subtitle))/2, h/2-3, subtitle, ColorWhite)

	for i, choice := range l.Choices {
		color := ColorWhite
		descColor := ColorDim
		prefix := "  "
		if i == l.selected {
			color = ColorGreen
			descColor = ColorCyan
			prefix = "> "
		}
		name := fmt.Sprintf("%s%s", prefix, choice.Name)
		t.WriteStr((w-len(name))/2, h/2-1+i*2, name, color)
		t.WriteStr((w-len(choice.Description))/2, h/2+i*2, choice.Description, descColor)
	}
}
