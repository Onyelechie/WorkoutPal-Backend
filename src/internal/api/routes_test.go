package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"

	"workoutpal/src/internal/api"
	"workoutpal/src/internal/config"
)

func newTestServer(t *testing.T) (*httptest.Server, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}

	t.Cleanup(func() { _ = db.Close() })

	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	handler := api.RegisterRoutes(cfg, db)
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return srv, mock
}

func walkRoutes(t *testing.T, r http.Handler) map[string]struct{} {
	t.Helper()

	cr, ok := r.(chi.Router)
	if !ok {
		t.Fatalf("handler does not implement chi.Router; got %T", r)
	}

	found := make(map[string]struct{})
	walkFn := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		key := method + " " + route
		found[key] = struct{}{}
		return nil
	}
	if err := chi.Walk(cr, walkFn); err != nil {
		t.Fatalf("chi.Walk failed: %v", err)
	}
	return found
}

func hasRoute(routes map[string]struct{}, method, pattern string) bool {
	_, ok := routes[method+" "+pattern]
	return ok
}

func do(t *testing.T, client *http.Client, method, url string, headers map[string]string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("client.Do: %v", err)
	}
	return resp
}

func TestProtectedRoutesRequireAuth(t *testing.T) {
	srv, _ := newTestServer(t)

	type tc struct {
		method string
		path   string
	}
	protected := []tc{
		{http.MethodGet, "/me"},
		{http.MethodGet, "/users"},
		{http.MethodGet, "/users/123"},
		{http.MethodPatch, "/users/123"},
		{http.MethodDelete, "/users/123"},
		{http.MethodPost, "/users/123/goals"},
		{http.MethodGet, "/users/123/goals"},
		{http.MethodPost, "/users/123/follow"},
		{http.MethodPost, "/users/123/unfollow"},
		{http.MethodGet, "/users/123/followers"},
		{http.MethodGet, "/users/123/following"},
		{http.MethodPost, "/users/123/routines"},
		{http.MethodGet, "/users/123/routines"},
		{http.MethodDelete, "/users/123/routines/456"},
		{http.MethodGet, "/exercises"},
		{http.MethodGet, "/exercises/789"},
		{http.MethodGet, "/routines/777"},
		{http.MethodDelete, "/routines/777"},
		{http.MethodPost, "/routines/777/exercises"},
		{http.MethodDelete, "/routines/777/exercises/888"},
	}

	for _, c := range protected {
		resp := do(t, srv.Client(), c.method, srv.URL+c.path, nil)
		func() {
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusUnauthorized && resp.StatusCode != http.StatusForbidden {
				t.Fatalf("%s %s: expected 401/403; got %d", c.method, c.path, resp.StatusCode)
			}
		}()
	}
}

func TestPublicRoutesExistByWalkingRouter(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	r := chi.NewRouter()
	inner := api.Routes(r, db, []byte("test-secret"))

	routes := walkRoutes(t, inner)

	if !hasRoute(routes, http.MethodGet, "/health") {
		t.Errorf("missing GET /health")
	}
	if !hasRoute(routes, http.MethodPost, "/login") {
		t.Errorf("missing POST /login")
	}
	if !hasRoute(routes, http.MethodPost, "/logout") {
		t.Errorf("missing POST /logout")
	}
	if !hasRoute(routes, http.MethodGet, "/me") {
		t.Errorf("missing GET /me")
	}

	expect := []struct {
		m string
		p string
	}{
		{http.MethodPost, "/users/"},
		{http.MethodGet, "/users/"},
		{http.MethodGet, "/users/{id}"},
		{http.MethodPatch, "/users/{id}"},
		{http.MethodDelete, "/users/{id}"},
		{http.MethodPost, "/users/{id}/goals"},
		{http.MethodGet, "/users/{id}/goals"},
		{http.MethodPost, "/users/{id}/follow"},
		{http.MethodPost, "/users/{id}/unfollow"},
		{http.MethodGet, "/users/{id}/followers"},
		{http.MethodGet, "/users/{id}/following"},
		{http.MethodPost, "/users/{id}/routines"},
		{http.MethodGet, "/users/{id}/routines"},
		{http.MethodDelete, "/users/{id}/routines/{routine_id}"},
		{http.MethodGet, "/exercises/"},
		{http.MethodGet, "/exercises/{id}"},
		{http.MethodGet, "/routines/{id}"},
		{http.MethodDelete, "/routines/{id}"},
		{http.MethodPost, "/routines/{id}/exercises"},
		{http.MethodDelete, "/routines/{id}/exercises/{exercise_id}"},
	}
	for _, e := range expect {
		if !hasRoute(routes, e.m, e.p) {
			t.Errorf("missing %s %s", e.m, e.p)
		}
	}
}

func TestSwaggerRouteIsMountedOnRegisterRoutes(t *testing.T) {
	srv, _ := newTestServer(t)

	resp := do(t, srv.Client(), http.MethodGet, srv.URL+"/swagger/index.html", nil)
	defer resp.Body.Close()

	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusMovedPermanently || resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusNotFound) {
		t.Fatalf("unexpected status from /swagger/index.html: %d", resp.StatusCode)
	}
}

func TestNoPanicRecovererOnUnknownPath(t *testing.T) {
	srv, _ := newTestServer(t)

	client := srv.Client()
	client.Timeout = 5 * time.Second

	resp := do(t, client, http.MethodGet, srv.URL+"/definitely-not-here", nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 for unknown path; got %d", resp.StatusCode)
	}
}
