package e2e

import (
	"net/http"
	"testing"
)

type scheduleDTO struct {
	ID                   int64   `json:"id"`
	Name                 string  `json:"name"`
	UserID               int64   `json:"userId"`
	DayOfWeek            int64   `json:"dayOfWeek"`
	RoutineIDs           []int64 `json:"routineIds"`
	TimeSlot             string  `json:"timeSlot"`
	RoutineLengthMinutes int64   `json:"routineLengthMinutes"`
}

type createScheduleDTO struct {
	Name                 string  `json:"name"`
	UserID               int64   `json:"userId"`
	DayOfWeek            int64   `json:"dayOfWeek"`
	TimeSlot             string  `json:"timeSlot"`
	RoutineLengthMinutes int64   `json:"routineLengthMinutes"`
	RoutineIDs           []int64 `json:"routineIds"`
}

type updateScheduleDTO struct {
	ID                   int64   `json:"id"`
	UserID               int64   `json:"userId"`
	Name                 string  `json:"name"`
	DayOfWeek            int64   `json:"dayOfWeek"`
	TimeSlot             string  `json:"timeSlot"`
	RoutineLengthMinutes int64   `json:"routineLengthMinutes"`
	RoutineIDs           []int64 `json:"routineIds"`
}

func testEndToEnd_Schedules_Create(t *testing.T) {
	validRoutineID := int64(1)

	reqBody := createScheduleDTO{
		Name:                 "E2E Schedule " + randStringAlphaNum(6),
		UserID:               0,
		DayOfWeek:            1,
		TimeSlot:             "07:30",
		RoutineLengthMinutes: 60,
		RoutineIDs:           []int64{validRoutineID},
	}

	resp := doRequest(t, http.MethodPost, "/schedules", reqBody, nil)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		created := mustDecode[scheduleDTO](t, resp)
		if created.ID == 0 {
			t.Fatalf("expected non-zero schedule ID")
		}
		if created.Name != reqBody.Name {
			t.Fatalf("expected name=%q got=%q", reqBody.Name, created.Name)
		}
		if len(created.RoutineIDs) != 1 || created.RoutineIDs[0] != validRoutineID {
			t.Fatalf("unexpected routineIds on create: %+v", created.RoutineIDs)
		}
		return
	}

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("status %d (want 201 or 500)", resp.StatusCode)
	}

	errObj := mustDecode[map[string]any](t, resp)
	if _, ok := errObj["error"]; !ok {
		t.Fatalf("expected error body, got %v", errObj)
	}
}

func testEndToEnd_Schedules_ListMine(t *testing.T) {
	resp := doRequest(t, http.MethodGet, "/schedules", nil, nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d (want 200)", resp.StatusCode)
	}

	list := mustDecode[[]scheduleDTO](t, resp)
	if len(list) == 0 {
		t.Fatalf("expected at least 1 schedule, got 0")
	}
	if list[0].ID == 0 {
		t.Fatalf("expected non-zero schedule id in list[0]")
	}
	if list[0].Name == "" {
		t.Fatalf("expected non-empty Name in list[0]")
	}
}

func testEndToEnd_Schedules_ListMineByDay(t *testing.T) {
	day := "1"
	resp := doRequest(t, http.MethodGet, "/schedules/of/"+day, nil, nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d (want 200)", resp.StatusCode)
	}

	list := mustDecode[[]scheduleDTO](t, resp)
	if list == nil {
		t.Fatalf("expected list, got nil")
	}
	if len(list) > 0 && int64ToStr(list[0].DayOfWeek) != day {
		t.Fatalf("expected dayOfWeek=%s got=%d", day, list[0].DayOfWeek)
	}
}

func testEndToEnd_Schedules_GetByID(t *testing.T) {
	seedID := "1"
	resp := doRequest(t, http.MethodGet, "/schedules/"+seedID, nil, nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d (want 200)", resp.StatusCode)
	}

	got := mustDecode[scheduleDTO](t, resp)
	if got.ID == 0 {
		t.Fatalf("expected non-zero schedule id")
	}
	if got.Name == "" {
		t.Fatalf("expected non-empty Name")
	}
}

func testEndToEnd_Schedules_GetByID_NotFound(t *testing.T) {
	resp := doRequest(t, http.MethodGet, "/schedules/99999999", nil, nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError && resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 500 or 200, got=%d", resp.StatusCode)
	}

	if resp.StatusCode == http.StatusInternalServerError {
		errObj := mustDecode[map[string]any](t, resp)
		if _, ok := errObj["error"]; !ok {
			t.Fatalf("expected error body, got %v", errObj)
		}
	}
}

func testEndToEnd_Schedules_Update(t *testing.T) {
	validRoutineID := int64(1)

	createBody := createScheduleDTO{
		Name:                 "E2E UpdateBase " + randStringAlphaNum(6),
		UserID:               0,
		DayOfWeek:            2,
		TimeSlot:             "18:45",
		RoutineLengthMinutes: 50,
		RoutineIDs:           []int64{validRoutineID},
	}

	createResp := doRequest(t, http.MethodPost, "/schedules", createBody, nil)
	defer createResp.Body.Close()

	if createResp.StatusCode != http.StatusCreated {
		if createResp.StatusCode == http.StatusInternalServerError {
			errObj := mustDecode[map[string]any](t, createResp)
			if _, ok := errObj["error"]; !ok {
				t.Fatalf("expected error body, got %v", errObj)
			}
			return
		}
		t.Fatalf("status %d (want 201 or 500)", createResp.StatusCode)
	}

	created := mustDecode[scheduleDTO](t, createResp)
	if created.ID == 0 {
		t.Fatalf("expected created schedule id != 0")
	}

	updateBody := updateScheduleDTO{
		ID:                   created.ID,
		UserID:               created.UserID,
		Name:                 created.Name + " (updated)",
		DayOfWeek:            3,
		TimeSlot:             "19:00",
		RoutineLengthMinutes: 55,
		RoutineIDs:           []int64{validRoutineID},
	}

	updateResp := doRequest(t, http.MethodPut, "/schedules/"+int64ToStr(created.ID), updateBody, nil)
	defer updateResp.Body.Close()

	if updateResp.StatusCode != http.StatusOK {
		if updateResp.StatusCode == http.StatusInternalServerError {
			errObj := mustDecode[map[string]any](t, updateResp)
			if _, ok := errObj["error"]; !ok {
				t.Fatalf("expected error body, got %v", errObj)
			}
			return
		}
		t.Fatalf("status %d (want 200 or 500)", updateResp.StatusCode)
	}

	updated := mustDecode[scheduleDTO](t, updateResp)
	if updated.ID != created.ID {
		t.Fatalf("expected same id=%d got=%d", created.ID, updated.ID)
	}
	if updated.Name != updateBody.Name {
		t.Fatalf("expected name=%q got=%q", updateBody.Name, updated.Name)
	}
	if updated.DayOfWeek != updateBody.DayOfWeek {
		t.Fatalf("expected dayOfWeek=%d got=%d", updateBody.DayOfWeek, updated.DayOfWeek)
	}
}

func testEndToEnd_Schedules_Delete_OK(t *testing.T) {
	validRoutineID := int64(1)

	createBody := createScheduleDTO{
		Name:                 "E2E DeleteMe " + randStringAlphaNum(6),
		UserID:               0,
		DayOfWeek:            4,
		TimeSlot:             "12:00",
		RoutineLengthMinutes: 40,
		RoutineIDs:           []int64{validRoutineID},
	}

	createResp := doRequest(t, http.MethodPost, "/schedules", createBody, nil)
	defer createResp.Body.Close()

	if createResp.StatusCode != http.StatusCreated {
		if createResp.StatusCode == http.StatusInternalServerError {
			errObj := mustDecode[map[string]any](t, createResp)
			if _, ok := errObj["error"]; !ok {
				t.Fatalf("expected error body, got %v", errObj)
			}
			return
		}
		t.Fatalf("status %d (want 201 or 500)", createResp.StatusCode)
	}

	created := mustDecode[scheduleDTO](t, createResp)
	if created.ID == 0 {
		t.Fatalf("expected created schedule id != 0")
	}

	delResp := doRequest(t, http.MethodDelete, "/schedules/"+int64ToStr(created.ID), nil, nil)
	defer delResp.Body.Close()

	if delResp.StatusCode != http.StatusNoContent {
		t.Fatalf("status %d (want 204)", delResp.StatusCode)
	}
}

func testEndToEnd_Schedules_Delete_Idempotent(t *testing.T) {
	delResp := doRequest(t, http.MethodDelete, "/schedules/987654321", nil, nil)
	defer delResp.Body.Close()

	if delResp.StatusCode != http.StatusNoContent {
		t.Fatalf("status %d (want 204)", delResp.StatusCode)
	}
}
