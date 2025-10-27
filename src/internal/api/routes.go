package api

import (
	"database/sql"
	"net/http"
	"workoutpal/src/internal/api/docs"
	"workoutpal/src/internal/config"
	"workoutpal/src/internal/dependency"
	"workoutpal/src/internal/handler"
	middleware2 "workoutpal/src/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterRoutes(cfg *config.Config, db *sql.DB) http.Handler {
	r := chi.NewRouter()

	// --- Global middleware ---
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// --- Real Routes ---
	r.Route("/", func(r chi.Router) {
		appDep := dependency.NewAppDependencies(db)
		Routes(r, appDep, []byte(cfg.JWTSecret))
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

func Routes(r chi.Router, appDep dependency.AppDependencies, secret []byte) http.Handler {
	// --- Init Handlers ---
	userHandler := handler.NewUserHandler(appDep.UserService)
	goalHandler := handler.NewGoalHandler(appDep.GoalService)
	relationshipHandler := handler.NewRelationshipHandler(appDep.RelationshipService)
	routineHandler := handler.NewRoutineHandler(appDep.RoutineService)
	exerciseHandler := handler.NewExerciseHandler(appDep.ExerciseService)
	authHandler := handler.NewAuthHandler(appDep.UserService, appDep.AuthService, secret)
	scheduleHandler := handler.NewScheduleHandler(appDep.ScheduleService)

	// --- Init Middleware ---
	var idMiddleware = middleware2.IdMiddleware()
	var authMiddleware = middleware2.AuthMiddleware(secret)

	// Health check
	r.Get("/health", handler.HealthCheck)

	// Auth routes
	r.Post("/login", authHandler.Login)
	r.Post("/logout", authHandler.Logout)
	r.With(authMiddleware).Get("/me", authHandler.Me)

	// --- Register Routes ---
	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateNewUser)

		r.With(authMiddleware).Group(func(r chi.Router) {
			r.Get("/", userHandler.ReadAllUsers)
			r.With(idMiddleware).Get("/{id}", userHandler.ReadUserByID)
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
			r.With(idMiddleware).Post("/{id}/routines", routineHandler.CreateUserRoutine)
			r.With(idMiddleware).Get("/{id}/routines", routineHandler.ReadUserRoutines)
			r.With(idMiddleware).Delete("/{id}/routines/{routine_id}", routineHandler.DeleteUserRoutine)
		})
	})

	// Exercises
	r.With(authMiddleware).Route("/exercises", func(r chi.Router) {
		r.Get("/", exerciseHandler.ReadExercises)
		r.With(idMiddleware).Get("/{id}", exerciseHandler.ReadExerciseByID)
	})

	// Routines
	r.With(authMiddleware).Route("/routines", func(r chi.Router) {
		r.With(idMiddleware).Get("/{id}", routineHandler.ReadRoutineWithExercises)
		r.With(idMiddleware).Delete("/{id}", routineHandler.DeleteRoutine)
		r.With(idMiddleware).Post("/{id}/exercises", routineHandler.AddExerciseToRoutine)
		r.With(idMiddleware).Delete("/{id}/exercises/{exercise_id}", routineHandler.RemoveExerciseFromRoutine)
	})

	// Schedules
	r.With(authMiddleware).Route("/schedules", func(r chi.Router) {
		r.Get("/", scheduleHandler.ReadUserSchedules)
		r.Get("/{dayOfWeek}", scheduleHandler.ReadUserSchedulesByDay)
		r.Post("/", scheduleHandler.CreateSchedule)
		r.With(idMiddleware).Get("/{id}", scheduleHandler.ReadScheduleByID)
		r.With(idMiddleware).Put("/{id}", scheduleHandler.UpdateSchedule)
		r.With(idMiddleware).Delete("/{id}", scheduleHandler.DeleteSchedule)
	})

	return r
}
