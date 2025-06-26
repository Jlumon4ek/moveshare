package models

import "time"

type NumberOfBedrooms string

const (
	OneBedroom    NumberOfBedrooms = "1"
	TwoBedrooms   NumberOfBedrooms = "2"
	ThreeBedrooms NumberOfBedrooms = "3"
	FourBedrooms  NumberOfBedrooms = "4"
	FivePlus      NumberOfBedrooms = "5+"
	OfficeBedroom NumberOfBedrooms = "office"
)

type TruckSize string

const (
	SmallTruck  TruckSize = "small"
	MediumTruck TruckSize = "medium"
	LargeTruck  TruckSize = "large"
)

type Job struct {
	ID                            string           `json:"id" db:"id"`
	JobTitle                      string           `json:"title" db:"title"`
	NumberOfBedrooms              NumberOfBedrooms `json:"number_of_bedrooms" db:"number_of_bedrooms"`
	AdditionalServices            string           `json:"additional_services" db:"additional_services"`
	DescriptionAdditionalServices string           `json:"description_additional_services" db:"description_additional_services"`
	TruckSize                     TruckSize        `json:"truck_size" db:"truck_size"`
	PickupDateTime                time.Time        `json:"pickup_datetime" db:"pickup_datetime"`
	DeliveryDateTime              time.Time        `json:"delivery_datetime" db:"delivery_datetime"`
	CutAmount                     float64          `json:"cut_amount" db:"cut_amount"`
	PaymentAmount                 float64          `json:"payment_amount" db:"payment_amount"`
}

// CreateJobRequest используется для создания новой Job через API (без ID)
type CreateJobRequest struct {
	JobTitle                      string           `json:"title"`
	NumberOfBedrooms              NumberOfBedrooms `json:"number_of_bedrooms"`
	AdditionalServices            string           `json:"additional_services"`
	DescriptionAdditionalServices string           `json:"description_additional_services"`
	TruckSize                     TruckSize        `json:"truck_size"`
	PickupDateTime                time.Time        `json:"pickup_datetime"`
	DeliveryDateTime              time.Time        `json:"delivery_datetime"`
	CutAmount                     float64          `json:"cut_amount"`
	PaymentAmount                 float64          `json:"payment_amount"`
}

// JobFilter для фильтрации и поиска
type JobFilter struct {
	NumberOfBedrooms string     // "1", "2", "office" и т.д.
	DateStart        *time.Time // >=
	DateEnd          *time.Time // <=
	TruckSize        string     // "small", "medium", "large"
	PayoutMin        *float64   // >=
	PayoutMax        *float64   // <=
}

// JobListResponse для ответа на GET /jobs
type JobListResponse struct {
	Jobs  []*Job `json:"jobs"`
	Total int    `json:"total"`
}
