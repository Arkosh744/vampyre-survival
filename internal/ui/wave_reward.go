package ui

import (
	"fmt"

	"github.com/arkosh/vampyre-survival/internal/reward"
)

type WaveRewardScreen struct {
	Choices  []reward.Reward
	selected int
}

func NewWaveRewardScreen(choices []reward.Reward) *WaveRewardScreen {
	return &WaveRewardScreen{Choices: choices}
}

func (s *WaveRewardScreen) Up() {
	s.selected--
	if s.selected < 0 {
		s.selected = len(s.Choices) - 1
	}
}

func (s *WaveRewardScreen) Down() {
	s.selected++
	if s.selected >= len(s.Choices) {
		s.selected = 0
	}
}

func (s *WaveRewardScreen) SelectedReward() reward.Reward {
	if s.selected >= 0 && s.selected < len(s.Choices) {
		return s.Choices[s.selected]
	}
	return reward.Reward{}
}

func rarityColor(r reward.RewardRarity) string {
	switch r {
	case reward.RarityCommon:
		return ColorWhite
	case reward.RarityUncommon:
		return ColorGreen
	case reward.RarityRare:
		return ColorCyan
	case reward.RarityEpic:
		return ColorMagenta
	default:
		return ColorWhite
	}
}

func rarityLabel(r reward.RewardRarity) string {
	switch r {
	case reward.RarityCommon:
		return "Common"
	case reward.RarityUncommon:
		return "Uncommon"
	case reward.RarityRare:
		return "Rare"
	case reward.RarityEpic:
		return "Epic"
	default:
		return ""
	}
}

func (s *WaveRewardScreen) Draw(t *Terminal) {
	w := t.Width()
	h := t.Height()

	title := "=== WAVE COMPLETE ==="
	t.WriteStr((w-len(title))/2, h/2-6, title, ColorMagenta)

	subtitle := "Choose a reward:"
	t.WriteStr((w-len(subtitle))/2, h/2-4, subtitle, ColorWhite)

	for i, choice := range s.Choices {
		nameColor := rarityColor(choice.Rarity)
		descColor := ColorDim
		prefix := "  "
		if i == s.selected {
			nameColor = ColorYellow
			descColor = ColorCyan
			prefix = "> "
		}
		tag := fmt.Sprintf("[%s]", rarityLabel(choice.Rarity))
		name := fmt.Sprintf("%s%s %s", prefix, choice.Name, tag)
		t.WriteStr((w-len(name))/2, h/2-2+i*3, name, nameColor)
		t.WriteStr((w-len(choice.Description))/2, h/2-1+i*3, choice.Description, descColor)
	}
}
