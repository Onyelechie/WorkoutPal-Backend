package main

import (
	"database/sql"
	"log"
	"net/http"
	"workoutpal/src/internal/api"
	"workoutpal/src/internal/config"
)

func main() {
	cfg := config.Load()
	db, err := connectToDatabase(cfg)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return
	}

	r := api.RegisterRoutes(cfg, db)
	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}

func connectToDatabase(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		panic("Failed to connect to PostgreSQL: " + err.Error())
	}
	if err = db.Ping(); err != nil {
		panic("Failed to ping PostgreSQL: " + err.Error())
	}
	return db, nil
}
