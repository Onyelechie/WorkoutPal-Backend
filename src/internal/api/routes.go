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

	return r
}

func Routes(r chi.Router) http.Handler {
	// --- Init Repositories ---
	userRepository := repository.NewUserRepository()

	// --- Init Services ---
	userService := service.NewUserService(userRepository)

	// --- Init Handlers ---
	userHandler := handler.NewUserHandler(userService)

	// --- Init Middleware ---
	var idMiddleware = middleware2.IdMiddleware()

	// --- Register Routes ---
	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateNewUser)
		r.Get("/", userHandler.ReadAllUsers)
		r.With(idMiddleware).Patch("/{id}", userHandler.UpdateUser)
		r.With(idMiddleware).Delete("/{id}", userHandler.DeleteUser)
	})

	return r
}
