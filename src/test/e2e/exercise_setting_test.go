package e2e

import (
	"net/http"
	"testing"
)

type exerciseSetting struct {
	UserID           int64 `json:"userId"`
	ExerciseID       int64 `json:"exerciseId"`
	WorkoutRoutineID int64 `json:"workoutRoutineId"`
	Weight           int   `json:"weight"`
	Reps             int   `json:"reps"`
	Sets             int   `json:"sets"`
	BreakInterval    int   `json:"breakInterval"`
}

type createExerciseSettingReq struct {
	ExerciseID       int64 `json:"exerciseId"`
	WorkoutRoutineID int64 `json:"workoutRoutineId"`
	Weight           int   `json:"weight"`
	Reps             int   `json:"reps"`
	Sets             int   `json:"sets"`
	BreakInterval    int   `json:"breakInterval"`
}

type updateExerciseSettingReq struct {
	ExerciseID       int64 `json:"exerciseId"`
	WorkoutRoutineID int64 `json:"workoutRoutineId"`
	Weight           int   `json:"weight"`
	Reps             int   `json:"reps"`
	Sets             int   `json:"sets"`
	BreakInterval    int   `json:"breakInterval"`
}

func testEndToEnd_ExerciseSettings_CreateAndRead(t *testing.T) {
	exerciseID := int64(1)
	workoutRoutineID := int64(1)

	createBody := createExerciseSettingReq{
		ExerciseID:       exerciseID,
		WorkoutRoutineID: workoutRoutineID,
		Weight:           50,
		Reps:             8,
		Sets:             3,
		BreakInterval:    90,
	}

	createResp := doRequest(t, http.MethodPost, "/exercise-settings", createBody, nil)
	defer createResp.Body.Close()
	mustStatus(t, createResp, http.StatusCreated)

	created := mustDecode[exerciseSetting](t, createResp)
	if created.ExerciseID != exerciseID || created.WorkoutRoutineID != workoutRoutineID {
		t.Fatalf("expected exerciseId=%d workoutRoutineId=%d, got %+v", exerciseID, workoutRoutineID, created)
	}

	url := "/exercise-settings?exercise_id=" + int64ToStr(exerciseID) +
		"&workout_routine_id=" + int64ToStr(workoutRoutineID)

	readResp := doRequest(t, http.MethodGet, url, nil, nil)
	defer readResp.Body.Close()
	mustStatus(t, readResp, http.StatusOK)

	read := mustDecode[exerciseSetting](t, readResp)
	if read.ExerciseID != exerciseID || read.WorkoutRoutineID != workoutRoutineID {
		t.Fatalf("expected exerciseId=%d workoutRoutineId=%d, got %+v", exerciseID, workoutRoutineID, read)
	}
	if read.Weight != createBody.Weight || read.Reps != createBody.Reps ||
		read.Sets != createBody.Sets || read.BreakInterval != createBody.BreakInterval {
		t.Fatalf("expected values %+v, got %+v", createBody, read)
	}
}

func testEndToEnd_ExerciseSettings_Update(t *testing.T) {
	exerciseID := int64(2)
	workoutRoutineID := int64(1)

	createBody := createExerciseSettingReq{
		ExerciseID:       exerciseID,
		WorkoutRoutineID: workoutRoutineID,
		Weight:           40,
		Reps:             6,
		Sets:             3,
		BreakInterval:    60,
	}
	createResp := doRequest(t, http.MethodPost, "/exercise-settings", createBody, nil)
	defer createResp.Body.Close()
	mustStatus(t, createResp, http.StatusCreated)

	_ = mustDecode[exerciseSetting](t, createResp)

	updateBody := updateExerciseSettingReq{
		ExerciseID:       exerciseID,
		WorkoutRoutineID: workoutRoutineID,
		Weight:           70,
		Reps:             10,
		Sets:             4,
		BreakInterval:    120,
	}

	updateResp := doRequest(t, http.MethodPut, "/exercise-settings", updateBody, nil)
	defer updateResp.Body.Close()
	mustStatus(t, updateResp, http.StatusOK)

	updated := mustDecode[exerciseSetting](t, updateResp)
	if updated.Weight != updateBody.Weight ||
		updated.Reps != updateBody.Reps ||
		updated.Sets != updateBody.Sets ||
		updated.BreakInterval != updateBody.BreakInterval {
		t.Fatalf("expected updated values %+v, got %+v", updateBody, updated)
	}

	url := "/exercise-settings?exercise_id=" + int64ToStr(exerciseID) +
		"&workout_routine_id=" + int64ToStr(workoutRoutineID)

	readResp := doRequest(t, http.MethodGet, url, nil, nil)
	defer readResp.Body.Close()
	mustStatus(t, readResp, http.StatusOK)

	read := mustDecode[exerciseSetting](t, readResp)
	if read.Weight != updateBody.Weight ||
		read.Reps != updateBody.Reps ||
		read.Sets != updateBody.Sets ||
		read.BreakInterval != updateBody.BreakInterval {
		t.Fatalf("expected persisted updated values %+v, got %+v", updateBody, read)
	}
}
