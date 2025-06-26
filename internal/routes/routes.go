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

	jobRepo := repository.NewJobRepository(db)
	jobService := services.NewJobService(jobRepo)
	jobHandler := handlers.NewJobHandler(jobService)

	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)

	r.HandleFunc("/sign-up", authHandler.SignUp).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	jobs := r.PathPrefix("/jobs").Subrouter()
	jobs.Use(middleware.AuthMiddleware(jwtService))
	jobs.HandleFunc("", jobHandler.CreateJob).Methods("POST")
	jobs.HandleFunc("", jobHandler.GetJobs).Methods("GET")
	jobs.HandleFunc("/{id}", jobHandler.DeleteJob).Methods("DELETE")

	return r
}
