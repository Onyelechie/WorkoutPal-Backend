package repository

import (
	"database/sql"
	"errors"
	"workoutpal/src/internal/domain/repository"
	"workoutpal/src/internal/model"
)

type goalRepository struct {
	db *sql.DB
}

func NewGoalRepository(db *sql.DB) repository.GoalRepository {
	return &goalRepository{db: db}
}

func (g *goalRepository) CreateGoal(userID int64, request model.CreateGoalRequest) (*model.Goal, error) {
	var goal model.Goal
	err := g.db.QueryRow(`
		INSERT INTO goals (user_id, name, description, deadline, status) 
		VALUES ($1, $2, $3, $4, 'active') 
		RETURNING id, user_id, name, description, deadline, created_at, status`,
		userID, request.Name, request.Description, request.Deadline).Scan(
		&goal.ID, &goal.UserID, &goal.Name, &goal.Description, &goal.Deadline, &goal.CreatedAt, &goal.Status)

	if err != nil {
		return nil, err
	}
	return &goal, nil
}

func (g *goalRepository) ReadUserGoals(userID int64) ([]*model.Goal, error) {
	rows, err := g.db.Query("SELECT id, user_id, name, description, deadline, created_at, status FROM goals WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goals []*model.Goal
	for rows.Next() {
		var goal model.Goal
		err := rows.Scan(&goal.ID, &goal.UserID, &goal.Name, &goal.Description, &goal.Deadline, &goal.CreatedAt, &goal.Status)
		if err != nil {
			return nil, err
		}
		goals = append(goals, &goal)
	}
	return goals, nil
}

func (g *goalRepository) UpdateGoal(request model.UpdateGoalRequest) (*model.Goal, error) {
	var goal model.Goal
	err := g.db.QueryRow(`
		UPDATE goals SET name=$2, description=$3, deadline=$4, status=$5
		WHERE id=$1 RETURNING id, user_id, name, description, deadline, created_at, status`,
		request.ID, request.Name, request.Description, request.Deadline, request.Status).Scan(
		&goal.ID, &goal.UserID, &goal.Name, &goal.Description, &goal.Deadline, &goal.CreatedAt, &goal.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("goal not found")
		}
		return nil, err
	}
	return &goal, nil
}

func (g *goalRepository) DeleteGoal(goalID int64) error {
	result, err := g.db.Exec("DELETE FROM goals WHERE id = $1", goalID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("goal not found")
	}
	return nil
}
