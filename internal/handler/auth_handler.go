package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Sujon37/school-management-system/internal/auth"
	"github.com/Sujon37/school-management-system/internal/response"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authService *auth.AuthService
	jwtSecret   string
}

func NewAuthHandler(authService *auth.AuthService, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		jwtSecret:   jwtSecret,
	}
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req loginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	validate := validator.New()

	if err := validate.Struct(req); err != nil {
		response.WriteError(
			w,
			http.StatusBadRequest,
			"invalid email or password format",
		)
		return
	}

	user, err := h.authService.Login(r.Context(), req.Email, req.Password)

	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) ||
			errors.Is(err, auth.ErrInactiveAccount) {

			response.WriteError(w, http.StatusUnauthorized, "invalid email or password")
			return
		}

		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	token, err := auth.GenerateToken(user, h.jwtSecret)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "could not generate token")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{
			"status":     "success",
			"token":      token,
			"token_type": "Bearer",
		},
	)

}
