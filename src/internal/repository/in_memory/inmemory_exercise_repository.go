package in_memory

import (
	"errors"
	"workoutpal/src/internal/domain/repository"
	"workoutpal/src/internal/model"
)

type inMemoryExerciseRepository struct {
	data map[int64]*model.Exercise
}

func NewInMemoryExerciseRepository() repository.ExerciseRepository {
	return &inMemoryExerciseRepository{
		data: map[int64]*model.Exercise{
			1: {ID: 1, Name: "Push-ups", Targets: []string{"chest", "shoulders", "triceps"}, Intensity: "medium", Expertise: "beginner"},
			2: {ID: 2, Name: "Pull-ups", Targets: []string{"back", "biceps"}, Intensity: "high", Expertise: "intermediate"},
			3: {ID: 3, Name: "Squats", Targets: []string{"legs", "glutes"}, Intensity: "medium", Expertise: "beginner"},
			4: {ID: 4, Name: "Deadlifts", Targets: []string{"back", "legs", "glutes"}, Intensity: "high", Expertise: "advanced"},
		},
	}
}

func (e *inMemoryExerciseRepository) ReadExerciseByID(id int64) (*model.Exercise, error) {
	if ex, ok := e.data[id]; ok {
		// return a copy to avoid external mutation of our map values
		cp := *ex
		return &cp, nil
	}
	return nil, errors.New("exercise not found")
}

func (e *inMemoryExerciseRepository) ReadAllExercises() ([]*model.Exercise, error) {
	out := make([]*model.Exercise, 0, len(e.data))
	for _, ex := range e.data {
		cp := *ex
		out = append(out, &cp)
	}
	return out, nil
}
