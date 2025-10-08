package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"workoutpal/src/internal/config"
	_ "github.com/lib/pq"
)

type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Message  string `json:"message"`
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status: "ok",
	}

	// Test PostgreSQL connection
	cfg := config.Load()
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		response.Database = "in-memory"
		response.Message = "PostgreSQL not available, using fallback"
	} else if err = db.Ping(); err != nil {
		response.Database = "in-memory"
		response.Message = "PostgreSQL not reachable, using fallback"
		db.Close()
	} else {
		response.Database = "postgresql"
		response.Message = "Connected to PostgreSQL"
		db.Close()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}