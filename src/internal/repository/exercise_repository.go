package repository

import (
	"database/sql"
	"strings"
	"workoutpal/src/internal/config"
	"workoutpal/src/internal/domain/repository"
	"workoutpal/src/internal/model"

	_ "github.com/lib/pq"
)

type exerciseRepository struct {
	db *sql.DB
}

func NewExerciseRepository() repository.ExerciseRepository {
	cfg := config.Load()
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return NewInMemoryExerciseRepository()
	}
	if err = db.Ping(); err != nil {
		return NewInMemoryExerciseRepository()
	}
	return &exerciseRepository{db: db}
}

func (e *exerciseRepository) GetAllExercises() ([]model.Exercise, error) {
	rows, err := e.db.Query("SELECT id, name, description, targets, image FROM exercises")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []model.Exercise
	for rows.Next() {
		var exercise model.Exercise
		var targetsStr string
		var image, demo sql.NullString
		
		err := rows.Scan(&exercise.ID, &exercise.Name, &exercise.Description, &targetsStr, &image)
		if err != nil {
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
		
		exercises = append(exercises, exercise)
	}
	return exercises, nil
}

// In-memory fallback
type inMemoryExerciseRepository struct{}

func NewInMemoryExerciseRepository() repository.ExerciseRepository {
	return &inMemoryExerciseRepository{}
}

func (e *inMemoryExerciseRepository) GetAllExercises() ([]model.Exercise, error) {
	return []model.Exercise{
		{ID: 1, Name: "Push-ups", Targets: []string{"chest", "shoulders", "triceps"}, Intensity: "medium", Expertise: "beginner"},
		{ID: 2, Name: "Pull-ups", Targets: []string{"back", "biceps"}, Intensity: "high", Expertise: "intermediate"},
		{ID: 3, Name: "Squats", Targets: []string{"legs", "glutes"}, Intensity: "medium", Expertise: "beginner"},
		{ID: 4, Name: "Deadlifts", Targets: []string{"back", "legs", "glutes"}, Intensity: "high", Expertise: "advanced"},
	}, nil
}