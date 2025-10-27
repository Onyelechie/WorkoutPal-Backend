package config

import (
	"fmt"
)

func main() {
	cfg := Load()
	fmt.Printf("DatabaseURL: %s\n", cfg.DatabaseURL)
}