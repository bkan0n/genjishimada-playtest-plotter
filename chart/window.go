// chart/window.go
package chart

const minWindowSize = 5

// CalculateWindow returns the min and max difficulty indices to display
func CalculateWindow(votes map[string]int) (minIdx, maxIdx int) {
	// Find min and max indices with votes
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

	// No votes case
	if maxVoted == -1 {
		return 0, minWindowSize - 1
	}

	// Expand by 1 on each side
	minIdx = minVoted - 1
	maxIdx = maxVoted + 1

	// Clamp to valid range
	if minIdx < 0 {
		minIdx = 0
	}
	if maxIdx > len(DifficultyLevels)-1 {
		maxIdx = len(DifficultyLevels) - 1
	}

	// Ensure minimum window size
	windowSize := maxIdx - minIdx + 1
	if windowSize < minWindowSize {
		deficit := minWindowSize - windowSize
		expandLeft := deficit / 2
		expandRight := deficit - expandLeft

		minIdx -= expandLeft
		maxIdx += expandRight

		// Re-clamp and adjust
		if minIdx < 0 {
			maxIdx -= minIdx // add the overflow to maxIdx
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
