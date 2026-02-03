package chart

const minWindowSize = 5

func CalculateWindow(votes map[string]int) (minIdx, maxIdx int) {
	minVoted := len(DifficultyLevels)
	maxVoted := -1

	for level, count := range votes {
		if count > 0 {
			if idx, ok := DifficultyIndex(level); ok {
				if idx < minVoted {
					minVoted = idx
				}
				if idx > maxVoted {
					maxVoted = idx
				}
			}
		}
	}

	if maxVoted == -1 {
		return 0, minWindowSize - 1
	}

	minIdx = minVoted - 1
	maxIdx = maxVoted + 1

	if minIdx < 0 {
		minIdx = 0
	}
	if maxIdx > len(DifficultyLevels)-1 {
		maxIdx = len(DifficultyLevels) - 1
	}

	windowSize := maxIdx - minIdx + 1
	if windowSize < minWindowSize {
		deficit := minWindowSize - windowSize
		expandLeft := deficit / 2
		expandRight := deficit - expandLeft

		minIdx -= expandLeft
		maxIdx += expandRight

		if minIdx < 0 {
			maxIdx -= minIdx
			minIdx = 0
		}
		if maxIdx > len(DifficultyLevels)-1 {
			minIdx -= maxIdx - (len(DifficultyLevels) - 1)
			maxIdx = len(DifficultyLevels) - 1
		}
		if minIdx < 0 {
			minIdx = 0
		}
	}

	return minIdx, maxIdx
}
