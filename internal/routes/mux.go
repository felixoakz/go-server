package routes

import (
	"goserve/internal/handlers"
	"net/http"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.HandleRoot)
	mux.HandleFunc("POST /users", handlers.CreateUser)
	mux.HandleFunc("GET /users/{id}", handlers.GetUser)
	mux.HandleFunc("DELETE /users/{id}", handlers.DeleteUser)

	return mux
}
