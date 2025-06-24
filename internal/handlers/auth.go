package handlers

import (
	"encoding/json"
	"log/slog"
	"moveshare/internal/models"
	"moveshare/internal/services"
	"net/http"
)

type AuthHandler struct {
	AuthService services.AuthService
	JWTService  services.JWTService
}

// SignUp godoc
// @Summary Регистрация пользователя
// @Description Создание нового пользователя с email, username и password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param input body models.SignUpRequest true "User registration data"
// @Success 201 {object} models.User
// @Failure 400
// @Failure 409
// @Failure 500
// @Router /sign-up [post]
func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req models.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode request body",
			slog.String("error", err.Error()),
			slog.String("path", r.URL.Path))
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	slog.Info("Creating new user",
		slog.String("email", req.Email),
		slog.String("username", req.Username))

	user, err := h.AuthService.CreateUser(req)
	if err != nil {
		switch err {
		case services.ErrUserExists:
			slog.Warn("User already exists",
				slog.String("email", req.Email),
				slog.String("username", req.Username))
			http.Error(w, "user already exists", http.StatusConflict)
		case services.ErrInvalidInput:
			slog.Warn("Invalid input data",
				slog.String("email", req.Email),
				slog.String("username", req.Username))
			http.Error(w, "invalid input data", http.StatusBadRequest)
		default:
			slog.Error("Failed to create user",
				slog.String("error", err.Error()),
				slog.String("email", req.Email))
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	slog.Info("User created successfully",
		slog.Int("user_id", user.ID),
		slog.String("email", user.Email),
		slog.String("username", user.Username))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Login godoc
// @Summary Авторизация пользователя
// @Description Логин по email и password, возвращает JWT access_token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param input body models.LoginRequest true "Login data"
// @Success 200 {object} models.LoginResponse
// @Failure 400
// @Failure 401
// @Router /login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.AuthService.Authenticate(req)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := h.JWTService.GenerateToken(user.ID, user.Email)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	resp := models.LoginResponse{AccessToken: token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
