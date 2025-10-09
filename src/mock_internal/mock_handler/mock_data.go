package mock_handler

import "workoutpal/src/internal/model"

var user = model.User{
	ID:           1,
	Username:     "bigworkoutman",
	Email:        "email@email.email",
	Password:     "supersecure",
	Name:         "Sagev Wen Outfall",
	Age:          57,
	Height:       3.7,
	HeightMetric: "meters",
	Weight:       281,
	WeightMetric: "kg",
	Avatar:       "",
}

var post = model.Post{
	ID:       1,
	PostedBy: "bigworkoutman",
	Title:    "Big Workout",
	Caption:  "Bigger gains",
	Date:     "01/01/0001",
	Content:  "https://some-cdn-link.com",
	Likes:    100320234,
	Comments: []model.Comment{
		{
			CommentedBy: "bigworkoutman",
			Comment:     "sheeeeesh",
			Date:        "01/02/0003",
		},
	},
}

var exercise = model.Exercise{
	ID:                  1,
	Name:                "pushups",
	Targets:             []string{"chest", "shoulder", "arms"},
	Intensity:           "high",
	Expertise:           "beginner",
	Image:               "https://some-cdn-link.com",
	Demo:                "https://some-cdn-link.com",
	RecommendedCount:    5000,
	RecommendedSets:     15,
	RecommendedDuration: 1000,
	Custom:              false,
}

var workout = model.Workout{
	ID:        1,
	Name:      "Big Workout",
	Frequency: "Daily",
	NextRound: "01/02/0004",
	Exercises: []model.RegisteredExercise{
		{
			StartTime: "00:00",
			EndTime:   "01:30",
			Count:     5000,
			Sets:      15,
			Duration:  1000,
			Exercise: model.Exercise{
				ID:        1,
				Name:      "pushups",
				Targets:   []string{"chest", "shoulder", "arms"},
				Intensity: "high",
				Expertise: "beginner",
				Image:     "https://some-cdn-link.com",
				Demo:      "https://some-cdn-link.com",
			},
		},
	},
}

var routine = model.ExerciseRoutine{
	ID:          1,
	UserID:      1,
	Name:        "Morning Routine",
	Description: "Daily morning workout",
}
