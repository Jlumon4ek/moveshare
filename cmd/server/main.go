package main

import (
	"context"
	"log/slog"
	"moveshare/internal/config"
	"moveshare/internal/db"
	"moveshare/internal/routes"
	"moveshare/internal/services"
	"net/http"
	"os"
	"time"

	_ "moveshare/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	config.SetupLogger()
	cfg, err := config.LoadDatabaseSettings()
	if err != nil {
		slog.Error("Failed to load database settings", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.Info("DB config",
		slog.String("user", cfg.User),
		slog.String("pass", cfg.Password),
		slog.String("host", cfg.Host),
		slog.Int("port", cfg.Port),
		slog.String("db", cfg.DB),
	)

	database, err := db.NewPostgresDB(cfg)
	if err != nil {
		slog.Error("Failed to connect to database (sql.Open)", slog.String("error", err.Error()))
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := database.PingContext(ctx); err != nil {
		slog.Error("Database connection check failed (ping)", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("✅ Database connection successful")

	defer database.Close()

	jwtService, err := services.NewJWTService("internal/keys/private.pem")
	if err != nil {
		slog.Error("Failed to load JWT private key", slog.String("error", err.Error()))
		os.Exit(1)
	}

	r := routes.NewRouter(database, jwtService)
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	slog.Info("🌟 Server started", slog.String("address", ":8080"))
	http.ListenAndServe(":8080", r)
}
