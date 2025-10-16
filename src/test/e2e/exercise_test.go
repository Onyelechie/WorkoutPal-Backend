package e2e

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"testing"
	"time"
)

type exercise struct {
	ID                  int64    `json:"id"`
	Name                string   `json:"name"`
	Description         string   `json:"description"`
	Targets             []string `json:"targets"`
	Intensity           string   `json:"intensity"`
	Expertise           string   `json:"expertise"`
	Image               string   `json:"image"`
	Demo                string   `json:"demo"`
	RecommendedCount    int      `json:"recommendedCount"`
	RecommendedSets     int      `json:"recommendedSets"`
	RecommendedDuration int      `json:"recommendedDuration"`
	Custom              bool     `json:"custom"`
}

func randStringAlphaNum(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	out := make([]byte, n)
	for i := range out {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		out[i] = letters[num.Int64()]
	}
	return string(out)
}

func uniqueExerciseName() string {
	return "E2E_Exercise_" + time.Now().Format("20060102_150405") + "_" + randStringAlphaNum(4)
}

func newExercisePayload() map[string]any {
	return map[string]any{
		"name":                uniqueExerciseName(),
		"description":         "Auto-generated e2e exercise",
		"targets":             []string{"chest", "triceps"},
		"intensity":           "moderate",
		"expertise":           "beginner",
		"image":               "",
		"demo":                "",
		"recommendedCount":    10,
		"recommendedSets":     3,
		"recommendedDuration": 45,
	}
}

func testEndToEnd_Exercises_List_WithQueryParams(t *testing.T) {
	q := "?target=chest&intensity=moderate&expertise=beginner"
	resp := doRequest(t, http.MethodGet, "/exercises"+q, nil, nil)
	mustStatus(t, resp, http.StatusOK)
	defer resp.Body.Close()

	// Ensure it returns a JSON array
	var list []map[string]any
	list = mustDecode[[]map[string]any](t, resp)
	if list == nil {
		t.Fatalf("expected a list response, got nil")
	}
}

func testEndToEnd_Exercises_GetByID(t *testing.T) {

	resp := doRequest(t, http.MethodGet, "/exercises/"+int64ToStr(1), nil, nil)
	mustStatus(t, resp, http.StatusOK)
	defer resp.Body.Close()

	got := mustDecode[exercise](t, resp)
	if got.ID != 1 {
		t.Fatalf("expected id=%d, got=%d", 1, got.ID)
	}
	if got.Name == "" {
		t.Fatalf("expected non-empty name")
	}
}

func testEndToEnd_Exercises_GetByID_NotFound(t *testing.T) {
	resp := doRequest(t, http.MethodGet, "/exercises/99999999", nil, nil)
	if resp.StatusCode != http.StatusInternalServerError {
		var m map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&m)
		_ = resp.Body.Close()
		t.Fatalf("expected 500, got=%d body=%v", resp.StatusCode, m)
	}
	_ = resp.Body.Close()
}
