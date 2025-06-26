package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"moveshare/internal/models"
	"strings"
)

type JobRepository interface {
	CreateJob(job *models.Job) (*models.Job, error)
	GetJobs(filter models.JobFilter, limit, offset int) ([]*models.Job, int, error)
	DeleteJob(id string) error
}

type jobRepository struct {
	db *sql.DB
}

func NewJobRepository(db *sql.DB) JobRepository {
	return &jobRepository{db: db}
}

func (r *jobRepository) CreateJob(job *models.Job) (*models.Job, error) {
	_, err := r.db.Exec(
		`INSERT INTO jobs 
(id, title, number_of_bedrooms, additional_services, description_additional_services, truck_size, pickup_datetime, delivery_datetime, cut_amount, payment_amount)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		job.ID, job.JobTitle, job.NumberOfBedrooms, job.AdditionalServices, job.DescriptionAdditionalServices,
		job.TruckSize, job.PickupDateTime, job.DeliveryDateTime, job.CutAmount, job.PaymentAmount,
	)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (r *jobRepository) GetJobs(filter models.JobFilter, limit, offset int) ([]*models.Job, int, error) {
	var (
		where  []string
		args   []interface{}
		argIdx = 1
	)

	if filter.NumberOfBedrooms != "" {
		where = append(where, fmt.Sprintf("number_of_bedrooms = $%d", argIdx))
		args = append(args, filter.NumberOfBedrooms)
		argIdx++
	}
	if filter.TruckSize != "" {
		where = append(where, fmt.Sprintf("truck_size = $%d", argIdx))
		args = append(args, filter.TruckSize)
		argIdx++
	}
	if filter.DateStart != nil {
		where = append(where, fmt.Sprintf("pickup_datetime >= $%d", argIdx))
		args = append(args, filter.DateStart)
		argIdx++
	}
	if filter.DateEnd != nil {
		where = append(where, fmt.Sprintf("delivery_datetime <= $%d", argIdx))
		args = append(args, filter.DateEnd)
		argIdx++
	}
	if filter.PayoutMin != nil {
		where = append(where, fmt.Sprintf("payment_amount >= $%d", argIdx))
		args = append(args, filter.PayoutMin)
		argIdx++
	}
	if filter.PayoutMax != nil {
		where = append(where, fmt.Sprintf("payment_amount <= $%d", argIdx))
		args = append(args, filter.PayoutMax)
		argIdx++
	}

	whereClause := ""
	if len(where) > 0 {
		whereClause = "WHERE " + strings.Join(where, " AND ")
	}

	// Основной запрос
	query := fmt.Sprintf(`SELECT id, title, number_of_bedrooms, additional_services, description_additional_services, truck_size, pickup_datetime, delivery_datetime, cut_amount, payment_amount
FROM jobs %s ORDER BY pickup_datetime DESC LIMIT $%d OFFSET $%d`, whereClause, argIdx, argIdx+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var jobs []*models.Job
	for rows.Next() {
		var job models.Job
		err := rows.Scan(
			&job.ID,
			&job.JobTitle,
			&job.NumberOfBedrooms,
			&job.AdditionalServices,
			&job.DescriptionAdditionalServices,
			&job.TruckSize,
			&job.PickupDateTime,
			&job.DeliveryDateTime,
			&job.CutAmount,
			&job.PaymentAmount,
		)
		if err != nil {
			return nil, 0, err
		}
		jobs = append(jobs, &job)
	}

	// Считаем total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM jobs %s", whereClause)
	var total int
	if err := r.db.QueryRow(countQuery, args[:argIdx-1]...).Scan(&total); err != nil {
		return nil, 0, err
	}

	return jobs, total, nil
}

func (r *jobRepository) DeleteJob(id string) error {
	res, err := r.db.Exec("DELETE FROM jobs WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("job not found")
	}
	return nil
}
