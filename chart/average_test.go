// chart/average_test.go
package chart

import (
	"math"
	"testing"
)

func TestCalculateWeightedAverage(t *testing.T) {
	votes := map[string]int{
		"Easy":     2,
		"Easy +":   5,
		"Medium -": 8,
		"Medium":   12,
		"Medium +": 6,
		"Hard -":   3,
		"Hard":     1,
	}
	avg := CalculateWeightedAverage(votes)
	// Expected: (2*1.47 + 5*2.06 + 8*2.65 + 12*3.23 + 6*3.83 + 3*4.42 + 1*5.0) / 37
	// = (2.94 + 10.30 + 21.20 + 38.76 + 22.98 + 13.26 + 5.0) / 37
	// = 114.44 / 37 = 3.093...
	expected := 114.44 / 37.0
	if math.Abs(avg-expected) > 0.01 {
		t.Errorf("CalculateWeightedAverage = %f, want %f", avg, expected)
	}
}

func TestCalculateWeightedAverageSingleVote(t *testing.T) {
	votes := map[string]int{"Hell": 1}
	avg := CalculateWeightedAverage(votes)
	if avg != 9.71 {
		t.Errorf("CalculateWeightedAverage = %f, want 9.71", avg)
	}
}

func TestAverageToLabel(t *testing.T) {
	tests := []struct {
		avg      float64
		expected string
	}{
		{0.5, "Easy -"},
		{1.5, "Easy"},
		{3.0, "Medium"},
		{5.0, "Hard"},
		{9.5, "Hell"},
		{10.0, "Hell"},
	}
	for _, tt := range tests {
		label := AverageToLabel(tt.avg)
		if label != tt.expected {
			t.Errorf("AverageToLabel(%f) = %q, want %q", tt.avg, label, tt.expected)
		}
	}
}
