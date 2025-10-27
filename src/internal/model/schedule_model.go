package model

type Schedule struct {
	ID                   int64   `json:"id"`
	Name                 string  `json:"name"`
	UserID               int64   `json:"userId"`
	DayOfWeek            int64   `json:"dayOfWeek"`
	RoutineIDs           []int64 `json:"routineIds"`
	TimeSlot             string  `json:"timeSlot"`
	RoutineLengthMinutes int64   `json:"routineLengthMinutes"`
}

type CreateScheduleRequest struct {
	Name                 string  `json:"name"`
	UserID               int64   `json:"userId"`
	DayOfWeek            int64   `json:"dayOfWeek"`
	RoutineIDs           []int64 `json:"routineIds"`
	TimeSlot             string  `json:"timeSlot"`
	RoutineLengthMinutes int64   `json:"routineLengthMinutes"`
}

type UpdateScheduleRequest struct {
	ID                   int64   `json:"id"`
	Name                 string  `json:"name"`
	UserID               int64   `json:"userId"`
	DayOfWeek            int64   `json:"dayOfWeek"`
	RoutineIDs           []int64 `json:"routineIds"`
	TimeSlot             string  `json:"timeSlot"`
	RoutineLengthMinutes int64   `json:"routineLengthMinutes"`
}

type DeleteScheduleRequest struct {
	ID int64 `json:"id"`
}
