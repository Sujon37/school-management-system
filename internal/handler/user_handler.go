package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Sujon37/school-management-system/internal/response"
	"github.com/Sujon37/school-management-system/internal/user"
	"github.com/go-playground/validator"
)

type UserHandler struct {
	userService *user.UserService
}

func NewUserHandler(userService *user.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

type createUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required,oneof=teacher student"`
}

func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request fields")
		return
	}

	input := user.CreateUserInput{
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	createUser, err := h.userService.CreateUser(r.Context(), input)
	if err != nil {

		if errors.Is(err, user.ErrEmailExists) {
				
			response.WriteError(w, http.StatusConflict, err.Error())
			return
		}

		if errors.Is(err, user.ErrInvalidRole) || errors.Is(err, user.ErrInvalidInput) {
				
			response.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		response.WriteError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	response.WriteJSON(w, http.StatusCreated, map[string]string{
		"id":    createUser.ID,
		"email": createUser.Email,
		"role":  createUser.Role,
	})
}
