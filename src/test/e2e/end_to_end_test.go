package e2e

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"testing"
	"time"
	"workoutpal/src/internal/api"
	"workoutpal/src/internal/config"
)

var (
	baseURL     = envOr("E2E_BASE_URL", "http://localhost:"+envOr("PORT", "8080"))
	loginEmail  = envOr("E2E_LOGIN_EMAIL", "john@example.com")
	loginPass   = envOr("E2E_LOGIN_PASSWORD", "password123")
	client      *http.Client
	startedOnce bool
)

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func doRequest(t *testing.T, method, path string, body any, headers map[string]string) *http.Response {
	t.Helper()

	var data *bytes.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		data = bytes.NewReader(jsonBytes)
	} else {
		data = bytes.NewReader(nil)
	}

	req, err := http.NewRequest(method, baseURL+path, data)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	return resp
}

func mustStatus(t *testing.T, resp *http.Response, want int) {
	t.Helper()
	if resp.StatusCode != want {
		var m map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&m)
		_ = resp.Body.Close()
		t.Fatalf("status %d (want %d). body=%v", resp.StatusCode, want, m)
	}
}

func mustDecode[T any](t *testing.T, resp *http.Response) T {
	t.Helper()
	defer resp.Body.Close()
	var v T
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		t.Fatalf("decode json: %v", err)
	}
	return v
}

func connectToDatabase(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		panic("Failed to connect to PostgreSQL: " + err.Error())
	}
	if err = db.Ping(); err != nil {
		panic("Failed to ping PostgreSQL: " + err.Error())
	}
	return db, nil
}

func startUpServer() {
	if startedOnce {
		return
	}
	startedOnce = true

	go func() {
		cfg := config.Load()

		db, err := connectToDatabase(cfg)
		if err != nil {
			log.Fatalf("Failed to connect to DB: %v", err)
		}

		r := api.RegisterRoutes(cfg, db)

		log.Printf("[E2E] Starting test server on port %s", cfg.Port)
		if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()
}

func waitForHealth(t *testing.T) {
	t.Helper()
	c := &http.Client{Timeout: 2 * time.Second}
	deadline := time.Now().Add(12 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := c.Get(baseURL + "/health")
		if err == nil && resp != nil && resp.StatusCode == http.StatusOK {
			_ = resp.Body.Close()
			return
		}
		if resp != nil {
			_ = resp.Body.Close()
		}
		time.Sleep(250 * time.Millisecond)
	}
	t.Fatalf("server did not become healthy at %s/health", baseURL)
}

func TestEndToEnd(t *testing.T) {
	os.Setenv("COOKIE_SECURE", "false")
	go startUpServer()

	// This is to keep the cookie once it is set
	jar, _ := cookiejar.New(nil)
	client = &http.Client{Timeout: 5 * time.Second, Jar: jar}

	waitForHealth(t)

	// Auth tests
	t.Run("Login", testEndToEnd_Login)
	t.Run("Me", testEndToEnd_Me)
	t.Run("Logout", testEndToEnd_Logout)
	t.Run("MeFail", testEndToEnd_MeFail)
	t.Run("Login", testEndToEnd_Login)

	// User Tests
	t.Run("Users_Create", testEndToEnd_Users_Create)
	t.Run("Users_List", testEndToEnd_Users_List)
	t.Run("Users_GetByID", testEndToEnd_Users_GetByID)
	t.Run("Users_UpdateByID", testEndToEnd_Users_UpdateByID)
	t.Run("Users_Create_Invalid", testEndToEnd_Users_Create_Invalid)
	t.Run("Users_Create_Duplicate", testEndToEnd_Users_Create_Duplicate)
	t.Run("Users_Delete_Success", testEndToEnd_Users_Delete_Success)
	t.Run("Users_GetByID_NotFound", testEndToEnd_Users_GetByID_NotFound)

	// Exercise Tests
	t.Run("Exercises_List_WithQueryParams", testEndToEnd_Exercises_List_WithQueryParams)
	t.Run("Exercises_GetByID", testEndToEnd_Exercises_GetByID)
	t.Run("Exercises_GetByID_NotFound", testEndToEnd_Exercises_GetByID_NotFound)

	// Routine Tests
	t.Run("Routines_Create", testEndToEnd_Routines_Create)
	t.Run("Routines_ListForUser", testEndToEnd_Routines_ListForUser)
	t.Run("Routines_ReadRoutineWithExercises", testEndToEnd_Routines_ReadRoutineWithExercises)
	t.Run("Routines_AddExerciseToRoutine", testEndToEnd_Routines_AddExerciseToRoutine)
	t.Run("Routines_RemoveExerciseFromRoutine", testEndToEnd_Routines_RemoveExerciseFromRoutine)
	t.Run("Routines_DeleteRoutine", testEndToEnd_Routines_DeleteRoutine)
	t.Run("Routines_DeleteUserRoutine", testEndToEnd_Routines_DeleteUserRoutine)

	// Schedule Tests
	t.Run("Schedules_GetByID_NotFound", testEndToEnd_Schedules_GetByID_NotFound)
	t.Run("Schedules_Create", testEndToEnd_Schedules_Create)
	t.Run("Schedules_ListMine", testEndToEnd_Schedules_ListMine)
	t.Run("Schedules_ListMineByDay", testEndToEnd_Schedules_ListMineByDay)
	t.Run("Schedules_GetByID", testEndToEnd_Schedules_GetByID)
	t.Run("Schedules_Update", testEndToEnd_Schedules_Update)
	t.Run("Schedules_Delete_OK", testEndToEnd_Schedules_Delete_OK)
	t.Run("Schedules_Delete_Idempotent", testEndToEnd_Schedules_Delete_Idempotent)
}
