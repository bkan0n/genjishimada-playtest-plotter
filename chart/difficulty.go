package chart

import "strconv"

var DifficultyLevels = []string{
	"Easy -", "Easy", "Easy +",
	"Medium -", "Medium", "Medium +",
	"Hard -", "Hard", "Hard +",
	"Very Hard -", "Very Hard", "Very Hard +",
	"Extreme -", "Extreme", "Extreme +",
	"Hell",
}

var difficultyIndexMap = func() map[string]int {
	m := make(map[string]int, len(DifficultyLevels))
	for i, level := range DifficultyLevels {
		m[level] = i
	}
	return m
}()

func DifficultyIndex(name string) (int, bool) {
	idx, ok := difficultyIndexMap[name]
	return idx, ok
}

var DifficultyColors = map[string]string{
	"Easy -":      "#66ff66",
	"Easy":        "#4dcc4d",
	"Easy +":      "#33cc33",
	"Medium -":    "#99ff33",
	"Medium":      "#99e600",
	"Medium +":    "#80cc00",
	"Hard -":      "#ffd633",
	"Hard":        "#ffb300",
	"Hard +":      "#ff9900",
	"Very Hard -": "#ff8000",
	"Very Hard":   "#e67e00",
	"Very Hard +": "#cc6600",
	"Extreme -":   "#ff4d00",
	"Extreme":     "#e04300",
	"Extreme +":   "#b92d00",
	"Hell":        "#990000",
}

// DifficultyMidpoints for weighted average calculation
var DifficultyMidpoints = map[string]float64{
	"Easy -":      0.89,
	"Easy":        1.47,
	"Easy +":      2.06,
	"Medium -":    2.65,
	"Medium":      3.23,
	"Medium +":    3.83,
	"Hard -":      4.42,
	"Hard":        5.0,
	"Hard +":      5.58,
	"Very Hard -": 6.17,
	"Very Hard":   6.76,
	"Very Hard +": 7.36,
	"Extreme -":   7.95,
	"Extreme":     8.53,
	"Extreme +":   9.12,
	"Hell":        9.71,
}

// DifficultyRange defines the range for a difficulty level
type DifficultyRange struct {
	Lower          float64
	Upper          float64
	UpperInclusive bool
}

var DifficultyRanges = map[string]DifficultyRange{
	"Easy -":      {0.0, 1.18, false},
	"Easy":        {1.18, 1.76, false},
	"Easy +":      {1.76, 2.35, false},
	"Medium -":    {2.35, 2.94, false},
	"Medium":      {2.94, 3.53, false},
	"Medium +":    {3.53, 4.12, false},
	"Hard -":      {4.12, 4.71, false},
	"Hard":        {4.71, 5.29, false},
	"Hard +":      {5.29, 5.88, false},
	"Very Hard -": {5.88, 6.47, false},
	"Very Hard":   {6.47, 7.06, false},
	"Very Hard +": {7.06, 7.65, false},
	"Extreme -":   {7.65, 8.24, false},
	"Extreme":     {8.24, 8.82, false},
	"Extreme +":   {8.82, 9.41, false},
	"Hell":        {9.41, 10.0, true},
}

func ParseHexColor(hex string) (r, g, b uint8) {
	if len(hex) != 7 || hex[0] != '#' {
		return 0, 0, 0
	}
	rVal, _ := strconv.ParseUint(hex[1:3], 16, 8)
	gVal, _ := strconv.ParseUint(hex[3:5], 16, 8)
	bVal, _ := strconv.ParseUint(hex[5:7], 16, 8)
	return uint8(rVal), uint8(gVal), uint8(bVal)
}
