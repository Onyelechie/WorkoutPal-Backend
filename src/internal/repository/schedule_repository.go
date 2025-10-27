package repository

import (
	"context"
	"database/sql"
	"workoutpal/src/internal/domain/repository"
	"workoutpal/src/internal/model"
)

type scheduleRepository struct {
	db *sql.DB
}

func NewScheduleRepository(db *sql.DB) repository.ScheduleRepository {
	return &scheduleRepository{db: db}
}

func (s *scheduleRepository) getRoutineIDsForSchedule(ctx context.Context, scheduleID int64) ([]int64, error) {
	const q = `
		SELECT routine_id
		FROM schedule_routine
		WHERE schedule_id = $1
		ORDER BY position ASC;
	`

	rows, err := s.db.QueryContext(ctx, q, scheduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routineIDs []int64
	for rows.Next() {
		var rid int64
		if err := rows.Scan(&rid); err != nil {
			return nil, err
		}
		routineIDs = append(routineIDs, rid)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return routineIDs, nil
}

func (s *scheduleRepository) hydrateSchedule(ctx context.Context, base *model.Schedule) (*model.Schedule, error) {
	routineIDs, err := s.getRoutineIDsForSchedule(ctx, base.ID)
	if err != nil {
		return nil, err
	}
	base.RoutineIDs = routineIDs
	return base, nil
}

func (s *scheduleRepository) insertScheduleRoutines(ctx context.Context, tx *sql.Tx, scheduleID int64, routineIDs []int64) error {
	const q = `
		INSERT INTO schedule_routine (schedule_id, routine_id, position)
		VALUES ($1, $2, $3);
	`
	for i, rid := range routineIDs {
		if _, err := tx.ExecContext(ctx, q, scheduleID, rid, i); err != nil {
			return err
		}
	}
	return nil
}

func (s *scheduleRepository) replaceScheduleRoutines(ctx context.Context, tx *sql.Tx, scheduleID int64, routineIDs []int64) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM schedule_routine WHERE schedule_id = $1`, scheduleID); err != nil {
		return err
	}
	return s.insertScheduleRoutines(ctx, tx, scheduleID, routineIDs)
}

func scanScheduleRow(row Scanner) (*model.Schedule, error) {
	var sch model.Schedule
	if err := row.Scan(
		&sch.ID,
		&sch.Name,
		&sch.UserID,
		&sch.DayOfWeek,
		&sch.TimeSlot,
		&sch.RoutineLengthMinutes,
	); err != nil {
		return nil, err
	}
	return &sch, nil
}

type Scanner interface {
	Scan(dest ...any) error
}

func (s *scheduleRepository) ReadUserSchedules(userId int64) ([]*model.Schedule, error) {
	ctx := context.Background()

	const q = `
		SELECT id, name, user_id, day_of_week, time_slot, routine_length_minutes
		FROM schedule
		WHERE user_id = $1
		ORDER BY day_of_week ASC, time_slot ASC;
	`

	rows, err := s.db.QueryContext(ctx, q, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []*model.Schedule

	for rows.Next() {
		base, err := scanScheduleRow(rows)
		if err != nil {
			return nil, err
		}
		full, err := s.hydrateSchedule(ctx, base)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, full)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}

func (s *scheduleRepository) ReadUserSchedulesByDay(userId int64, dayOfWeek int64) ([]*model.Schedule, error) {
	ctx := context.Background()

	const q = `
		SELECT id, name, user_id, day_of_week, time_slot, routine_length_minutes
		FROM schedule
		WHERE user_id = $1
		  AND day_of_week = $2
		ORDER BY time_slot ASC;
	`

	rows, err := s.db.QueryContext(ctx, q, userId, dayOfWeek)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []*model.Schedule

	for rows.Next() {
		base, err := scanScheduleRow(rows)
		if err != nil {
			return nil, err
		}
		full, err := s.hydrateSchedule(ctx, base)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, full)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}

func (s *scheduleRepository) ReadScheduleByID(id int64) (*model.Schedule, error) {
	ctx := context.Background()

	const q = `
		SELECT id, name, user_id, day_of_week, time_slot, routine_length_minutes
		FROM schedule
		WHERE id = $1;
	`

	row := s.db.QueryRowContext(ctx, q, id)

	base, err := scanScheduleRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return s.hydrateSchedule(ctx, base)
}

func (s *scheduleRepository) CreateSchedule(request model.CreateScheduleRequest) (*model.Schedule, error) {
	ctx := context.Background()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	const insertSchedule = `
		INSERT INTO schedule (name, user_id, day_of_week, time_slot, routine_length_minutes)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, user_id, day_of_week, time_slot, routine_length_minutes;
	`

	row := tx.QueryRowContext(ctx, insertSchedule,
		request.Name,
		request.UserID,
		request.DayOfWeek,
		request.TimeSlot,
		request.RoutineLengthMinutes,
	)

	base, err := scanScheduleRow(row)
	if err != nil {
		return nil, err
	}

	if err := s.insertScheduleRoutines(ctx, tx, base.ID, request.RoutineIDs); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.hydrateSchedule(ctx, base)
}

func (s *scheduleRepository) UpdateSchedule(request model.UpdateScheduleRequest) (*model.Schedule, error) {
	ctx := context.Background()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	const updateSchedule = `
		UPDATE schedule
		SET name = $1,
		    day_of_week = $2,
		    time_slot = $3,
		    routine_length_minutes = $4
		WHERE id = $5
		  AND user_id = $6
		RETURNING id, name, user_id, day_of_week, time_slot, routine_length_minutes;
	`

	row := tx.QueryRowContext(ctx, updateSchedule,
		request.Name,
		request.DayOfWeek,
		request.TimeSlot,
		request.RoutineLengthMinutes,
		request.ID,
		request.UserID,
	)

	base, err := scanScheduleRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if err := s.replaceScheduleRoutines(ctx, tx, base.ID, request.RoutineIDs); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.hydrateSchedule(ctx, base)
}

func (s *scheduleRepository) DeleteSchedule(request model.DeleteScheduleRequest) error {
	ctx := context.Background()

	const q = `
		DELETE FROM schedule
		WHERE id = $1;
	`

	res, err := s.db.ExecContext(ctx, q, request.ID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
