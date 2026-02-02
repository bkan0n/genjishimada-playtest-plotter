// chart/render.go
package chart

import (
	"bytes"
	"image"
	"strings"

	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/ungerik/go-cairo"
)

// Cached encoder options - created once at startup
var webpEncoderOptions *encoder.Options

func init() {
	opts, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 85)
	if err != nil {
		panic("failed to create webp encoder options: " + err.Error())
	}
	webpEncoderOptions = opts
}

const (
	CanvasWidth  = 1000
	CanvasHeight = 500
	LeftMargin   = 60
	RightMargin  = 45
	TopMargin    = 75
	BottomMargin = 70
	BarGap       = 9
	BarRadius    = 20
)

var (
	BackgroundColor   = [3]float64{0.168, 0.176, 0.192} // #2b2d31
	TextColor         = [3]float64{1.0, 1.0, 1.0}       // white
	TextShadowColor   = [4]float64{0, 0, 0, 0.5}        // black with 50% opacity
	TextShadowOffsetX = 1.5
	TextShadowOffsetY = 1.5
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
	maxVotes := calculateMaxVotes(votes, minIdx, maxIdx)

	// Draw chart elements
	drawYAxisLines(surface, maxVotes, minIdx, maxIdx) // Draw grid lines first (behind bars)
	drawBars(surface, votes, minIdx, maxIdx, maxVotes)
	drawXAxisLabels(surface, minIdx, maxIdx)
	drawYAxis(surface, maxVotes)
	drawVoteCounts(surface, votes, minIdx, maxIdx, maxVotes)
	drawAverageLine(surface, avg, avgLabel, minIdx, maxIdx)

	// Convert to image.RGBA
	img := surfaceToImage(surface)

	// Encode to WebP with pre-allocated buffer
	buf := bytes.NewBuffer(make([]byte, 0, 50*1024))
	if err := webp.Encode(buf, img, webpEncoderOptions); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// calculateMaxVotes finds the maximum vote count in the visible window
func calculateMaxVotes(votes map[string]int, minIdx, maxIdx int) int {
	maxVotes := 0
	for i := minIdx; i <= maxIdx; i++ {
		level := DifficultyLevels[i]
		if v := votes[level]; v > maxVotes {
			maxVotes = v
		}
	}
	if maxVotes == 0 {
		return 1
	}
	return maxVotes
}

func surfaceToImage(surface *cairo.Surface) *image.RGBA {
	width := surface.GetWidth()
	height := surface.GetHeight()
	data := surface.GetData()

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Optimized single-loop pixel conversion
	// Cairo ARGB32 (little endian): B, G, R, A → Go RGBA: R, G, B, A
	pixelCount := width * height
	for i := 0; i < pixelCount; i++ {
		off := i * 4
		img.Pix[off] = data[off+2]   // R
		img.Pix[off+1] = data[off+1] // G
		img.Pix[off+2] = data[off]   // B
		img.Pix[off+3] = data[off+3] // A
	}
	return img
}

// drawTextWithShadow draws text with a drop shadow effect
func drawTextWithShadow(surface *cairo.Surface, text string, x, y float64) {
	// Draw shadow first
	surface.SetSourceRGBA(TextShadowColor[0], TextShadowColor[1], TextShadowColor[2], TextShadowColor[3])
	surface.MoveTo(x+TextShadowOffsetX, y+TextShadowOffsetY)
	surface.ShowText(text)

	// Draw main text
	surface.SetSourceRGB(TextColor[0], TextColor[1], TextColor[2])
	surface.MoveTo(x, y)
	surface.ShowText(text)
}

const (
	ShadowOffsetX = 3
	ShadowOffsetY = 3
	ShadowAlpha   = 0.3
)

func drawBars(surface *cairo.Surface, votes map[string]int, minIdx, maxIdx, maxVotes int) {
	// Calculate dimensions
	chartWidth := float64(CanvasWidth - LeftMargin - RightMargin)
	chartHeight := float64(CanvasHeight - TopMargin - BottomMargin)
	numBars := maxIdx - minIdx + 1
	barWidth := (chartWidth - float64(numBars-1)*BarGap) / float64(numBars)

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
	surface.SelectFontFace("Bank Sans EF CY", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_NORMAL)
	surface.SetFontSize(16)

	chartWidth := float64(CanvasWidth - LeftMargin - RightMargin)
	numBars := maxIdx - minIdx + 1
	barWidth := (chartWidth - float64(numBars-1)*BarGap) / float64(numBars)

	for i := minIdx; i <= maxIdx; i++ {
		level := DifficultyLevels[i]
		x := float64(LeftMargin) + float64(i-minIdx)*(barWidth+BarGap) + barWidth/2

		// Get text extents for centering
		upperLevel := strings.ToUpper(level)
		extents := surface.TextExtents(upperLevel)
		textX := x - extents.Width/2

		drawTextWithShadow(surface, upperLevel, textX, float64(CanvasHeight-BottomMargin+30))
	}
}

func drawYAxisLines(surface *cairo.Surface, maxVotes, minIdx, maxIdx int) {
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

func drawYAxis(surface *cairo.Surface, maxVotes int) {
	surface.SelectFontFace("Bank Sans EF CY", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_NORMAL)
	surface.SetFontSize(12)

	chartHeight := float64(CanvasHeight - TopMargin - BottomMargin)

	// Draw 5 tick marks
	for i := 0; i <= 4; i++ {
		value := (maxVotes * i) / 4
		y := float64(TopMargin) + chartHeight - (float64(value)/float64(maxVotes))*chartHeight

		label := formatInt(value)

		extents := surface.TextExtents(label)
		drawTextWithShadow(surface, label, float64(LeftMargin)-extents.Width-10, y+extents.Height/2)
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

func drawVoteCounts(surface *cairo.Surface, votes map[string]int, minIdx, maxIdx, maxVotes int) {
	surface.SelectFontFace("Bank Sans EF CY", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_BOLD)
	surface.SetFontSize(13)

	chartWidth := float64(CanvasWidth - LeftMargin - RightMargin)
	chartHeight := float64(CanvasHeight - TopMargin - BottomMargin)
	numBars := maxIdx - minIdx + 1
	barWidth := (chartWidth - float64(numBars-1)*BarGap) / float64(numBars)

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
		drawTextWithShadow(surface, label, x-extents.Width/2, y)
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

	// Draw white dashed line
	surface.SetSourceRGB(1, 1, 1)
	surface.SetLineWidth(2)
	dashes := []float64{8, 5}
	surface.SetDash(dashes, len(dashes), 0)
	surface.MoveTo(x, float64(TopMargin))
	surface.LineTo(x, float64(TopMargin)+chartHeight)
	surface.Stroke()

	// Draw label with shadow
	surface.SelectFontFace("Bank Sans EF CY", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_BOLD)
	surface.SetFontSize(13)

	labelText := "AVG: " + formatFloat(avg) + " (" + strings.ToUpper(avgLabel) + ")"
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

	drawTextWithShadow(surface, labelText, labelX, labelY)
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
