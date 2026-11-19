package httpx

import (
	"encoding/json"
	"net/http"

	"github.com/Sujon37/school-management-system/internal/handler"
	"github.com/Sujon37/school-management-system/internal/middleware"
)

func NewRouter(jwtSecret string, authHandler *handler.AuthHandler, userHandler *handler.UserHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", healthCheck)
	mux.HandleFunc("POST /api/auth/login", authHandler.LoginHandler)

	mux.Handle(
		"POST /api/users",
		middleware.Auth(jwtSecret)(
			middleware.RequireRole("admin")(
				http.HandlerFunc(userHandler.CreateUserHandler),
			),
		),
	)

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
