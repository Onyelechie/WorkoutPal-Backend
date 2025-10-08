package main

import (
	"log"
	"net/http"
	"workoutpal/src/internal/api"
	"workoutpal/src/internal/config"
)

func main() {
	cfg := config.Load()
	r := api.RegisterRoutes()
	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
