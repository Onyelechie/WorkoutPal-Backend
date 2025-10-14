package api

import (
	"net/http"
	"workoutpal/src/internal/api/docs"
	"workoutpal/src/internal/handler"
	middleware2 "workoutpal/src/internal/middleware"
	"workoutpal/src/internal/repository"
	"workoutpal/src/internal/service"
	mock_api "workoutpal/src/mock_internal/api"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	// --- Global middleware ---
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// --- Mock Routes ---
	r.Route("/mock", func(r chi.Router) {
		mock_api.MockRoutes(r)
	})

	// --- Real Routes ---
	r.Route("/", func(r chi.Router) {
		Routes(r)
	})

	// Swagger Docs
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.InstanceName(docs.SwaggerInfo.InstanceName()),
	))

	// CORS middleware
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{
			// TODO read this from config
			"http://localhost:4200",
			"http://localhost:5173",
		},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
	}).Handler(r)

	return corsHandler
}

func Routes(r chi.Router) http.Handler {
	// --- Init Repositories ---
	userRepository := repository.NewUserRepository()
	exerciseRepository := repository.NewExerciseRepository()

	// --- Init Services ---
	userService := service.NewUserService(userRepository)
	exerciseService := service.NewExerciseService(exerciseRepository)
	authService := service.NewAuthService(userRepository)

	// --- Init Handlers ---
	userHandler := handler.NewUserHandler(userService)
	goalHandler := handler.NewGoalHandler(userService)
	relationshipHandler := handler.NewRelationshipHandler(userService)
	workoutHandler := handler.NewWorkoutHandler(userService)
	exerciseHandler := handler.NewExerciseHandler(exerciseService)
	authHandler := handler.NewAuthHandler(userService, authService)

	// --- Init Middleware ---
	var idMiddleware = middleware2.IdMiddleware()

	// Health check
	r.Get("/health", handler.HealthCheck)

	// Auth routes
	r.Post("/login", authHandler.Login)
	r.Post("/logout", authHandler.Logout)
	r.Route("/auth", func(r chi.Router) {
		//r.Post("/google", authHandler.GoogleAuth)
	})

	// --- Register Routes ---
	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateNewUser)

		r.With(middleware2.AuthMiddleware([]byte("secret"))).Group(func(r chi.Router) {
			r.Get("/", userHandler.ReadAllUsers)
			r.With(idMiddleware).Get("/{id}", userHandler.GetUserByID)
			r.With(idMiddleware).Patch("/{id}", userHandler.UpdateUser)
			r.With(idMiddleware).Delete("/{id}", userHandler.DeleteUser)
			// User Goals
			r.With(idMiddleware).Post("/{id}/goals", goalHandler.CreateUserGoal)
			r.With(idMiddleware).Get("/{id}/goals", goalHandler.GetUserGoals)
			// User Followers
			r.With(idMiddleware).Post("/{id}/follow", relationshipHandler.FollowUser)
			r.With(idMiddleware).Post("/{id}/unfollow", relationshipHandler.UnfollowUser)
			r.With(idMiddleware).Get("/{id}/followers", relationshipHandler.ReadFollowers)
			r.With(idMiddleware).Get("/{id}/following", relationshipHandler.ReadFollowings)
			// User Routines
			r.With(idMiddleware).Post("/{id}/routines", workoutHandler.CreateUserRoutine)
			r.With(idMiddleware).Get("/{id}/routines", workoutHandler.GetUserRoutines)
			r.With(idMiddleware).Delete("/{id}/routines/{routine_id}", workoutHandler.DeleteUserRoutine)
		})
	})

	// Exercises
	r.Route("/exercises", func(r chi.Router) {
		r.Get("/", exerciseHandler.ReadExercises)
	})

	// Routines
	r.Route("/routines", func(r chi.Router) {
		r.With(idMiddleware).Get("/{id}", workoutHandler.GetRoutineWithExercises)
		r.With(idMiddleware).Delete("/{id}", workoutHandler.DeleteRoutine)
		r.With(idMiddleware).Post("/{id}/exercises", workoutHandler.AddExerciseToRoutine)
		r.With(idMiddleware).Delete("/{id}/exercises/{exercise_id}", workoutHandler.RemoveExerciseFromRoutine)
	})

	return r
}
