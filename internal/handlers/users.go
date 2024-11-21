package handlers

import (
	"bytes"
	"encoding/json"
	"goserve/internal/models"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var (
	userCache  = make(map[int]models.User)
	cacheMutex sync.RWMutex
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
    // Log the request method, URL, and headers
    log.Printf("Received %s request for %s", r.Method, r.URL.Path)
    log.Printf("Request Headers: %v", r.Header)

    // Print the request body (for debugging purposes, be cautious in production)
    var bodyContent []byte
    if r.Body != nil {
        bodyContent, _ = io.ReadAll(r.Body)
        r.Body = io.NopCloser(bytes.NewReader(bodyContent)) // Reset body so it can be read again later
    }
    log.Printf("Request Body: %s", string(bodyContent))

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	cacheMutex.Lock()
	userCache[len(userCache)+1] = user
	cacheMutex.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cacheMutex.RLock()
	user, ok := userCache[id]
	cacheMutex.RUnlock()

	if !ok {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "aplication/json")
	j, err := json.Marshal(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := userCache[id]; !ok {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	cacheMutex.Lock()
	delete(userCache, id)
	cacheMutex.Unlock()

	w.WriteHeader(http.StatusNoContent)
}
