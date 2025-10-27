package repository

import (
	"database/sql"
	"regexp"
	"testing"
	"workoutpal/src/internal/model"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestScheduleRepository_ReadUserSchedules_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewScheduleRepository(db)

	userID := int64(10)

	rowsSchedules := sqlmock.NewRows([]string{
		"id", "name", "user_id", "day_of_week", "time_slot", "routine_length_minutes",
	}).AddRow(1, "Morning Lift", userID, 1, "07:30", 90).
		AddRow(2, "Cardio", userID, 3, "18:00", 45)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, user_id, day_of_week, time_slot, routine_length_minutes
		FROM schedule
		WHERE user_id = $1
		ORDER BY day_of_week ASC, time_slot ASC;
	`)).
		WithArgs(userID).
		WillReturnRows(rowsSchedules)

	rowsRoutine1 := sqlmock.NewRows([]string{"routine_id"}).
		AddRow(100).
		AddRow(101)
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT routine_id
		FROM schedule_routine
		WHERE schedule_id = $1
		ORDER BY position ASC;
	`)).
		WithArgs(int64(1)).
		WillReturnRows(rowsRoutine1)

	rowsRoutine2 := sqlmock.NewRows([]string{"routine_id"}).
		AddRow(200)
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT routine_id
		FROM schedule_routine
		WHERE schedule_id = $1
		ORDER BY position ASC;
	`)).
		WithArgs(int64(2)).
		WillReturnRows(rowsRoutine2)

	got, err := repo.ReadUserSchedules(userID)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 schedules, got %d", len(got))
	}

	if got[0].ID != 1 || got[0].Name != "Morning Lift" || got[0].UserID != userID {
		t.Fatalf("bad first schedule: %#v", got[0])
	}
	if got[0].DayOfWeek != 1 || got[0].TimeSlot != "07:30" || got[0].RoutineLengthMinutes != 90 {
		t.Fatalf("bad first schedule fields: %#v", got[0])
	}
	if len(got[0].RoutineIDs) != 2 || got[0].RoutineIDs[0] != 100 || got[0].RoutineIDs[1] != 101 {
		t.Fatalf("bad first schedule routines: %#v", got[0].RoutineIDs)
	}

	if got[1].ID != 2 || got[1].Name != "Cardio" || got[1].UserID != userID {
		t.Fatalf("bad second schedule: %#v", got[1])
	}
	if len(got[1].RoutineIDs) != 1 || got[1].RoutineIDs[0] != 200 {
		t.Fatalf("bad second schedule routines: %#v", got[1].RoutineIDs)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestScheduleRepository_ReadUserSchedules_QueryError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewScheduleRepository(db)

	userID := int64(10)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, user_id, day_of_week, time_slot, routine_length_minutes
		FROM schedule
		WHERE user_id = $1
		ORDER BY day_of_week ASC, time_slot ASC;
	`)).
		WithArgs(userID).
		WillReturnError(assertErr)

	got, err := repo.ReadUserSchedules(userID)
	if got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
	if err == nil || err.Error() != assertErr.Error() {
		t.Fatalf("expected %q, got %v", assertErr.Error(), err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestScheduleRepository_ReadUserSchedules_ScanError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewScheduleRepository(db)

	userID := int64(10)

	rowsSchedules := sqlmock.NewRows([]string{
		"id", "name", "user_id", "day_of_week", "time_slot", "routine_length_minutes",
	}).AddRow("bad", "X", userID, 1, "07:30", 90)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, user_id, day_of_week, time_slot, routine_length_minutes
		FROM schedule
		WHERE user_id = $1
		ORDER BY day_of_week ASC, time_slot ASC;
	`)).
		WithArgs(userID).
		WillReturnRows(rowsSchedules)

	got, err := repo.ReadUserSchedules(userID)
	if got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
	if err == nil {
		t.Fatalf("expected scan error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestScheduleRepository_ReadUserSchedulesByDay_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewScheduleRepository(db)

	userID := int64(10)
	day := int64(3)

	rowsSchedules := sqlmock.NewRows([]string{
		"id", "name", "user_id", "day_of_week", "time_slot", "routine_length_minutes",
	}).AddRow(5, "Evening Cardio", userID, day, "19:00", 60)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, user_id, day_of_week, time_slot, routine_length_minutes
		FROM schedule
		WHERE user_id = $1
		  AND day_of_week = $2
		ORDER BY time_slot ASC;
	`)).
		WithArgs(userID, day).
		WillReturnRows(rowsSchedules)

	rowsRoutine := sqlmock.NewRows([]string{"routine_id"}).
		AddRow(99).
		AddRow(100)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT routine_id
		FROM schedule_routine
		WHERE schedule_id = $1
		ORDER BY position ASC;
	`)).
		WithArgs(int64(5)).
		WillReturnRows(rowsRoutine)

	got, err := repo.ReadUserSchedulesByDay(userID, day)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 schedule, got %d", len(got))
	}
	if got[0].ID != 5 || got[0].Name != "Evening Cardio" {
		t.Fatalf("bad schedule: %#v", got[0])
	}
	if len(got[0].RoutineIDs) != 2 || got[0].RoutineIDs[0] != 99 || got[0].RoutineIDs[1] != 100 {
		t.Fatalf("bad routine ids: %#v", got[0].RoutineIDs)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestScheduleRepository_ReadScheduleByID_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewScheduleRepository(db)

	rowsSchedule := sqlmock.NewRows([]string{
		"id", "name", "user_id", "day_of_week", "time_slot", "routine_length_minutes",
	}).AddRow(8, "Leg Day", 99, 2, "06:00", 75)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, user_id, day_of_week, time_slot, routine_length_minutes
		FROM schedule
		WHERE id = $1;
	`)).
		WithArgs(int64(8)).
		WillReturnRows(rowsSchedule)

	rowsRoutine := sqlmock.NewRows([]string{"routine_id"}).
		AddRow(501).
		AddRow(777)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT routine_id
		FROM schedule_routine
		WHERE schedule_id = $1
		ORDER BY position ASC;
	`)).
		WithArgs(int64(8)).
		WillReturnRows(rowsRoutine)

	got, err := repo.ReadScheduleByID(8)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got == nil || got.ID != 8 || got.Name != "Leg Day" || got.UserID != 99 {
		t.Fatalf("bad schedule: %#v", got)
	}
	if len(got.RoutineIDs) != 2 || got.RoutineIDs[0] != 501 || got.RoutineIDs[1] != 777 {
		t.Fatalf("bad routine ids: %#v", got.RoutineIDs)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestScheduleRepository_ReadScheduleByID_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewScheduleRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, user_id, day_of_week, time_slot, routine_length_minutes
		FROM schedule
		WHERE id = $1;
	`)).
		WithArgs(int64(123)).
		WillReturnError(sql.ErrNoRows)

	got, err := repo.ReadScheduleByID(123)
	if got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestScheduleRepository_ReadScheduleByID_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewScheduleRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, user_id, day_of_week, time_slot, routine_length_minutes
		FROM schedule
		WHERE id = $1;
	`)).
		WithArgs(int64(55)).
		WillReturnError(assertErr)

	got, err := repo.ReadScheduleByID(55)
	if got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
	if err == nil || err.Error() != assertErr.Error() {
		t.Fatalf("expected %q, got %v", assertErr.Error(), err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestScheduleRepository_CreateSchedule_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewScheduleRepository(db)

	req := model.CreateScheduleRequest{
		Name:                 "AM Upper",
		UserID:               77,
		DayOfWeek:            1,
		TimeSlot:             "07:15",
		RoutineLengthMinutes: 80,
		RoutineIDs:           []int64{10, 11, 12},
	}

	mock.ExpectBegin()

	rowInsert := sqlmock.NewRows([]string{
		"id", "name", "user_id", "day_of_week", "time_slot", "routine_length_minutes",
	}).AddRow(900, req.Name, req.UserID, req.DayOfWeek, req.TimeSlot, req.RoutineLengthMinutes)

	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO schedule (name, user_id, day_of_week, time_slot, routine_length_minutes)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, user_id, day_of_week, time_slot, routine_length_minutes;
	`)).
		WithArgs(req.Name, req.UserID, req.DayOfWeek, req.TimeSlot, req.RoutineLengthMinutes).
		WillReturnRows(rowInsert)

	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO schedule_routine (schedule_id, routine_id, position)
		VALUES ($1, $2, $3);
	`)).
		WithArgs(int64(900), int64(10), 0).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO schedule_routine (schedule_id, routine_id, position)
		VALUES ($1, $2, $3);
	`)).
		WithArgs(int64(900), int64(11), 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO schedule_routine (schedule_id, routine_id, position)
		VALUES ($1, $2, $3);
	`)).
		WithArgs(int64(900), int64(12), 2).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	finalRoutineRows := sqlmock.NewRows([]string{"routine_id"}).
		AddRow(10).
		AddRow(11).
		AddRow(12)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT routine_id
		FROM schedule_routine
		WHERE schedule_id = $1
		ORDER BY position ASC;
	`)).
		WithArgs(int64(900)).
		WillReturnRows(finalRoutineRows)

	got, err := repo.CreateSchedule(req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got == nil || got.ID != 900 || got.Name != req.Name || got.UserID != req.UserID {
		t.Fatalf("unexpected schedule: %#v", got)
	}
	if len(got.RoutineIDs) != 3 || got.RoutineIDs[2] != 12 {
		t.Fatalf("unexpected routineIDs: %#v", got.RoutineIDs)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestScheduleRepository_CreateSchedule_InsertError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewScheduleRepository(db)

	req := model.CreateScheduleRequest{
		Name:                 "AM Upper",
		UserID:               77,
		DayOfWeek:            1,
		TimeSlot:             "07:15",
		RoutineLengthMinutes: 80,
		RoutineIDs:           []int64{10},
	}

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO schedule (name, user_id, day_of_week, time_slot, routine_length_minutes)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, user_id, day_of_week, time_slot, routine_length_minutes;
	`)).
		WithArgs(req.Name, req.UserID, req.DayOfWeek, req.TimeSlot, req.RoutineLengthMinutes).
		WillReturnError(assertErr)

	mock.ExpectRollback()

	got, err := repo.CreateSchedule(req)
	if got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
	if err == nil || err.Error() != assertErr.Error() {
		t.Fatalf("expected %q got %v", assertErr.Error(), err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestScheduleRepository_UpdateSchedule_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewScheduleRepository(db)

	req := model.UpdateScheduleRequest{
		ID:                   321,
		Name:                 "Updated Split",
		UserID:               44,
		DayOfWeek:            5,
		TimeSlot:             "20:00",
		RoutineLengthMinutes: 50,
		RoutineIDs:           []int64{7, 8},
	}

	mock.ExpectBegin()

	rowUpdate := sqlmock.NewRows([]string{
		"id", "name", "user_id", "day_of_week", "time_slot", "routine_length_minutes",
	}).AddRow(req.ID, req.Name, req.UserID, req.DayOfWeek, req.TimeSlot, req.RoutineLengthMinutes)

	mock.ExpectQuery(regexp.QuoteMeta(`
		UPDATE schedule
		SET name = $1,
		    day_of_week = $2,
		    time_slot = $3,
		    routine_length_minutes = $4
		WHERE id = $5
		  AND user_id = $6
		RETURNING id, name, user_id, day_of_week, time_slot, routine_length_minutes;
	`)).
		WithArgs(req.Name, req.DayOfWeek, req.TimeSlot, req.RoutineLengthMinutes, req.ID, req.UserID).
		WillReturnRows(rowUpdate)

	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM schedule_routine WHERE schedule_id = $1`,
	)).
		WithArgs(req.ID).
		WillReturnResult(sqlmock.NewResult(0, 2))

	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO schedule_routine (schedule_id, routine_id, position)
		VALUES ($1, $2, $3);
	`)).
		WithArgs(req.ID, int64(7), 0).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO schedule_routine (schedule_id, routine_id, position)
		VALUES ($1, $2, $3);
	`)).
		WithArgs(req.ID, int64(8), 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	finalRoutineRows := sqlmock.NewRows([]string{"routine_id"}).
		AddRow(7).
		AddRow(8)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT routine_id
		FROM schedule_routine
		WHERE schedule_id = $1
		ORDER BY position ASC;
	`)).
		WithArgs(req.ID).
		WillReturnRows(finalRoutineRows)

	got, err := repo.UpdateSchedule(req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got == nil || got.ID != req.ID || got.Name != req.Name {
		t.Fatalf("unexpected schedule: %#v", got)
	}
	if len(got.RoutineIDs) != 2 || got.RoutineIDs[1] != 8 {
		t.Fatalf("unexpected routineIDs: %#v", got.RoutineIDs)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestScheduleRepository_UpdateSchedule_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewScheduleRepository(db)

	req := model.UpdateScheduleRequest{
		ID:                   321,
		Name:                 "Updated Split",
		UserID:               44,
		DayOfWeek:            5,
		TimeSlot:             "20:00",
		RoutineLengthMinutes: 50,
		RoutineIDs:           []int64{7, 8},
	}

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(`
		UPDATE schedule
		SET name = $1,
		    day_of_week = $2,
		    time_slot = $3,
		    routine_length_minutes = $4
		WHERE id = $5
		  AND user_id = $6
		RETURNING id, name, user_id, day_of_week, time_slot, routine_length_minutes;
	`)).
		WithArgs(req.Name, req.DayOfWeek, req.TimeSlot, req.RoutineLengthMinutes, req.ID, req.UserID).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectRollback()

	got, err := repo.UpdateSchedule(req)
	if got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestScheduleRepository_DeleteSchedule_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewScheduleRepository(db)

	mock.ExpectExec(regexp.QuoteMeta(`
		DELETE FROM schedule
		WHERE id = $1;
	`)).
		WithArgs(int64(55)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteSchedule(model.DeleteScheduleRequest{ID: 55})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestScheduleRepository_DeleteSchedule_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewScheduleRepository(db)

	mock.ExpectExec(regexp.QuoteMeta(`
		DELETE FROM schedule
		WHERE id = $1;
	`)).
		WithArgs(int64(99)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.DeleteSchedule(model.DeleteScheduleRequest{ID: 99})
	if err == nil {
		t.Fatalf("expected err, got nil")
	}
	if err.Error() != sql.ErrNoRows.Error() {
		t.Fatalf("expected sql.ErrNoRows, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
