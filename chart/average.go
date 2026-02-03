package chart

func CalculateWeightedAverage(votes map[string]int) float64 {
	var totalWeight float64
	var totalVotes int

	for level, count := range votes {
		if midpoint, ok := DifficultyMidpoints[level]; ok {
			totalWeight += midpoint * float64(count)
			totalVotes += count
		}
	}

	if totalVotes == 0 {
		return 0
	}
	return totalWeight / float64(totalVotes)
}

func AverageToLabel(avg float64) string {
	for _, level := range DifficultyLevels {
		r := DifficultyRanges[level]
		if avg >= r.Lower {
			if r.UpperInclusive {
				if avg <= r.Upper {
					return level
				}
			} else {
				if avg < r.Upper {
					return level
				}
			}
		}
	}
	return "Hell"
}
