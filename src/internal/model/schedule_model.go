package model

import "time"

// user can have multiple schedules per day but can have multiple routines per schedule
// assumption is that for one schedule, the list of routines are consecutive of each other. i.e. start RoutineID[1] right after RoutineID[0] is done.
type Schedule struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	UserID        int64     `json:"userId"`
	DayOfWeek     int64     `json:"dayOfWeek"`     // index from 0 to 6. Sunday to Saturday. e.g. Sunday = 0, Saturday = 6.
	RoutineID     []int64   `json:"routineId"`     // list of routines that are consecutive of each other in the time slot. index 0 starts first
	TimeSlot      time.Time `json:"timeSlot"`      // start time of list of routines
	RoutineLength time.Time `json:"routineLength"` // total length of list of routines in minutes
}

type CreateScheduleRequest struct {
	Name          string    `json:"name"`
	UserID        int64     `json:"userId"`
	DayOfWeek     int64     `json:"dayOfWeek"`
	RoutineID     []int64   `json:"routineId"`
	TimeSlot      time.Time `json:"timeSlot"`
	RoutineLength time.Time `json:"routineLength"`
}

type UpdateScheduleRequest struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	UserID        int64     `json:"userId"`
	DayOfWeek     int64     `json:"dayOfWeek"`
	RoutineID     []int64   `json:"routineId"`
	TimeSlot      time.Time `json:"timeSlot"`
	RoutineLength time.Time `json:"routineLength"`
}

type DeleteScheduleRequest struct {
	ID int64 `json:"id"`
}
