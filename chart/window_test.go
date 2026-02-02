// chart/window_test.go
package chart

import "testing"

func TestCalculateWindow(t *testing.T) {
	tests := []struct {
		name   string
		votes  map[string]int
		minIdx int
		maxIdx int
	}{
		{
			name:   "middle range",
			votes:  map[string]int{"Medium": 5, "Medium +": 3, "Hard -": 2},
			minIdx: 3,  // Medium - (one before Medium)
			maxIdx: 7,  // Hard (one after Hard -)
		},
		{
			name:   "single vote expands to minimum 5",
			votes:  map[string]int{"Medium": 1},
			minIdx: 2,  // Easy +
			maxIdx: 6,  // Hard -
		},
		{
			name:   "at start clamped",
			votes:  map[string]int{"Easy -": 5},
			minIdx: 0,  // Easy -
			maxIdx: 4,  // Medium
		},
		{
			name:   "at end clamped",
			votes:  map[string]int{"Hell": 5},
			minIdx: 11, // Very Hard +
			maxIdx: 15, // Hell
		},
		{
			name:   "wide range shows all voted",
			votes:  map[string]int{"Easy": 1, "Extreme +": 1},
			minIdx: 0,  // Easy - (one before Easy)
			maxIdx: 15, // Hell (one after Extreme +, clamped)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			minIdx, maxIdx := CalculateWindow(tt.votes)
			if minIdx != tt.minIdx || maxIdx != tt.maxIdx {
				t.Errorf("CalculateWindow() = (%d, %d), want (%d, %d)",
					minIdx, maxIdx, tt.minIdx, tt.maxIdx)
			}
		})
	}
}
