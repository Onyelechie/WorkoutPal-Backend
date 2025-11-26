package model

type ExerciseSetting struct {
	UserID           int64   `json:"userId"`
	ExerciseID       int64   `json:"exerciseId"`
	WorkoutRoutineID int64   `json:"workoutRoutineId"`
	Weight           float64 `json:"weight"`
	Reps             int64   `json:"reps"`
	Sets             int64   `json:"sets"`
	BreakInterval    int64   `json:"breakInterval"`
}

type ReadExerciseSettingRequest struct {
	UserID           int64 `json:"userId"`
	ExerciseID       int64 `json:"exerciseId"`
	WorkoutRoutineID int64 `json:"workoutRoutineId"`
}

type CreateExerciseSettingRequest struct {
	UserID           int64   `json:"userId"`
	ExerciseID       int64   `json:"exerciseId"`
	WorkoutRoutineID int64   `json:"workoutRoutineId"`
	Weight           float64 `json:"weight"`
	Reps             int64   `json:"reps"`
	Sets             int64   `json:"sets"`
	BreakInterval    int64   `json:"breakInterval"`
}

type UpdateExerciseSettingRequest struct {
	UserID           int64   `json:"userId"`
	ExerciseID       int64   `json:"exerciseId"`
	WorkoutRoutineID int64   `json:"workoutRoutineId"`
	Weight           float64 `json:"weight"`
	Reps             int64   `json:"reps"`
	Sets             int64   `json:"sets"`
	BreakInterval    int64   `json:"breakInterval"`
}
