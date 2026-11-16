package httpx

import (
	"encoding/json"
	"net/http"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	
	mux.HandleFunc("GET /api/health", healthCheck)
	return mux
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"status":  "OK",
		"message": "School Management System is running",
	}

	json.NewEncoder(w).Encode(response)
}
