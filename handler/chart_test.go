// handler/chart_test.go
package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateRequest(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantErr    bool
		errContain string
	}{
		{
			name:    "valid request",
			body:    `{"votes":{"Easy":5,"Medium":10}}`,
			wantErr: false,
		},
		{
			name:       "invalid json",
			body:       `{not json}`,
			wantErr:    true,
			errContain: "invalid JSON",
		},
		{
			name:       "missing votes",
			body:       `{}`,
			wantErr:    true,
			errContain: "missing votes",
		},
		{
			name:       "empty votes",
			body:       `{"votes":{}}`,
			wantErr:    true,
			errContain: "no votes",
		},
		{
			name:       "invalid difficulty",
			body:       `{"votes":{"NotReal":5}}`,
			wantErr:    true,
			errContain: "invalid difficulty",
		},
		{
			name:       "negative votes",
			body:       `{"votes":{"Easy":-1}}`,
			wantErr:    true,
			errContain: "invalid vote count",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/chart", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")

			votes, err := ParseAndValidate(req)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
					return
				}
				if tt.errContain != "" && !contains(err.Error(), tt.errContain) {
					t.Errorf("error %q should contain %q", err.Error(), tt.errContain)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if len(votes) == 0 {
					t.Error("expected votes, got empty map")
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}

func TestChartHandler(t *testing.T) {
	body := `{"votes":{"Medium":10,"Medium +":5}}`
	req := httptest.NewRequest(http.MethodPost, "/chart", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	ChartHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status: got %d want %d", rr.Code, http.StatusOK)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "image/webp" {
		t.Errorf("wrong content type: got %s want image/webp", contentType)
	}

	if rr.Body.Len() == 0 {
		t.Error("expected non-empty response body")
	}
}

func TestChartHandlerError(t *testing.T) {
	body := `{"votes":{}}`
	req := httptest.NewRequest(http.MethodPost, "/chart", bytes.NewBufferString(body))

	rr := httptest.NewRecorder()
	ChartHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("handler returned wrong status: got %d want %d", rr.Code, http.StatusBadRequest)
	}

	var errResp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}
	if _, ok := errResp["error"]; !ok {
		t.Error("expected error field in response")
	}
}
