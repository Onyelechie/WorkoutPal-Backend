package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=user password=password dbname=workoutpal sslmode=disable")
	if err != nil {
		fmt.Printf("❌ Failed to connect: %v\n", err)
		return
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		fmt.Printf("❌ Failed to ping: %v\n", err)
		return
	}

	fmt.Println("✅ PostgreSQL connection successful!")

	// Test query
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		fmt.Printf("❌ Query failed: %v\n", err)
		return
	}

	fmt.Printf("✅ Found %d users in database\n", count)
}