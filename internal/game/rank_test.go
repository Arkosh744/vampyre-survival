package game

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_CalculateRank_S(t *testing.T) {
	r := CalculateRank(RunStats{
		Time:      250,
		Kills:     400,
		HPPercent: 1.0,
	})
	require.Equal(t, RankS, r)
}

func Test_CalculateRank_A(t *testing.T) {
	r := CalculateRank(RunStats{
		Time:      350,
		Kills:     300,
		HPPercent: 0.8,
	})
	require.Equal(t, RankA, r)
}

func Test_CalculateRank_B(t *testing.T) {
	r := CalculateRank(RunStats{
		Time:      450,
		Kills:     150,
		HPPercent: 0.5,
	})
	require.Equal(t, RankB, r)
}

func Test_CalculateRank_C(t *testing.T) {
	// Time <540 → 20pts, kills 100 → 10pts, HP 0.1 → 3pts = 33 → C
	r := CalculateRank(RunStats{
		Time:      500,
		Kills:     100,
		HPPercent: 0.1,
	})
	require.Equal(t, RankC, r)
}

func Test_CalculateRank_D(t *testing.T) {
	r := CalculateRank(RunStats{
		Time:      700,
		Kills:     10,
		HPPercent: 0.05,
	})
	require.Equal(t, RankD, r)
}

func Test_Rank_String(t *testing.T) {
	require.Equal(t, "S", RankS.String())
	require.Equal(t, "A", RankA.String())
	require.Equal(t, "B", RankB.String())
	require.Equal(t, "C", RankC.String())
	require.Equal(t, "D", RankD.String())
}

func Test_CalculateRank_SlowButAlive(t *testing.T) {
	// Time ≥600 → 0pts, HP 100% → 30pts, kills 300 → 30pts = 60 → B
	r := CalculateRank(RunStats{
		Time:      650,
		Kills:     300,
		HPPercent: 1.0,
	})
	require.Equal(t, RankB, r)
}

func Test_CalculateRank_FastButHurt(t *testing.T) {
	// Time <300 → 40pts, HP 10% → 3pts, kills 200 → 20pts = 63 → B
	r := CalculateRank(RunStats{
		Time:      280,
		Kills:     200,
		HPPercent: 0.1,
	})
	require.Equal(t, RankB, r)
}
