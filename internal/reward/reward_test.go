package reward

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_AllRewards_Count20(t *testing.T) {
	require.Len(t, AllRewards(), 20)
}

func Test_RewardPool_Wave1_OnlyCommon(t *testing.T) {
	p := NewRewardPool()
	for i := 0; i < 50; i++ {
		choices := p.GetChoices(1)
		for _, c := range choices {
			require.Equal(t, RarityCommon, c.Rarity, "wave 1 should only yield Common rewards")
		}
		// Reset taken so pool doesn't deplete
		p.taken = make(map[RewardID]bool)
	}
}

func Test_RewardPool_Wave7_HasRare(t *testing.T) {
	p := NewRewardPool()
	foundRare := false
	for i := 0; i < 200; i++ {
		choices := p.GetChoices(7)
		for _, c := range choices {
			if c.Rarity == RarityRare {
				foundRare = true
			}
		}
		p.taken = make(map[RewardID]bool)
	}
	require.True(t, foundRare, "wave 7 should sometimes yield Rare rewards")
}

func Test_RewardPool_Returns2(t *testing.T) {
	p := NewRewardPool()
	choices := p.GetChoices(15)
	require.Len(t, choices, 2)
}

func Test_RewardPool_NoDuplicates(t *testing.T) {
	p := NewRewardPool()
	for i := 0; i < 100; i++ {
		choices := p.GetChoices(15)
		if len(choices) == 2 {
			require.NotEqual(t, choices[0].ID, choices[1].ID, "two choices must be different")
		}
		p.taken = make(map[RewardID]bool)
	}
}

func Test_RewardPool_MarkTaken_Excludes(t *testing.T) {
	p := NewRewardPool()
	// Take all common rewards
	for _, r := range allRewards {
		if r.Rarity == RarityCommon {
			p.MarkTaken(r.ID)
		}
	}
	// Wave 1 should return nil (no commons left, others locked by wave)
	choices := p.GetChoices(1)
	require.Nil(t, choices)
}
