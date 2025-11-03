package e2e

import (
	"net/http"
	"testing"
)

type achievementCatalog struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	BadgeIcon   string `json:"badgeIcon"`
	Description string `json:"description"`
}

type userAchievement struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"userId"`
	Title       string `json:"title"`
	BadgeIcon   string `json:"badgeIcon"`
	Description string `json:"description"`
	EarnedAt    string `json:"earnedAt"`
}

type createAchievementReq struct {
	UserID        int64 `json:"userId"`
	AchievementID int64 `json:"achievementId"`
}

func testEndToEnd_Achievements_CreateUnlock(t *testing.T) {
	const userID int64 = 1

	catResp := doRequest(t, http.MethodGet, "/achievements", nil, nil)
	defer catResp.Body.Close()
	mustStatus(t, catResp, http.StatusOK)

	catalog := mustDecode[[]achievementCatalog](t, catResp)
	if len(catalog) == 0 {
		t.Fatalf("expected at least one catalog achievement to unlock")
	}
	target := catalog[0]
	if target.ID == 0 || target.Title == "" {
		t.Fatalf("bad catalog row: %#v", target)
	}

	body := createAchievementReq{
		UserID:        userID,
		AchievementID: target.ID,
	}
	resp := doRequest(t, http.MethodPost, "/achievements", body, nil)
	defer resp.Body.Close()
	mustStatus(t, resp, http.StatusCreated)

	created := mustDecode[userAchievement](t, resp)
	if created.ID != target.ID {
		t.Fatalf("expected created.ID == achievementId (%d), got %d", target.ID, created.ID)
	}
	if created.UserID != userID {
		t.Fatalf("user mismatch: want %d got %d", userID, created.UserID)
	}
	if created.Title != target.Title || created.BadgeIcon != target.BadgeIcon {
		t.Fatalf("title/icon mismatch: created=%#v target=%#v", created, target)
	}
	if created.EarnedAt == "" {
		t.Fatalf("expected EarnedAt to be set")
	}
}

func testEndToEnd_Achievements_ListCatalog(t *testing.T) {
	resp := doRequest(t, http.MethodGet, "/achievements", nil, nil)
	defer resp.Body.Close()

	mustStatus(t, resp, http.StatusOK)

	list := mustDecode[[]achievementCatalog](t, resp)
	if len(list) == 0 {
		t.Fatalf("expected at least one catalog achievement")
	}
	if list[0].ID == 0 || list[0].Title == "" {
		t.Fatalf("expected valid id/title, got %#v", list[0])
	}
}

func testEndToEnd_Achievements_ReadUnlockedByUserID(t *testing.T) {
	const userID int64 = 1

	resp := doRequest(t, http.MethodGet, "/achievements/unlocked/"+int64ToStr(userID), nil, nil)
	defer resp.Body.Close()
	mustStatus(t, resp, http.StatusOK)

	list := mustDecode[[]userAchievement](t, resp)
	for _, ua := range list {
		if ua.UserID != userID {
			t.Fatalf("all rows must be for user %d, got %#v", userID, ua)
		}
		if ua.ID == 0 || ua.Title == "" || ua.EarnedAt == "" {
			t.Fatalf("unexpected unlocked row: %#v", ua)
		}
	}
}
