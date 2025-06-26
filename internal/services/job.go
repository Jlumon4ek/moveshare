package services

import (
	"errors"
	"moveshare/internal/models"
	"moveshare/internal/repository"

	"github.com/google/uuid"
)

var ErrJobNotFound = errors.New("job not found")

type JobService interface {
	CreateJob(req models.CreateJobRequest) (*models.Job, error)
	GetJobs(filter models.JobFilter, limit, offset int) ([]*models.Job, int, error)
	DeleteJob(id string) error
}

type jobService struct {
	repo repository.JobRepository
}

func NewJobService(repo repository.JobRepository) JobService {
	return &jobService{repo: repo}
}

func (s *jobService) CreateJob(req models.CreateJobRequest) (*models.Job, error) {
	job := &models.Job{
		ID:                            uuid.New().String(),
		JobTitle:                      req.JobTitle,
		NumberOfBedrooms:              req.NumberOfBedrooms,
		AdditionalServices:            req.AdditionalServices,
		DescriptionAdditionalServices: req.DescriptionAdditionalServices,
		TruckSize:                     req.TruckSize,
		PickupDateTime:                req.PickupDateTime,
		DeliveryDateTime:              req.DeliveryDateTime,
		CutAmount:                     req.CutAmount,
		PaymentAmount:                 req.PaymentAmount,
	}
	return s.repo.CreateJob(job)
}

func (s *jobService) GetJobs(filter models.JobFilter, limit, offset int) ([]*models.Job, int, error) {
	return s.repo.GetJobs(filter, limit, offset)
}

func (s *jobService) DeleteJob(id string) error {
	err := s.repo.DeleteJob(id)
	if err != nil {
		if err.Error() == "job not found" {
			return ErrJobNotFound
		}
		return err
	}
	return nil
}
