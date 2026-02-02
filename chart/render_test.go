// chart/render_test.go
package chart

import (
	"testing"
)

func TestRenderChart(t *testing.T) {
	votes := map[string]int{
		"Medium":   10,
		"Medium +": 5,
		"Hard -":   3,
	}

	imgData, err := RenderChart(votes)
	if err != nil {
		t.Fatalf("RenderChart error: %v", err)
	}
	if len(imgData) == 0 {
		t.Error("expected non-empty image data")
	}

	// Check WebP magic bytes
	if len(imgData) < 4 || string(imgData[0:4]) != "RIFF" {
		t.Error("expected WebP format (RIFF header)")
	}
}
