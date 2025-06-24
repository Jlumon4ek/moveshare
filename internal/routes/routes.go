package routes

import (
	"database/sql"
	"moveshare/internal/handlers"
	"moveshare/internal/middleware"
	"moveshare/internal/repository"
	"moveshare/internal/services"

	"github.com/gorilla/mux"
)

func NewRouter(db *sql.DB, jwtService services.JWTService) *mux.Router {
	userRepo := repository.NewUserRepository(db)
	authSvc := services.NewAuthService(userRepo)
	authHandler := &handlers.AuthHandler{
		AuthService: authSvc,
		JWTService:  jwtService,
	}

	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)
	r.HandleFunc("/sign-up", authHandler.SignUp).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")
	return r
}
