// chart/difficulty_test.go
package chart

import "testing"

func TestDifficultyLevelsOrdered(t *testing.T) {
	expected := []string{
		"Easy -", "Easy", "Easy +",
		"Medium -", "Medium", "Medium +",
		"Hard -", "Hard", "Hard +",
		"Very Hard -", "Very Hard", "Very Hard +",
		"Extreme -", "Extreme", "Extreme +",
		"Hell",
	}
	if len(DifficultyLevels) != 16 {
		t.Fatalf("expected 16 levels, got %d", len(DifficultyLevels))
	}
	for i, level := range expected {
		if DifficultyLevels[i] != level {
			t.Errorf("level %d: expected %q, got %q", i, level, DifficultyLevels[i])
		}
	}
}

func TestDifficultyIndex(t *testing.T) {
	tests := []struct {
		name  string
		index int
	}{
		{"Easy -", 0},
		{"Medium", 4},
		{"Hell", 15},
	}
	for _, tt := range tests {
		idx, ok := DifficultyIndex(tt.name)
		if !ok {
			t.Errorf("DifficultyIndex(%q) not found", tt.name)
			continue
		}
		if idx != tt.index {
			t.Errorf("DifficultyIndex(%q) = %d, want %d", tt.name, idx, tt.index)
		}
	}
}

func TestDifficultyIndexInvalid(t *testing.T) {
	_, ok := DifficultyIndex("NotALevel")
	if ok {
		t.Error("expected DifficultyIndex to return false for invalid level")
	}
}
