package main

import (
	"fmt"
	"goserve/internal/routes"
	"net/http"
)

func main() {
	mux := routes.SetupRoutes()

	fmt.Println("Server running on port :8080")

	http.ListenAndServe(":8080", mux)
}
