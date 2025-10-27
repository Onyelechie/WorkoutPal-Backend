package e2e

import (
	"net/http"
	"testing"
)

type exerciseRoutine struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"userId"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Exercises   []exercise `json:"exercises"`
	ExerciseIDs []int64    `json:"exerciseIds"`
	CreatedAt   string     `json:"createdAt"`
	IsActive    bool       `json:"isActive"`
}

type basicResponse struct {
	Message string `json:"message"`
}

type createRoutineReq struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ExerciseIDs []int64 `json:"exerciseIds"`
}

func testEndToEnd_Routines_Create(t *testing.T) {
	body := createRoutineReq{
		Name:        "E2E Routine " + randStringAlphaNum(6),
		Description: "auto-generated test routine",
		ExerciseIDs: []int64{1, 2},
	}

	resp := doRequest(t, http.MethodPost, "/users/1/routines", body, nil)
	defer resp.Body.Close()

	mustStatus(t, resp, http.StatusCreated)

	created := mustDecode[exerciseRoutine](t, resp)
	if created.ID == 0 {
		t.Fatalf("expected non-zero routine id")
	}
	if created.Name != body.Name {
		t.Fatalf("expected name=%q got=%q", body.Name, created.Name)
	}
	if len(created.ExerciseIDs) == 0 {
		t.Fatalf("expected routine to have exercises")
	}
}

func testEndToEnd_Routines_ListForUser(t *testing.T) {
	resp := doRequest(t, http.MethodGet, "/users/1/routines", nil, nil)
	defer resp.Body.Close()

	mustStatus(t, resp, http.StatusOK)

	list := mustDecode[[]exerciseRoutine](t, resp)
	if len(list) == 0 {
		t.Fatalf("expected at least one routine")
	}
	if list[0].ID == 0 {
		t.Fatalf("expected routine to have ID")
	}
	if list[0].Name == "" {
		t.Fatalf("expected routine to have Name")
	}
}

func testEndToEnd_Routines_ReadRoutineWithExercises(t *testing.T) {
	createBody := createRoutineReq{
		Name:        "E2E Detail " + randStringAlphaNum(6),
		Description: "routine to fetch with exercises",
		ExerciseIDs: []int64{1, 2},
	}

	createResp := doRequest(t, http.MethodPost, "/users/1/routines", createBody, nil)
	defer createResp.Body.Close()
	mustStatus(t, createResp, http.StatusCreated)

	created := mustDecode[exerciseRoutine](t, createResp)
	if created.ID == 0 {
		t.Fatalf("expected created routine id != 0")
	}

	getResp := doRequest(t, http.MethodGet, "/routines/"+int64ToStr(created.ID), nil, nil)
	defer getResp.Body.Close()
	mustStatus(t, getResp, http.StatusOK)

	got := mustDecode[exerciseRoutine](t, getResp)
	if got.ID != created.ID {
		t.Fatalf("expected id=%d got=%d", created.ID, got.ID)
	}
	if len(got.ExerciseIDs) == 0 {
		t.Fatalf("expected exercises on routine")
	}
}

func testEndToEnd_Routines_AddExerciseToRoutine(t *testing.T) {
	createBody := createRoutineReq{
		Name:        "E2E AddEx " + randStringAlphaNum(6),
		Description: "add exercise flow",
		ExerciseIDs: []int64{1},
	}

	createResp := doRequest(t, http.MethodPost, "/users/1/routines", createBody, nil)
	defer createResp.Body.Close()
	mustStatus(t, createResp, http.StatusCreated)

	created := mustDecode[exerciseRoutine](t, createResp)
	if created.ID == 0 {
		t.Fatalf("expected created routine id != 0")
	}

	addResp := doRequest(t, http.MethodPost, "/routines/"+int64ToStr(created.ID)+"/exercises?exercise_id=2", nil, nil)
	defer addResp.Body.Close()
	mustStatus(t, addResp, http.StatusOK)

	msg := mustDecode[basicResponse](t, addResp)
	if msg.Message == "" {
		t.Fatalf("expected success message, got empty")
	}

	checkResp := doRequest(t, http.MethodGet, "/routines/"+int64ToStr(created.ID), nil, nil)
	defer checkResp.Body.Close()
	mustStatus(t, checkResp, http.StatusOK)

	got := mustDecode[exerciseRoutine](t, checkResp)
	if len(got.ExerciseIDs) < 2 {
		t.Fatalf("expected >=2 exercises after add, got %d", len(got.ExerciseIDs))
	}
}

func testEndToEnd_Routines_RemoveExerciseFromRoutine(t *testing.T) {
	createBody := createRoutineReq{
		Name:        "E2E RemoveEx " + randStringAlphaNum(6),
		Description: "remove exercise flow",
		ExerciseIDs: []int64{1, 2},
	}

	createResp := doRequest(t, http.MethodPost, "/users/1/routines", createBody, nil)
	defer createResp.Body.Close()
	mustStatus(t, createResp, http.StatusCreated)

	created := mustDecode[exerciseRoutine](t, createResp)
	if created.ID == 0 {
		t.Fatalf("expected created routine id != 0")
	}
	if len(created.ExerciseIDs) < 2 {
		t.Fatalf("expected at least 2 exercises at creation")
	}

	exID := created.ExerciseIDs[0]

	delExResp := doRequest(t, http.MethodDelete, "/routines/"+int64ToStr(created.ID)+"/exercises/"+int64ToStr(exID), nil, nil)
	defer delExResp.Body.Close()
	mustStatus(t, delExResp, http.StatusOK)

	msg := mustDecode[basicResponse](t, delExResp)
	if msg.Message == "" {
		t.Fatalf("expected success message, got empty")
	}

	checkResp := doRequest(t, http.MethodGet, "/routines/"+int64ToStr(created.ID), nil, nil)
	defer checkResp.Body.Close()
	mustStatus(t, checkResp, http.StatusOK)

	got := mustDecode[exerciseRoutine](t, checkResp)
	if len(got.ExerciseIDs) != len(created.ExerciseIDs)-1 {
		t.Fatalf("expected %d exercises after removal, got %d", len(created.ExerciseIDs)-1, len(got.ExerciseIDs))
	}
}

func testEndToEnd_Routines_DeleteRoutine(t *testing.T) {
	createBody := createRoutineReq{
		Name:        "E2E DeleteRoutine " + randStringAlphaNum(6),
		Description: "delete flow",
		ExerciseIDs: []int64{1},
	}

	createResp := doRequest(t, http.MethodPost, "/users/1/routines", createBody, nil)
	defer createResp.Body.Close()
	mustStatus(t, createResp, http.StatusCreated)

	created := mustDecode[exerciseRoutine](t, createResp)
	if created.ID == 0 {
		t.Fatalf("expected created routine id != 0")
	}

	delResp := doRequest(t, http.MethodDelete, "/routines/"+int64ToStr(created.ID), nil, nil)
	defer delResp.Body.Close()
	mustStatus(t, delResp, http.StatusOK)

	msg := mustDecode[basicResponse](t, delResp)
	if msg.Message == "" {
		t.Fatalf("expected success message, got empty")
	}

	checkResp := doRequest(t, http.MethodGet, "/routines/"+int64ToStr(created.ID), nil, nil)
	defer checkResp.Body.Close()

	if checkResp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", checkResp.StatusCode)
	}
}

func testEndToEnd_Routines_DeleteUserRoutine(t *testing.T) {
	createBody := createRoutineReq{
		Name:        "E2E DeleteUserRoutine " + randStringAlphaNum(6),
		Description: "user-specific delete flow",
		ExerciseIDs: []int64{1},
	}

	createResp := doRequest(t, http.MethodPost, "/users/1/routines", createBody, nil)
	defer createResp.Body.Close()
	mustStatus(t, createResp, http.StatusCreated)

	created := mustDecode[exerciseRoutine](t, createResp)
	if created.ID == 0 {
		t.Fatalf("expected created routine id != 0")
	}

	delResp := doRequest(t, http.MethodDelete, "/users/1/routines/"+int64ToStr(created.ID), nil, nil)
	defer delResp.Body.Close()
	mustStatus(t, delResp, http.StatusOK)

	msg := mustDecode[basicResponse](t, delResp)
	if msg.Message == "" {
		t.Fatalf("expected success message, got empty")
	}
}
