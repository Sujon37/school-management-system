package httpx

import (
	"encoding/json"
	"net/http"

	"github.com/Sujon37/school-management-system/internal/handler"
	"github.com/Sujon37/school-management-system/internal/middleware"
	"github.com/Sujon37/school-management-system/internal/response"
)

func NewRouter(jwtSecret string, authHandler *handler.AuthHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", healthCheck)
	mux.HandleFunc("POST /api/auth/login", authHandler.LoginHandler)

	mux.Handle(
		"GET /api/protected",
		middleware.Auth(jwtSecret)(
			middleware.RequireRole("admin")(
				http.HandlerFunc(protectedHandler),
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

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "claims not found")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{
		"user_id": claims.Subject,
        "role":    claims.Role,
	})
}
