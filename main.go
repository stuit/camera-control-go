package main

import (
	"fmt"
	"log"
	"net/http"

	"camera-control/api"
)

func main() {
	mux := http.NewServeMux()
	api.RegisterRoutes(mux)

	fmt.Println("Server running at :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
