package main

import (
	"log"
	"net/http"
	"workoutpal/src/internal/api"
)

func main() {
	r := api.RegisterRoutes()
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", r))
}
