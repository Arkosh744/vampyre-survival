package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

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

// Draw renders the wave reward screen directly onto an Ebitengine image.
func (s *WaveRewardScreen) Draw(dst *ebiten.Image) {
	w := dst.Bounds().Dx()
	h := dst.Bounds().Dy()

	boxW := 280.0
	boxH := 180.0
	boxX := (float64(w) - boxW) / 2
	boxY := (float64(h) - boxH) / 2
	DrawFilledRect(dst, boxX, boxY, boxW, boxH, color.RGBA{R: 15, G: 10, B: 25, A: 240})

	title := "=== WAVE COMPLETE ==="
	titleW := len(title) * 7
	DrawText(dst, (w-titleW)/2, int(boxY)+10, title, ColorToRGBA(ColorMagenta))

	subtitle := "Choose a reward:"
	subW := len(subtitle) * 7
	DrawText(dst, (w-subW)/2, int(boxY)+24, subtitle, ColorToRGBA(ColorWhite))

	startY := int(boxY) + 40
	for i, choice := range s.Choices {
		nameClr := ColorToRGBA(rarityColor(choice.Rarity))
		descClr := ColorToRGBA(ColorDim)
		prefix := "  "
		if i == s.selected {
			nameClr = ColorToRGBA(ColorYellow)
			descClr = ColorToRGBA(ColorCyan)
			prefix = "> "
			DrawFilledRect(dst, boxX+10, float64(startY+i*36-2), boxW-20, 32, color.RGBA{R: 30, G: 20, B: 40, A: 200})
		}
		tag := fmt.Sprintf("[%s]", rarityLabel(choice.Rarity))
		name := fmt.Sprintf("%s%s %s", prefix, choice.Name, tag)
		nameW := len(name) * 7
		DrawText(dst, (w-nameW)/2, startY+i*36, name, nameClr)
		descW := len(choice.Description) * 7
		DrawText(dst, (w-descW)/2, startY+i*36+16, choice.Description, descClr)
	}
}
