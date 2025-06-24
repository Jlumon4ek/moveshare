package db

import (
	"context"
	"database/sql"
	"fmt"
	"moveshare/internal/models"
	"strings"

	"github.com/jmoiron/sqlx"
)

type JobRepository struct {
	db *sqlx.DB
}

func NewJobRepository(db *sqlx.DB) *JobRepository {
	return &JobRepository{db: db}
}

func (r *JobRepository) CreateJob(ctx context.Context, job *models.Job) error {
	query := `
		INSERT INTO jobs (title, user_id, origin, destination, distance, 
						  start_date, end_date, truck_size, weight, volume, 
						  payout, is_new, is_verified, is_protected, is_escrow)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(ctx, query,
		job.Title, job.UserID, job.Origin, job.Destination, job.Distance,
		job.StartDate, job.EndDate, job.TruckSize, job.Weight, job.Volume,
		job.Payout, job.IsNew, job.IsVerified, job.IsProtected, job.IsEscrow,
	).Scan(&job.ID, &job.CreatedAt)

	return err
}

// GetJobByID получает задание по ID
func (r *JobRepository) GetJobByID(ctx context.Context, id int64) (*models.Job, error) {
	job := &models.Job{}
	query := `
		SELECT * FROM jobs
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, job, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return job, nil
}

// GetJobsByUserID получает задания, созданные пользователем
func (r *JobRepository) GetJobsByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*models.Job, int, error) {
	jobs := []*models.Job{}
	offset := (page - 1) * pageSize

	query := `
		SELECT * FROM jobs
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	err := r.db.SelectContext(ctx, &jobs, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	// Получаем общее количество для пагинации
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM jobs WHERE user_id = $1`
	err = r.db.GetContext(ctx, &totalCount, countQuery, userID)
	if err != nil {
		return nil, 0, err
	}

	return jobs, totalCount, nil
}

// GetAvailableJobs получает доступные задания с фильтрацией
func (r *JobRepository) GetAvailableJobs(ctx context.Context, filter *models.JobFilter) ([]*models.Job, int, error) {
	jobs := []*models.Job{}

	// Строим SQL запрос с фильтрацией
	queryParams := []interface{}{}
	conditions := []string{"is_claimed = FALSE"}
	paramIndex := 1

	if filter.Query != "" {
		conditions = append(conditions, fmt.Sprintf("(title ILIKE $%d OR origin ILIKE $%d OR destination ILIKE $%d)",
			paramIndex, paramIndex, paramIndex))
		queryParams = append(queryParams, "%"+filter.Query+"%")
		paramIndex++
	}

	if filter.Origin != "" {
		conditions = append(conditions, fmt.Sprintf("origin ILIKE $%d", paramIndex))
		queryParams = append(queryParams, "%"+filter.Origin+"%")
		paramIndex++
	}

	if filter.Destination != "" {
		conditions = append(conditions, fmt.Sprintf("destination ILIKE $%d", paramIndex))
		queryParams = append(queryParams, "%"+filter.Destination+"%")
		paramIndex++
	}

	if filter.MinDistance > 0 {
		conditions = append(conditions, fmt.Sprintf("distance >= $%d", paramIndex))
		queryParams = append(queryParams, filter.MinDistance)
		paramIndex++
	}

	if filter.MaxDistance > 0 {
		conditions = append(conditions, fmt.Sprintf("distance <= $%d", paramIndex))
		queryParams = append(queryParams, filter.MaxDistance)
		paramIndex++
	}

	if !filter.StartDate.IsZero() {
		conditions = append(conditions, fmt.Sprintf("start_date >= $%d", paramIndex))
		queryParams = append(queryParams, filter.StartDate)
		paramIndex++
	}

	if !filter.EndDate.IsZero() {
		conditions = append(conditions, fmt.Sprintf("end_date <= $%d", paramIndex))
		queryParams = append(queryParams, filter.EndDate)
		paramIndex++
	}

	if len(filter.TruckSizes) > 0 {
		placeholders := make([]string, len(filter.TruckSizes))
		for i := range filter.TruckSizes {
			placeholders[i] = fmt.Sprintf("$%d", paramIndex)
			queryParams = append(queryParams, filter.TruckSizes[i])
			paramIndex++
		}
		conditions = append(conditions, fmt.Sprintf("truck_size IN (%s)", strings.Join(placeholders, ", ")))
	}

	if filter.MinPayout > 0 {
		conditions = append(conditions, fmt.Sprintf("payout >= $%d", paramIndex))
		queryParams = append(queryParams, filter.MinPayout)
		paramIndex++
	}

	if filter.MaxPayout > 0 {
		conditions = append(conditions, fmt.Sprintf("payout <= $%d", paramIndex))
		queryParams = append(queryParams, filter.MaxPayout)
		paramIndex++
	}

	// Формируем полный SQL запрос
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT * FROM jobs
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, paramIndex, paramIndex+1)

	// Добавляем параметры для пагинации
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 10 // По умолчанию 10 элементов на странице
	}

	page := filter.Page
	if page <= 0 {
		page = 1 // По умолчанию первая страница
	}

	offset := (page - 1) * pageSize
	queryParams = append(queryParams, pageSize, offset)

	err := r.db.SelectContext(ctx, &jobs, query, queryParams...)
	if err != nil {
		return nil, 0, err
	}

	// Получаем общее количество для пагинации
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM jobs %s", whereClause)
	var totalCount int
	err = r.db.GetContext(ctx, &totalCount, countQuery, queryParams[:len(queryParams)-2]...)
	if err != nil {
		return nil, 0, err
	}

	return jobs, totalCount, nil
}

// ClaimJob создает отклик на задание
func (r *JobRepository) ClaimJob(ctx context.Context, claim *models.JobClaim) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Проверяем, что задание существует и не занято
	var isClaimed bool
	checkQuery := "SELECT is_claimed FROM jobs WHERE id = $1"
	err = tx.GetContext(ctx, &isClaimed, checkQuery, claim.JobID)
	if err != nil {
		return err
	}

	if isClaimed {
		return fmt.Errorf("job is already claimed")
	}

	// Создаем отклик
	insertQuery := `
		INSERT INTO job_claims (job_id, user_id, status)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	err = tx.QueryRowContext(ctx, insertQuery, claim.JobID, claim.UserID, claim.Status).Scan(&claim.ID, &claim.CreatedAt)
	if err != nil {
		return err
	}

	// Обновляем статус задания как занятого
	updateQuery := "UPDATE jobs SET is_claimed = TRUE WHERE id = $1"
	_, err = tx.ExecContext(ctx, updateQuery, claim.JobID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetClaimedJobsByUserID получает задания, на которые откликнулся пользователь
func (r *JobRepository) GetClaimedJobsByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*models.Job, int, error) {
	offset := (page - 1) * pageSize
	jobs := []*models.Job{}

	query := `
		SELECT j.* FROM jobs j
		INNER JOIN job_claims jc ON j.id = jc.job_id
		WHERE jc.user_id = $1
		ORDER BY jc.created_at DESC
		LIMIT $2 OFFSET $3
	`

	err := r.db.SelectContext(ctx, &jobs, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	// Получаем общее количество для пагинации
	var totalCount int
	countQuery := `
		SELECT COUNT(*) FROM jobs j
		INNER JOIN job_claims jc ON j.id = jc.job_id
		WHERE jc.user_id = $1
	`
	err = r.db.GetContext(ctx, &totalCount, countQuery, userID)
	if err != nil {
		return nil, 0, err
	}

	return jobs, totalCount, nil
}

// CancelJobClaim отменяет отклик на задание
func (r *JobRepository) CancelJobClaim(ctx context.Context, jobID, userID int64) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Удаляем отклик
	deleteQuery := "DELETE FROM job_claims WHERE job_id = $1 AND user_id = $2"
	result, err := tx.ExecContext(ctx, deleteQuery, jobID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("job claim not found")
	}

	// Обновляем статус задания
	updateQuery := "UPDATE jobs SET is_claimed = FALSE WHERE id = $1"
	_, err = tx.ExecContext(ctx, updateQuery, jobID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
