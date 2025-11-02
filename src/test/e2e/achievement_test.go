package e2e

import (
	"net/http"
	"testing"
)

type achievement struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"userId"`
	Title       string `json:"title"`
	BadgeIcon   string `json:"badgeIcon"`
	Description string `json:"description"`
	EarnedAt    string `json:"earnedAt"`
}

type createAchievementReq struct {
	UserID      int64  `json:"userId"`
	Title       string `json:"title"`
	BadgeIcon   string `json:"badgeIcon"`
	Description string `json:"description"`
	EarnedAt    string `json:"earnedAt"`
}

func testEndToEnd_Achievements_Create(t *testing.T) {
	body := createAchievementReq{
		UserID:      1,
		Title:       "E2E Achievement " + randStringAlphaNum(6),
		BadgeIcon:   "🏅",
		Description: "auto-generated",
		EarnedAt:    "2025-01-01T00:00:00Z",
	}

	resp := doRequest(t, http.MethodPost, "/achievements", body, nil)
	defer resp.Body.Close()

	mustStatus(t, resp, http.StatusCreated)

	created := mustDecode[achievement](t, resp)
	if created.ID == 0 {
		t.Fatalf("expected non-zero id")
	}
	if created.Title != body.Title {
		t.Fatalf("title mismatch")
	}
	if created.UserID != body.UserID {
		t.Fatalf("user mismatch")
	}
}

func testEndToEnd_Achievements_List(t *testing.T) {
	resp := doRequest(t, http.MethodGet, "/achievements", nil, nil)
	defer resp.Body.Close()

	mustStatus(t, resp, http.StatusOK)

	list := mustDecode[[]achievement](t, resp)
	if len(list) == 0 {
		t.Fatalf("expected at least one achievement")
	}
	if list[0].ID == 0 {
		t.Fatalf("expected id")
	}
	if list[0].Title == "" {
		t.Fatalf("expected title")
	}
}

func testEndToEnd_Achievements_Delete(t *testing.T) {
	body := createAchievementReq{
		UserID:      1,
		Title:       "E2E Delete " + randStringAlphaNum(6),
		BadgeIcon:   "🔥",
		Description: "to delete",
		EarnedAt:    "2025-01-01T00:00:00Z",
	}

	createResp := doRequest(t, http.MethodPost, "/achievements", body, nil)
	defer createResp.Body.Close()
	mustStatus(t, createResp, http.StatusCreated)

	created := mustDecode[achievement](t, createResp)
	delResp := doRequest(t, http.MethodDelete, "/achievements/"+int64ToStr(created.ID), nil, nil)
	defer delResp.Body.Close()
	mustStatus(t, delResp, http.StatusOK)

	msg := mustDecode[basicResponse](t, delResp)
	if msg.Message == "" {
		t.Fatalf("expected success message")
	}
}
