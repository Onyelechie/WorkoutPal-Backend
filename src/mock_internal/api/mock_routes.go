package api

import (
	"workoutpal/src/mock_internal/mock_handler"

	"github.com/go-chi/chi/v5"
)

func MockRoutes(r chi.Router) {
	// --- Init Mock Handlers ---
	mockAuthHandler := mock_handler.NewMockAuthHandler()
	mockUserHandler := mock_handler.NewMockUserHandler()
	mockRelationshipHandler := mock_handler.NewMockRelationshipHandler()
	mockPostHandler := mock_handler.NewMockPostHandler()
	mockExerciseHandler := mock_handler.NewMockExerciseHandler()
	mockWorkoutHandler := mock_handler.NewMockWorkoutHandler()

	// --- Init Mock Routes ---
	r.Post("/login", mockAuthHandler.Login)
	r.Post("/logout", mockAuthHandler.Logout)

	// Users
	r.Route("/users", func(r chi.Router) {
		r.Get("/", mockUserHandler.ReadAllUsers)
		r.Post("/", mockUserHandler.CreateNewUser)
		r.Patch("/{id}", mockUserHandler.UpdateUser)
		r.Delete("/{id}", mockUserHandler.DeleteUser)
	})

	// Relationships
	r.Route("/relationships", func(r chi.Router) {
		r.Post("/follow", mockRelationshipHandler.FollowUser)
		r.Post("/unfollow", mockRelationshipHandler.UnfollowUser)
	})

	// Followers/Followings
	r.Route("/users/{id}", func(r chi.Router) {
		r.Get("/followers", mockRelationshipHandler.ReadFollowers)
		r.Get("/followings", mockRelationshipHandler.ReadFollowings)
	})

	// Posts
	r.Route("/posts", func(r chi.Router) {
		r.Get("/", mockPostHandler.ReadPosts)
		r.Post("/", mockPostHandler.CreatePost)
		r.Post("/{id}/comment", mockPostHandler.CommentOnPost)
		r.Post("/{id}/like", mockPostHandler.LikePost)
	})

	// Exercises
	r.Route("/exercises", func(r chi.Router) {
		r.Get("/", mockExerciseHandler.ReadExercises)
		r.Post("/", mockExerciseHandler.CreateExercise)
	})

	// Workouts
	r.Route("/workouts", func(r chi.Router) {
		r.Get("/", mockWorkoutHandler.ReadWorkouts)
		r.Post("/", mockWorkoutHandler.CreateWorkout)
		r.Patch("/{id}", mockWorkoutHandler.UpdateWorkout)
	})
}
