// handler/chart.go
package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/genjishimada/playtest-plotter/chart"
)

// ChartRequest represents the incoming request body
type ChartRequest struct {
	Votes map[string]int `json:"votes"`
}

// ParseAndValidate parses and validates the chart request
func ParseAndValidate(r *http.Request) (map[string]int, error) {
	var req ChartRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.New("invalid JSON")
	}

	if req.Votes == nil {
		return nil, errors.New("missing votes field")
	}

	totalVotes := 0
	for level, count := range req.Votes {
		if _, ok := chart.DifficultyIndex(level); !ok {
			return nil, fmt.Errorf("invalid difficulty: %s", level)
		}
		if count < 0 {
			return nil, fmt.Errorf("invalid vote count for %s", level)
		}
		totalVotes += count
	}

	if totalVotes == 0 {
		return nil, errors.New("no votes provided")
	}

	return req.Votes, nil
}

// ChartHandler handles POST /chart requests
func ChartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	votes, err := ParseAndValidate(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	imgData, err := chart.RenderChart(votes)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate chart")
		return
	}

	w.Header().Set("Content-Type", "image/webp")
	w.WriteHeader(http.StatusOK)
	w.Write(imgData)
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
