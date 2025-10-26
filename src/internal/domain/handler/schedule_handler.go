package handler

import "net/http"

type ScheduleHandler interface {
	ReadUserSchedules(w http.ResponseWriter, r *http.Request)
	ReadUserSchedulesByDay(w http.ResponseWriter, r *http.Request)
	ReadScheduleByID(w http.ResponseWriter, r *http.Request)
	CreateSchedule(w http.ResponseWriter, r *http.Request)
	UpdateSchedule(w http.ResponseWriter, r *http.Request)
	DeleteSchedule(w http.ResponseWriter, r *http.Request)
}
