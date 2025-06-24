package models

import (
	"time"
)

type User struct {
	ID        int       `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"-" db:"password_hash"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type SignUpRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}
