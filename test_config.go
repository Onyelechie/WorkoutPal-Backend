package main

import (
	"fmt"
	"workoutpal/src/internal/config"
)

func main() {
	cfg := config.Load()
	fmt.Printf("DatabaseURL: %s\n", cfg.DatabaseURL)
}