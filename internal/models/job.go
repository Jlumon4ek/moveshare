package models

import (
	"time"
)

type TruckSize string

const (
	Small  TruckSize = "small"
	Medium TruckSize = "medium"
	Large  TruckSize = "large"
)

type Job struct {
	ID          int64     `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UserID      int64     `json:"user_id" db:"user_id"`
	Origin      string    `json:"origin" db:"origin"`
	Destination string    `json:"destination" db:"destination"`
	Distance    float64   `json:"distance" db:"distance"`
	StartDate   time.Time `json:"start_date" db:"start_date"`
	EndDate     time.Time `json:"end_date" db:"end_date"`
	TruckSize   TruckSize `json:"truck_size" db:"truck_size"`
	Weight      float64   `json:"weight" db:"weight"` // в lbs
	Volume      float64   `json:"volume" db:"volume"` // в ft³
	Payout      float64   `json:"payout" db:"payout"` // в USD
	IsNew       bool      `json:"is_new" db:"is_new"`
	IsClaimed   bool      `json:"is_claimed" db:"is_claimed"`
	IsVerified  bool      `json:"is_verified" db:"is_verified"`
	IsProtected bool      `json:"is_protected" db:"is_protected"`
	IsEscrow    bool      `json:"is_escrow" db:"is_escrow"`
}

type JobClaim struct {
	ID        int64     `json:"id" db:"id"`
	JobID     int64     `json:"job_id" db:"job_id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Status    string    `json:"status" db:"status"`
}

type JobFilter struct {
	Query       string    `json:"query"`
	Origin      string    `json:"origin"`
	Destination string    `json:"destination"`
	MinDistance float64   `json:"min_distance"`
	MaxDistance float64   `json:"max_distance"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	TruckSizes  []string  `json:"truck_sizes"`
	MinPayout   float64   `json:"min_payout"`
	MaxPayout   float64   `json:"max_payout"`
	Page        int       `json:"page"`
	PageSize    int       `json:"page_size"`
}
