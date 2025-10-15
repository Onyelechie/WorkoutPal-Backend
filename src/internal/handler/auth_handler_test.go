package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"

	"workoutpal/src/internal/model"
	mock_service "workoutpal/src/mock_internal/domain/service"
)

func newAuthHandlerMocks(t *testing.T) (*mock_service.MockUserService, *mock_service.MockAuthService, *authHandler, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	userSvc := mock_service.NewMockUserService(ctrl)
	authSvc := mock_service.NewMockAuthService(ctrl)

	h := &authHandler{
		userService: userSvc,
		authService: authSvc,
	}
	return userSvc, authSvc, h, ctrl
}

func mustJSONBody(t *testing.T, v any) *bytes.Buffer {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return bytes.NewBuffer(b)
}

func TestAuthHandler_Login_OK(t *testing.T) {
	_, authSvc, h, _ := newAuthHandlerMocks(t)

	reqBody := model.LoginRequest{Email: "A@B.com", Password: "pw"}
	authSvc.
		EXPECT().
		Authenticate(gomock.Any(), gomock.AssignableToTypeOf(model.LoginRequest{})).
		DoAndReturn(func(_ interface{}, got model.LoginRequest) (*model.User, error) {
			if got.Email != strings.ToLower(reqBody.Email) {
				t.Fatalf("expected lowercased email, got %q", got.Email)
			}
			return &model.User{ID: 123, Email: got.Email, Name: "Tester"}, nil
		})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/login", mustJSONBody(t, reqBody))
	r.Header.Set("Content-Type", "application/json")

	h.Login(w, r)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	cookies := resp.Cookies()
	var tokenCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "access_token" {
			tokenCookie = c
			break
		}
	}
	if tokenCookie == nil || tokenCookie.Value == "" {
		t.Fatalf("expected access_token cookie to be set")
	}
	if !tokenCookie.HttpOnly {
		t.Fatalf("expected HttpOnly cookie")
	}

	var gotUser model.User
	if err := json.NewDecoder(resp.Body).Decode(&gotUser); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if gotUser.ID != 123 || gotUser.Email != strings.ToLower(reqBody.Email) {
		t.Fatalf("unexpected user in response: %+v", gotUser)
	}
}

func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
	_, _, h, _ := newAuthHandlerMocks(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("{"))
	r.Header.Set("Content-Type", "application/json")

	h.Login(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestAuthHandler_Login_AuthFail(t *testing.T) {
	_, authSvc, h, _ := newAuthHandlerMocks(t)

	authSvc.
		EXPECT().
		Authenticate(gomock.Any(), gomock.AssignableToTypeOf(model.LoginRequest{})).
		Return(nil, errors.New("invalid credentials"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/login", mustJSONBody(t, model.LoginRequest{Email: "a@b.com", Password: "bad"}))
	r.Header.Set("Content-Type", "application/json")

	h.Login(w, r)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", w.Code)
	}
}

func TestAuthHandler_Me_Unauthorized_NoClaims(t *testing.T) {
	userSvc, _, h, _ := newAuthHandlerMocks(t)

	userSvc.EXPECT().ReadUserByEmail(gomock.Any()).Times(0)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/me", nil)

	h.Me(w, r)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", w.Code)
	}
}

func TestAuthHandler_Logout_ClearsCookie(t *testing.T) {
	_, _, h, _ := newAuthHandlerMocks(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/logout", nil)

	h.Logout(w, r)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	var cleared *http.Cookie
	for _, c := range resp.Cookies() {
		if c.Name == "access_token" {
			cleared = c
			break
		}
	}
	if cleared == nil {
		t.Fatalf("expected access_token cookie to be set (for clearing)")
	}
	if cleared.MaxAge != -1 {
		t.Fatalf("expected MaxAge -1 (cleared), got %d", cleared.MaxAge)
	}

	var br struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&br); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if br.Message != "success" {
		t.Fatalf("message = %q, want %q", br.Message, "success")
	}
}

var _ = jwt.MapClaims{}
