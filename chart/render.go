// chart/render.go
package chart

import (
	"bytes"
	"image"

	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/ungerik/go-cairo"
)

const (
	CanvasWidth  = 1400
	CanvasHeight = 700
	LeftMargin   = 80
	RightMargin  = 60
	TopMargin    = 80
	BottomMargin = 120
	BarGap       = 12
	BarRadius    = 16
)

var (
	BackgroundColor = [3]float64{0.168, 0.176, 0.192} // #2b2d31
	TextColor       = [3]float64{1.0, 1.0, 1.0}       // white
)

// RenderChart generates a WebP chart image from vote data
func RenderChart(votes map[string]int) ([]byte, error) {
	// Create Cairo surface
	surface := cairo.NewSurface(cairo.FORMAT_ARGB32, CanvasWidth, CanvasHeight)
	defer surface.Finish()

	// Fill background
	surface.SetSourceRGB(BackgroundColor[0], BackgroundColor[1], BackgroundColor[2])
	surface.Rectangle(0, 0, CanvasWidth, CanvasHeight)
	surface.Fill()

	// Calculate data
	avg := CalculateWeightedAverage(votes)
	avgLabel := AverageToLabel(avg)
	minIdx, maxIdx := CalculateWindow(votes)

	// Draw chart elements
	drawYAxisLines(surface, votes, minIdx, maxIdx) // Draw grid lines first (behind bars)
	drawBars(surface, votes, minIdx, maxIdx)
	drawXAxisLabels(surface, minIdx, maxIdx)
	drawYAxis(surface, votes, minIdx, maxIdx)
	drawVoteCounts(surface, votes, minIdx, maxIdx)
	drawAverageLine(surface, avg, avgLabel, minIdx, maxIdx)

	// Convert to image.RGBA
	img := surfaceToImage(surface)

	// Encode to WebP
	var buf bytes.Buffer
	options, _ := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 90)
	if err := webp.Encode(&buf, img, options); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func surfaceToImage(surface *cairo.Surface) *image.RGBA {
	width := surface.GetWidth()
	height := surface.GetHeight()
	data := surface.GetData()

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Cairo uses ARGB32 (premultiplied), need to convert to RGBA
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			i := (y*width + x) * 4
			// Cairo ARGB32: B, G, R, A (little endian)
			b := data[i]
			g := data[i+1]
			r := data[i+2]
			a := data[i+3]
			j := (y*width + x) * 4
			img.Pix[j] = r
			img.Pix[j+1] = g
			img.Pix[j+2] = b
			img.Pix[j+3] = a
		}
	}
	return img
}

const (
	ShadowOffsetX = 4
	ShadowOffsetY = 4
	ShadowAlpha   = 0.3
)

func drawBars(surface *cairo.Surface, votes map[string]int, minIdx, maxIdx int) {
	// Calculate dimensions
	chartWidth := float64(CanvasWidth - LeftMargin - RightMargin)
	chartHeight := float64(CanvasHeight - TopMargin - BottomMargin)
	numBars := maxIdx - minIdx + 1
	barWidth := (chartWidth - float64(numBars-1)*BarGap) / float64(numBars)

	// Find max votes for scaling
	maxVotes := 0
	for i := minIdx; i <= maxIdx; i++ {
		level := DifficultyLevels[i]
		if v := votes[level]; v > maxVotes {
			maxVotes = v
		}
	}
	if maxVotes == 0 {
		maxVotes = 1
	}

	// Draw shadows first (so they appear behind all bars)
	surface.SetSourceRGBA(0, 0, 0, ShadowAlpha)
	for i := minIdx; i <= maxIdx; i++ {
		level := DifficultyLevels[i]
		voteCount := votes[level]
		if voteCount == 0 {
			continue
		}

		x := float64(LeftMargin) + float64(i-minIdx)*(barWidth+BarGap)
		barHeight := (float64(voteCount) / float64(maxVotes)) * chartHeight
		y := float64(TopMargin) + chartHeight - barHeight

		// Draw shadow (offset down and right)
		drawRoundedTopRect(surface, x+ShadowOffsetX, y+ShadowOffsetY, barWidth, barHeight, BarRadius)
		surface.Fill()
	}

	// Draw each bar
	for i := minIdx; i <= maxIdx; i++ {
		level := DifficultyLevels[i]
		voteCount := votes[level]

		x := float64(LeftMargin) + float64(i-minIdx)*(barWidth+BarGap)
		barHeight := (float64(voteCount) / float64(maxVotes)) * chartHeight
		y := float64(TopMargin) + chartHeight - barHeight

		// Set bar color
		r, g, b := ParseHexColor(DifficultyColors[level])
		surface.SetSourceRGB(float64(r)/255, float64(g)/255, float64(b)/255)

		// Draw rounded rectangle (top corners only)
		drawRoundedTopRect(surface, x, y, barWidth, barHeight, BarRadius)
		surface.Fill()
	}
}

func drawRoundedTopRect(surface *cairo.Surface, x, y, w, h, r float64) {
	if h < r {
		r = h
	}
	if r < 0 {
		r = 0
	}

	surface.MoveTo(x, y+h)                            // bottom-left
	surface.LineTo(x, y+r)                            // left side up to curve
	surface.Arc(x+r, y+r, r, 3.14159, 1.5*3.14159)    // top-left curve
	surface.LineTo(x+w-r, y)                          // top side
	surface.Arc(x+w-r, y+r, r, 1.5*3.14159, 2*3.14159) // top-right curve
	surface.LineTo(x+w, y+h)                          // right side
	surface.ClosePath()
}

func drawXAxisLabels(surface *cairo.Surface, minIdx, maxIdx int) {
	surface.SetSourceRGB(TextColor[0], TextColor[1], TextColor[2])
	surface.SelectFontFace("Inter", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_NORMAL)
	surface.SetFontSize(18)

	chartWidth := float64(CanvasWidth - LeftMargin - RightMargin)
	numBars := maxIdx - minIdx + 1
	barWidth := (chartWidth - float64(numBars-1)*BarGap) / float64(numBars)

	for i := minIdx; i <= maxIdx; i++ {
		level := DifficultyLevels[i]
		x := float64(LeftMargin) + float64(i-minIdx)*(barWidth+BarGap) + barWidth/2

		// Get text extents for centering
		extents := surface.TextExtents(level)
		textX := x - extents.Width/2

		surface.MoveTo(textX, float64(CanvasHeight-BottomMargin+30))
		surface.ShowText(level)
	}
}

func drawYAxisLines(surface *cairo.Surface, votes map[string]int, minIdx, maxIdx int) {
	// Find max votes
	maxVotes := 0
	for i := minIdx; i <= maxIdx; i++ {
		level := DifficultyLevels[i]
		if v := votes[level]; v > maxVotes {
			maxVotes = v
		}
	}
	if maxVotes == 0 {
		maxVotes = 1
	}

	chartWidth := float64(CanvasWidth - LeftMargin - RightMargin)
	chartHeight := float64(CanvasHeight - TopMargin - BottomMargin)

	// Draw faint horizontal grid lines
	surface.SetSourceRGBA(1, 1, 1, 0.15) // White with 15% opacity
	surface.SetLineWidth(1)

	for i := 0; i <= 4; i++ {
		value := (maxVotes * i) / 4
		y := float64(TopMargin) + chartHeight - (float64(value)/float64(maxVotes))*chartHeight

		surface.MoveTo(float64(LeftMargin), y)
		surface.LineTo(float64(LeftMargin)+chartWidth, y)
		surface.Stroke()
	}
}

func drawYAxis(surface *cairo.Surface, votes map[string]int, minIdx, maxIdx int) {
	// Find max votes
	maxVotes := 0
	for i := minIdx; i <= maxIdx; i++ {
		level := DifficultyLevels[i]
		if v := votes[level]; v > maxVotes {
			maxVotes = v
		}
	}
	if maxVotes == 0 {
		maxVotes = 1
	}

	surface.SetSourceRGB(TextColor[0], TextColor[1], TextColor[2])
	surface.SelectFontFace("Inter", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_NORMAL)
	surface.SetFontSize(16)

	chartHeight := float64(CanvasHeight - TopMargin - BottomMargin)

	// Draw 5 tick marks
	for i := 0; i <= 4; i++ {
		value := (maxVotes * i) / 4
		y := float64(TopMargin) + chartHeight - (float64(value)/float64(maxVotes))*chartHeight

		label := formatInt(value)

		extents := surface.TextExtents(label)
		surface.MoveTo(float64(LeftMargin)-extents.Width-10, y+extents.Height/2)
		surface.ShowText(label)
	}
}

func formatInt(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
}

func drawVoteCounts(surface *cairo.Surface, votes map[string]int, minIdx, maxIdx int) {
	surface.SetSourceRGB(TextColor[0], TextColor[1], TextColor[2])
	surface.SelectFontFace("Inter", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_BOLD)
	surface.SetFontSize(18)

	chartWidth := float64(CanvasWidth - LeftMargin - RightMargin)
	chartHeight := float64(CanvasHeight - TopMargin - BottomMargin)
	numBars := maxIdx - minIdx + 1
	barWidth := (chartWidth - float64(numBars-1)*BarGap) / float64(numBars)

	maxVotes := 0
	for i := minIdx; i <= maxIdx; i++ {
		level := DifficultyLevels[i]
		if v := votes[level]; v > maxVotes {
			maxVotes = v
		}
	}
	if maxVotes == 0 {
		maxVotes = 1
	}

	for i := minIdx; i <= maxIdx; i++ {
		level := DifficultyLevels[i]
		voteCount := votes[level]
		if voteCount == 0 {
			continue
		}

		x := float64(LeftMargin) + float64(i-minIdx)*(barWidth+BarGap) + barWidth/2
		barHeight := (float64(voteCount) / float64(maxVotes)) * chartHeight
		y := float64(TopMargin) + chartHeight - barHeight - 15

		label := formatInt(voteCount)
		extents := surface.TextExtents(label)
		surface.MoveTo(x-extents.Width/2, y)
		surface.ShowText(label)
	}
}

func drawAverageLine(surface *cairo.Surface, avg float64, avgLabel string, minIdx, maxIdx int) {
	chartWidth := float64(CanvasWidth - LeftMargin - RightMargin)
	chartHeight := float64(CanvasHeight - TopMargin - BottomMargin)

	// Calculate x position based on average value
	minValue := DifficultyRanges[DifficultyLevels[minIdx]].Lower
	maxValue := DifficultyRanges[DifficultyLevels[maxIdx]].Upper

	xRatio := (avg - minValue) / (maxValue - minValue)
	x := float64(LeftMargin) + xRatio*chartWidth

	// Get the difficulty color for the average
	r, g, b := ParseHexColor(DifficultyColors[avgLabel])
	rf, gf, bf := float64(r)/255, float64(g)/255, float64(b)/255

	dashes := []float64{8, 5}

	// Draw white outline/stroke first (thicker, behind)
	surface.SetSourceRGB(1, 1, 1)
	surface.SetLineWidth(6)
	surface.SetDash(dashes, len(dashes), 0)
	surface.MoveTo(x, float64(TopMargin))
	surface.LineTo(x, float64(TopMargin)+chartHeight)
	surface.Stroke()

	// Draw colored line on top
	surface.SetSourceRGB(rf, gf, bf)
	surface.SetLineWidth(3)
	surface.SetDash(dashes, len(dashes), 0)
	surface.MoveTo(x, float64(TopMargin))
	surface.LineTo(x, float64(TopMargin)+chartHeight)
	surface.Stroke()

	// Draw label in difficulty color (no need to reset dash for text)
	surface.SelectFontFace("Inter", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_BOLD)
	surface.SetFontSize(18)

	labelText := "Avg: " + formatFloat(avg) + " (" + avgLabel + ")"
	extents := surface.TextExtents(labelText)

	labelX := x - extents.Width/2
	labelY := float64(TopMargin) - 45

	// Keep label in bounds
	if labelX < float64(LeftMargin) {
		labelX = float64(LeftMargin)
	}
	if labelX+extents.Width > float64(CanvasWidth-RightMargin) {
		labelX = float64(CanvasWidth-RightMargin) - extents.Width
	}

	surface.MoveTo(labelX, labelY)
	surface.ShowText(labelText)
}

func formatFloat(f float64) string {
	// Simple formatting to 2 decimal places
	intPart := int(f)
	decPart := int((f - float64(intPart)) * 100)
	if decPart < 0 {
		decPart = -decPart
	}
	result := formatInt(intPart) + "."
	if decPart < 10 {
		result += "0"
	}
	result += formatInt(decPart)
	return result
}
