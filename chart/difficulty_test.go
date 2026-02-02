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

func TestParseHexColor(t *testing.T) {
	tests := []struct {
		hex     string
		r, g, b uint8
	}{
		{"#66ff66", 0x66, 0xff, 0x66},
		{"#990000", 0x99, 0x00, 0x00},
		{"#ffb300", 0xff, 0xb3, 0x00},
	}
	for _, tt := range tests {
		r, g, b := ParseHexColor(tt.hex)
		if r != tt.r || g != tt.g || b != tt.b {
			t.Errorf("ParseHexColor(%q) = (%d, %d, %d), want (%d, %d, %d)",
				tt.hex, r, g, b, tt.r, tt.g, tt.b)
		}
	}
}
