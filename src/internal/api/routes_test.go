package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"workoutpal/src/internal/config"
	"workoutpal/src/internal/dependency"
	"workoutpal/src/internal/model"
	mock_repository "workoutpal/src/mock_internal/domain/repository"
	mock_service "workoutpal/src/mock_internal/domain/service"
)

/* ----------------------------- helpers ----------------------------- */

func do(ts *httptest.Server, method, path string, body io.Reader, headers map[string]string) *http.Response {
	req, _ := http.NewRequest(method, ts.URL+path, body)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, _ := http.DefaultClient.Do(req)
	return resp
}

func chiRouterWithGlobalMiddleware() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	return r
}

/* ------------------------------- tests ------------------------------ */

// Basic sanity: the pure RegisterRoutes (which wraps with CORS & swagger) builds and /health works.
func TestRegisterRoutes_Health_OK(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}

	h := RegisterRoutes(cfg, &sql.DB{})
	ts := httptest.NewServer(h)
	defer ts.Close()

	resp := do(ts, http.MethodGet, "/health", nil, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /health status = %d, want 200", resp.StatusCode)
	}
}

func TestRegisterRoutes_CORS_AllowsAllowedOrigin_OnSimpleGET(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{JWTSecret: "test-secret"}
	h := RegisterRoutes(cfg, &sql.DB{})
	ts := httptest.NewServer(h)
	defer ts.Close()

	headers := map[string]string{"Origin": "http://localhost:4200"}
	resp := do(ts, http.MethodGet, "/health", nil, headers)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /health status = %d, want 200", resp.StatusCode)
	}

	if got := resp.Header.Get("Access-Control-Allow-Origin"); got != "http://localhost:4200" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want %q", got, "http://localhost:4200")
	}
}

// Unit-test the DI seam: /login must call AuthService.Authenticate and succeed.
func TestRoutes_Login_CallsAuthService_AndReturns200(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	// Build mocks for the services the handlers need.
	mockAuth := mock_service.NewMockAuthService(ctrl)
	mockUserSvc := mock_service.NewMockUserService(ctrl)
	// Other services are needed only to construct handlers; no expectations required.
	mockGoalSvc := mock_service.NewMockGoalService(ctrl)
	mockRelSvc := mock_service.NewMockRelationshipService(ctrl)
	mockRoutineSvc := mock_service.NewMockRoutineService(ctrl)
	mockExerciseSvc := mock_service.NewMockExerciseService(ctrl)

	// Some handlers take repositories directly; provide a mock user repo as required by your AuthHandler ctor.
	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)

	deps := dependency.AppDependencies{
		UserRepository:      mockUserRepo,
		UserService:         mockUserSvc,
		GoalService:         mockGoalSvc,
		RelationshipService: mockRelSvc,
		RoutineService:      mockRoutineSvc,
		ExerciseService:     mockExerciseSvc,
		AuthService:         mockAuth,
	}

	r := chiRouterWithGlobalMiddleware()
	h := Routes(r, deps, []byte("test-secret"))
	ts := httptest.NewServer(h)
	defer ts.Close()

	reqBody := model.LoginRequest{Email: "a@b.com", Password: "pw"}
	mockAuth.
		EXPECT().
		Authenticate(gomock.Any(), gomock.AssignableToTypeOf(model.LoginRequest{})).
		Return(&model.User{ID: 1, Email: reqBody.Email}, nil)

	buf, _ := json.Marshal(reqBody)
	resp := do(ts, http.MethodPost, "/login", bytes.NewBuffer(buf), map[string]string{
		"Content-Type": "application/json",
	})
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		t.Fatal("POST /login returned 404; route not registered")
	}
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("POST /login status = %d, want 200; body=%s", resp.StatusCode, string(b))
	}
}

// Protected endpoints without Authorization header should be rejected by AuthMiddleware.
func TestRoutes_ProtectedEndpoints_UnauthorizedWithoutToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	deps := dependency.AppDependencies{
		UserRepository:      mock_repository.NewMockUserRepository(ctrl),
		UserService:         mock_service.NewMockUserService(ctrl),
		GoalService:         mock_service.NewMockGoalService(ctrl),
		RelationshipService: mock_service.NewMockRelationshipService(ctrl),
		RoutineService:      mock_service.NewMockRoutineService(ctrl),
		ExerciseService:     mock_service.NewMockExerciseService(ctrl),
		AuthService:         mock_service.NewMockAuthService(ctrl),
	}

	r := chiRouterWithGlobalMiddleware()
	h := Routes(r, deps, []byte("test-secret"))
	ts := httptest.NewServer(h)
	defer ts.Close()

	// /me requires auth
	resp := do(ts, http.MethodGet, "/me", nil, nil)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("GET /me status = %d, want 401", resp.StatusCode)
	}

	// /users GET requires auth
	resp = do(ts, http.MethodGet, "/users/", nil, nil)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("GET /users status = %d, want 401", resp.StatusCode)
	}
}

// Unprotected create route exists; with invalid body it should return a 4xx (but not 404), proving wiring.
func TestRoutes_CreateUser_RouteExists_Returns4xxOnBadBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	deps := dependency.AppDependencies{
		UserRepository:      mock_repository.NewMockUserRepository(ctrl),
		UserService:         mock_service.NewMockUserService(ctrl),
		GoalService:         mock_service.NewMockGoalService(ctrl),
		RelationshipService: mock_service.NewMockRelationshipService(ctrl),
		RoutineService:      mock_service.NewMockRoutineService(ctrl),
		ExerciseService:     mock_service.NewMockExerciseService(ctrl),
		AuthService:         mock_service.NewMockAuthService(ctrl),
	}

	r := chiRouterWithGlobalMiddleware()
	h := Routes(r, deps, []byte("test-secret"))
	ts := httptest.NewServer(h)
	defer ts.Close()

	// Bad/empty body to force handler to bail early without touching service/DB.
	resp := do(ts, http.MethodPost, "/users/", bytes.NewBufferString(`{}`), map[string]string{
		"Content-Type": "application/json",
	})

	if resp.StatusCode == http.StatusNotFound {
		t.Fatal("POST /users returned 404; route not registered")
	}
	if resp.StatusCode < 400 || resp.StatusCode > 499 {
		t.Fatalf("POST /users status = %d, want 4xx (invalid body)", resp.StatusCode)
	}
}
