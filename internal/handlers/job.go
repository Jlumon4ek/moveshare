package handlers

import (
	"encoding/json"
	"moveshare/internal/models"
	"moveshare/internal/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// JobHandler отвечает за обработку job-related endpoints
type JobHandler struct {
	JobService services.JobService
}

func NewJobHandler(jobService services.JobService) *JobHandler {
	return &JobHandler{JobService: jobService}
}

// CreateJob godoc
// @Summary Создание новой работы (Job)
// @Description Создать новую работу (Job) с параметрами перевозки
// @Tags jobs
// @Accept  json
// @Produce  json
// @Param input body models.CreateJobRequest true "Данные для новой работы"
// @Success 201 {object} models.Job
// @Failure 400 {string} string "invalid request"
// @Failure 500 {string} string "failed to create job"
// @Router /jobs [post]
// @Security BearerAuth
func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	var req models.CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	job, err := h.JobService.CreateJob(req)
	if err != nil {
		http.Error(w, "failed to create job", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(job)
}

// GetJobs godoc
// @Summary Получить список работ (Jobs) с фильтрами и пагинацией
// @Description Получить jobs по фильтрам: кол-во комнат/офис, даты, размер грузовика, диапазон оплаты, пагинация
// @Tags jobs
// @Accept  json
// @Produce  json
// @Param relocation_size query string false "Количество комнат или office"
// @Param date_start query string false "Дата начала (ISO8601)"
// @Param date_end query string false "Дата конца (ISO8601)"
// @Param truck_size query string false "Размер грузовика (small, medium, large)"
// @Param payout_min query number false "Минимальная оплата"
// @Param payout_max query number false "Максимальная оплата"
// @Param limit query int false "Лимит (по умолчанию 10)"
// @Param offset query int false "Смещение (по умолчанию 0)"
// @Success 200 {object} models.JobListResponse
// @Failure 500 {string} string "failed to fetch jobs"
// @Router /jobs [get]
// @Security BearerAuth
func (h *JobHandler) GetJobs(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := models.JobFilter{}

	if v := q.Get("relocation_size"); v != "" {
		filter.NumberOfBedrooms = v
	}
	if v := q.Get("date_start"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			filter.DateStart = &t
		}
	}
	if v := q.Get("date_end"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			filter.DateEnd = &t
		}
	}
	if v := q.Get("truck_size"); v != "" {
		filter.TruckSize = v
	}
	if v := q.Get("payout_min"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			filter.PayoutMin = &f
		}
	}
	if v := q.Get("payout_max"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			filter.PayoutMax = &f
		}
	}
	limit := 10
	offset := 0
	if v := q.Get("limit"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			limit = i
		}
	}
	if v := q.Get("offset"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			offset = i
		}
	}

	jobs, total, err := h.JobService.GetJobs(filter, limit, offset)
	if err != nil {
		http.Error(w, "failed to fetch jobs", http.StatusInternalServerError)
		return
	}
	resp := models.JobListResponse{
		Jobs:  jobs,
		Total: total,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DeleteJob godoc
// @Summary Удалить работу (Job) по id
// @Description Удаляет работу (Job) по её id
// @Tags jobs
// @Param id path string true "ID работы"
// @Success 204 {string} string "deleted"
// @Failure 400 {string} string "invalid id"
// @Failure 404 {string} string "job not found"
// @Failure 500 {string} string "failed to delete job"
// @Security BearerAuth
// @Router /jobs/{id} [delete]
func (h *JobHandler) DeleteJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	err := h.JobService.DeleteJob(id)
	if err != nil {
		if err == services.ErrJobNotFound {
			http.Error(w, "job not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to delete job", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
