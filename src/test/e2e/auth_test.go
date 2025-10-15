package e2e

import (
	"net/http"
	"testing"
	"workoutpal/src/internal/model"
)

func testEndToEnd_Login(t *testing.T) {
	loginBody := model.LoginRequest{
		Email:    "john@example.com",
		Password: "password123",
	}
	resp := doRequest(t, http.MethodPost, "/login", loginBody, nil)
	mustStatus(t, resp, http.StatusOK)
	_ = resp.Body.Close()
}

func testEndToEnd_Me(t *testing.T) {
	resp := doRequest(t, http.MethodGet, "/me", nil, nil)
	mustStatus(t, resp, http.StatusOK)
	me := mustDecode[model.User](t, resp)
	if me.Email != "john@example.com" {
		t.Fatalf("unexpected /me email. got=%s want=%s", me.Email, "john@example.com")
	}
	_ = resp.Body.Close()
}

func testEndToEnd_Logout(t *testing.T) {
	resp := doRequest(t, http.MethodPost, "/logout", nil, nil)
	mustStatus(t, resp, http.StatusOK)
	_ = resp.Body.Close()
}

func testEndToEnd_MeFail(t *testing.T) {
	resp := doRequest(t, http.MethodGet, "/me", nil, nil)
	if resp.StatusCode == http.StatusOK {
		_ = resp.Body.Close()
		t.Fatalf("/me succeeded after logout; expected 401")
	}
	_ = resp.Body.Close()
}
