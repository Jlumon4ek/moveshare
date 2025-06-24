package services

import (
	"context"

	"moveshare/internal/db"
	"moveshare/internal/models"
)

type JobService struct {
	jobRepo *db.JobRepository
}

func NewJobService(jobRepo *db.JobRepository) *JobService {
	return &JobService{jobRepo: jobRepo}
}

// CreateJob создает новое задание
func (s *JobService) CreateJob(ctx context.Context, job *models.Job) error {
	// Устанавливаем значения по умолчанию
	job.IsNew = true
	job.IsClaimed = false
	job.IsVerified = false // Можно установить true, если есть система верификации
	job.IsProtected = true
	job.IsEscrow = false

	return s.jobRepo.CreateJob(ctx, job)
}

// GetJobByID получает задание по ID
func (s *JobService) GetJobByID(ctx context.Context, id int64) (*models.Job, error) {
	return s.jobRepo.GetJobByID(ctx, id)
}

// GetJobsByUserID получает задания, созданные пользователем
func (s *JobService) GetJobsByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*models.Job, int, error) {
	return s.jobRepo.GetJobsByUserID(ctx, userID, page, pageSize)
}

// GetAvailableJobs получает доступные задания с фильтрацией
func (s *JobService) GetAvailableJobs(ctx context.Context, filter *models.JobFilter) ([]*models.Job, int, error) {
	return s.jobRepo.GetAvailableJobs(ctx, filter)
}

// ClaimJob создает отклик на задание
func (s *JobService) ClaimJob(ctx context.Context, jobID, userID int64) error {
	claim := &models.JobClaim{
		JobID:  jobID,
		UserID: userID,
		Status: "pending", // По умолчанию статус "в ожидании"
	}
	return s.jobRepo.ClaimJob(ctx, claim)
}

// GetClaimedJobsByUserID получает задания, на которые откликнулся пользователь
func (s *JobService) GetClaimedJobsByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*models.Job, int, error) {
	return s.jobRepo.GetClaimedJobsByUserID(ctx, userID, page, pageSize)
}

// CancelJobClaim отменяет отклик на задание
func (s *JobService) CancelJobClaim(ctx context.Context, jobID, userID int64) error {
	return s.jobRepo.CancelJobClaim(ctx, jobID, userID)
}
