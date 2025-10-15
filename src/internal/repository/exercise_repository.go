package repository

import (
	"database/sql"
	"errors"
	"strings"
	"workoutpal/src/internal/domain/repository"
	"workoutpal/src/internal/model"

	_ "github.com/lib/pq"
)

type exerciseRepository struct {
	db *sql.DB
}

func NewExerciseRepository(db *sql.DB) repository.ExerciseRepository {
	return &exerciseRepository{db: db}
}
func (e *exerciseRepository) ReadExerciseByID(id int64) (*model.Exercise, error) {
	var exercise model.Exercise
	var targetsStr string
	var image, demo sql.NullString

	err := e.db.QueryRow(
		"SELECT id, name, description, targets, image, demo FROM exercises WHERE id = $1",
		id,
	).Scan(&exercise.ID, &exercise.Name, &exercise.Description, &targetsStr, &image, &demo)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("exercise not found")
		}
		return nil, err
	}

	if targetsStr != "" {
		exercise.Targets = strings.Split(targetsStr, ",")
	}
	if image.Valid {
		exercise.Image = image.String
	}
	if demo.Valid {
		exercise.Demo = demo.String
	}
	return &exercise, nil
}

func (e *exerciseRepository) ReadAllExercises() ([]*model.Exercise, error) {
	rows, err := e.db.Query("SELECT id, name, description, targets, image, demo FROM exercises")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []*model.Exercise
	for rows.Next() {
		var exercise model.Exercise
		var targetsStr string
		var image, demo sql.NullString

		if err := rows.Scan(&exercise.ID, &exercise.Name, &exercise.Description, &targetsStr, &image, &demo); err != nil {
			return nil, err
		}
		if targetsStr != "" {
			exercise.Targets = strings.Split(targetsStr, ",")
		}
		if image.Valid {
			exercise.Image = image.String
		}
		if demo.Valid {
			exercise.Demo = demo.String
		}
		exercises = append(exercises, &exercise)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return exercises, nil
}
