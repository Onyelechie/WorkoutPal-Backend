package model

type Workout struct {
	ID        int64                `json:"id"`
	Name      string               `json:"name"`
	Frequency string               `json:"frequency"`
	NextRound string               `json:"nextRound"`
	Exercises []RegisteredExercise `json:"exercises"`
}

type RegisteredExercise struct {
	StartTime string   `json:"startTime"`
	EndTime   string   `json:"endTime"`
	Count     int      `json:"count"`
	Sets      int      `json:"sets"`
	Duration  int      `json:"duration"`
	Exercise  Exercise `json:"exercise"`
}

type CreateWorkoutRequest struct {
	Name      string               `json:"name"`
	Frequency string               `json:"frequency"`
	NextRound string               `json:"nextRound"`
	Exercises []RegisteredExercise `json:"exercises"`
}

type UpdateWorkoutRequest struct {
	Name      string               `json:"name"`
	Frequency string               `json:"frequency"`
	NextRound string               `json:"nextRound"`
	Exercises []RegisteredExercise `json:"exercises"`
}
