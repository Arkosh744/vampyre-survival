package game

type Rank int

const (
	RankS Rank = iota
	RankA
	RankB
	RankC
	RankD
)

func (r Rank) String() string {
	switch r {
	case RankS:
		return "S"
	case RankA:
		return "A"
	case RankB:
		return "B"
	case RankC:
		return "C"
	default:
		return "D"
	}
}

type RunStats struct {
	Time       float64
	Kills      int
	HPPercent  float64
	Level      int
	DamageDealt int
	BestStreak int
}

func CalculateRank(s RunStats) Rank {
	pts := 0

	// Time: 40pts max — faster is better
	switch {
	case s.Time < 300:
		pts += 40
	case s.Time < 420:
		pts += 30
	case s.Time < 540:
		pts += 20
	case s.Time < 600:
		pts += 10
	}

	// HP: 30pts max
	hp := int(s.HPPercent * 30)
	if hp > 30 {
		hp = 30
	}
	if hp < 0 {
		hp = 0
	}
	pts += hp

	// Kills: 30pts max (1pt per 10 kills)
	killPts := s.Kills / 10
	if killPts > 30 {
		killPts = 30
	}
	pts += killPts

	switch {
	case pts >= 85:
		return RankS
	case pts >= 65:
		return RankA
	case pts >= 45:
		return RankB
	case pts >= 25:
		return RankC
	default:
		return RankD
	}
}
