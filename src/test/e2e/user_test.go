package e2e

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"testing"
	"time"
)

type createdUser struct {
	ID           int64   `json:"id"`
	Username     string  `json:"username"`
	Email        string  `json:"email"`
	Name         string  `json:"name"`
	Age          int     `json:"age"`
	Height       float64 `json:"height"`
	HeightMetric string  `json:"heightMetric"`
	Weight       float64 `json:"weight"`
	WeightMetric string  `json:"weightMetric"`
}

// --- Helpers ---

func randString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	res := make([]byte, n)
	for i := range res {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		res[i] = letters[num.Int64()]
	}
	return string(res)
}

func uniqueUsername() string {
	return "e2e_" + time.Now().Format("20060102_150405") + "_" + randString(4)
}

func uniquePassword() string {
	return "pw_" + randString(16)
}

func uniqueEmail() string {
	return fmt.Sprintf("e2e+%s_%s@example.com", time.Now().Format("20060102_150405"), randString(4))
}

func newUserPayload() map[string]any {
	return map[string]any{
		"username":     uniqueUsername(),
		"password":     uniquePassword(),
		"name":         "E2E User",
		"email":        uniqueEmail(),
		"avatar":       "",
		"age":          25,
		"height":       175.0,
		"heightMetric": "cm",
		"weight":       70.0,
		"weightMetric": "kg",
	}
}

func createUser(t *testing.T) createdUser {
	t.Helper()
	body := newUserPayload()
	resp := doRequest(t, http.MethodPost, "/users", body, nil)
	mustStatus(t, resp, http.StatusCreated)
	defer resp.Body.Close()

	u := mustDecode[createdUser](t, resp)
	if u.ID == 0 || u.Username == "" || u.Email == "" {
		t.Fatalf("unexpected created user: %#v", u)
	}
	return u
}

func int64ToStr(n int64) string {
	return strconv.FormatInt(n, 10)
}

func testEndToEnd_Users_Create(t *testing.T) {
	body := newUserPayload()
	resp := doRequest(t, http.MethodPost, "/users", body, nil)
	mustStatus(t, resp, http.StatusCreated)

	var u createdUser
	u = mustDecode[createdUser](t, resp)
	_ = resp.Body.Close()

	if u.ID == 0 {
		t.Fatalf("expected non-zero created user id")
	}
	if u.Email == "" || u.Username == "" {
		t.Fatalf("expected username/email to be set, got username=%q email=%q", u.Username, u.Email)
	}
}

func testEndToEnd_Users_List(t *testing.T) {
	_ = createUser(t)

	resp := doRequest(t, http.MethodGet, "/users", nil, nil)
	mustStatus(t, resp, http.StatusOK)

	var users []map[string]any
	users = mustDecode[[]map[string]any](t, resp)
	_ = resp.Body.Close()

	if len(users) == 0 {
		t.Fatalf("expected at least one user")
	}
}

func testEndToEnd_Users_GetByID(t *testing.T) {
	created := createUser(t)

	resp := doRequest(t, http.MethodGet, "/users/"+int64ToStr(created.ID), nil, nil)
	mustStatus(t, resp, http.StatusOK)

	type getResp struct {
		ID       int64  `json:"id"`
		Email    string `json:"email"`
		Username string `json:"username"`
		Name     string `json:"name"`
	}
	u := mustDecode[getResp](t, resp)
	_ = resp.Body.Close()

	if u.ID != created.ID {
		t.Fatalf("expected id=%d, got=%d", created.ID, u.ID)
	}
	if u.Email == "" {
		t.Fatalf("expected non-empty email")
	}
}

func testEndToEnd_Users_UpdateByID(t *testing.T) {
	created := createUser(t)

	newUsername := "e2e_updated_" + randString(6)
	newEmail := "e2e.updated+" + time.Now().Format("20060102_150405") + "_" + randString(4) + "@example.com"

	body := map[string]any{
		"username":     newUsername,
		"password":     uniquePassword(),
		"name":         "Updated E2E",
		"email":        newEmail,
		"avatar":       "",
		"age":          26,
		"height":       176.5,
		"heightMetric": "cm",
		"weight":       72.0,
		"weightMetric": "kg",
	}

	resp := doRequest(t, http.MethodPatch, "/users/"+int64ToStr(created.ID), body, nil)
	mustStatus(t, resp, http.StatusOK)

	type updateResp struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Name     string `json:"name"`
		Age      int    `json:"age"`
	}
	u := mustDecode[updateResp](t, resp)
	_ = resp.Body.Close()

	if u.ID != created.ID {
		t.Fatalf("expected id=%d after update, got=%d", created.ID, u.ID)
	}
	if u.Username != newUsername {
		t.Fatalf("username not updated, got=%q want=%q", u.Username, newUsername)
	}
}

func testEndToEnd_Users_Create_Invalid(t *testing.T) {
	body := map[string]any{
		"username":     "",
		"password":     "short",
		"name":         "",
		"email":        "not-an-email",
		"avatar":       "",
		"age":          -2,
		"height":       -20.0,
		"heightMetric": "meters-but-wrong?",
		"weight":       -5.0,
		"weightMetric": "kg",
	}
	resp := doRequest(t, http.MethodPost, "/users", body, nil)
	if resp.StatusCode != http.StatusInternalServerError {
		var m map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&m)
		_ = resp.Body.Close()
		t.Fatalf("expected 500, got=%d body=%v", resp.StatusCode, m)
	}
	_ = resp.Body.Close()
}

func testEndToEnd_Users_Create_Duplicate(t *testing.T) {
	body := newUserPayload()

	resp1 := doRequest(t, http.MethodPost, "/users", body, nil)
	if resp1.StatusCode != http.StatusCreated && resp1.StatusCode != http.StatusOK {
		_ = resp1.Body.Close()
		t.Fatalf("initial create unexpected status=%d", resp1.StatusCode)
	}
	_ = resp1.Body.Close()

	resp2 := doRequest(t, http.MethodPost, "/users", body, nil)
	if resp2.StatusCode != http.StatusInternalServerError {
		var m map[string]any
		_ = json.NewDecoder(resp2.Body).Decode(&m)
		_ = resp2.Body.Close()
		t.Fatalf("expected 500 on duplicate, got=%d body=%v", resp2.StatusCode, m)
	}
	_ = resp2.Body.Close()
}

func testEndToEnd_Users_Delete_Success(t *testing.T) {
	created := createUser(t)

	delResp := doRequest(t, http.MethodDelete, "/users/"+int64ToStr(created.ID), nil, nil)
	mustStatus(t, delResp, http.StatusOK)
	_ = delResp.Body.Close()

	getResp := doRequest(t, http.MethodGet, "/users/"+int64ToStr(created.ID), nil, nil)
	if getResp.StatusCode != http.StatusInternalServerError {
		var m map[string]any
		_ = json.NewDecoder(getResp.Body).Decode(&m)
		_ = getResp.Body.Close()
		t.Fatalf("expected 500 after delete, got=%d body=%v", getResp.StatusCode, m)
	}
	_ = getResp.Body.Close()
}

func testEndToEnd_Users_GetByID_NotFound(t *testing.T) {
	resp := doRequest(t, http.MethodGet, "/users/99999999", nil, nil)
	if resp.StatusCode != http.StatusInternalServerError {
		var m map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&m)
		_ = resp.Body.Close()
		t.Fatalf("expected 500, got=%d body=%v", resp.StatusCode, m)
	}
	_ = resp.Body.Close()
}
