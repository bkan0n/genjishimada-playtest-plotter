// integration_test.go
//go:build integration

package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/genjishimada/playtest-plotter/handler"
)

func TestFullChartGeneration(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/chart", handler.ChartHandler)
	mux.HandleFunc("/health", handler.HealthHandler)

	server := httptest.NewServer(mux)
	defer server.Close()

	// Test health endpoint
	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("health check failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("health returned %d", resp.StatusCode)
	}

	// Test chart generation
	votes := map[string]interface{}{
		"votes": map[string]int{
			"Easy":     2,
			"Easy +":   5,
			"Medium -": 8,
			"Medium":   12,
			"Medium +": 6,
			"Hard -":   3,
			"Hard":     1,
		},
	}
	body, _ := json.Marshal(votes)

	resp, err = http.Post(server.URL+"/chart", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("chart request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("chart returned %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if resp.Header.Get("Content-Type") != "image/webp" {
		t.Errorf("wrong content type: %s", resp.Header.Get("Content-Type"))
	}

	imgData, _ := io.ReadAll(resp.Body)
	if len(imgData) < 100 {
		t.Error("image data too small")
	}

	// Verify WebP header
	if string(imgData[0:4]) != "RIFF" {
		t.Error("not a valid WebP file")
	}
}
