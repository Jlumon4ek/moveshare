package repository

import (
	"database/sql"
	"moveshare/internal/models"
	"time"
)

type UserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	UserExists(email, username string) (bool, error)
	GetUserByEmail(email string) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *models.User) (*models.User, error) {
	query := `
		INSERT INTO users (email, username, password_hash, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	user.CreatedAt = time.Now()

	err := r.db.QueryRow(query, user.Email, user.Username, user.Password, user.CreatedAt).
		Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) UserExists(email, username string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 OR username = $2)`
	var exists bool
	err := r.db.QueryRow(query, email, username).Scan(&exists)
	return exists, err
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, username, password_hash, created_at FROM users WHERE email = $1`
	var user models.User
	err := r.db.QueryRow(query, email).
		Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
