package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/cmd/web/app"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/common/oapi"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/jobman"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/userman"
)

type jobRequest struct {
	Title          string   `json:"title"`
	Location       string   `json:"location"`
	EmploymentType string   `json:"employment_type"`
	Seniority      string   `json:"seniority"`
	Status         string   `json:"status"`
	Description    string   `json:"description"`
	Requirements   []string `json:"requirements"`
}

type jobResponse struct {
	ID              int       `json:"id"`
	RecruiterID     int       `json:"recruiter_id"`
	CompanyID       int       `json:"company_id"`
	CompanyName     string    `json:"company_name"`
	Title           string    `json:"title"`
	Location        string    `json:"location"`
	EmploymentType  string    `json:"employment_type"`
	Seniority       string    `json:"seniority"`
	Status          string    `json:"status"`
	Description     string    `json:"description"`
	Requirements    []string  `json:"requirements"`
	ApplicantsCount int       `json:"applicants_count"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func mapJobResponse(job *jobman.JobPosting) *jobResponse {
	if job == nil {
		return nil
	}

	requirements := make([]string, 0, len(job.Requirements))
	for _, requirement := range job.Requirements {
		requirement.Description = strings.TrimSpace(requirement.Description)
		if requirement.Description == "" {
			continue
		}
		requirements = append(requirements, requirement.Description)
	}

	return &jobResponse{
		ID:              job.ID,
		RecruiterID:     int(job.PostedBy),
		CompanyID:       int(job.CompanyID),
		Title:           job.Title,
		Location:        job.Location,
		EmploymentType:  job.Type,
		Seniority:       job.Level,
		Status:          job.Status,
		Description:     job.AdditionalInfo,
		Requirements:    requirements,
		ApplicantsCount: job.AppsCount,
		CreatedAt:       job.CreatedAt,
		UpdatedAt:       job.UpdatedAt,
	}
}

func chosenJob(r *http.Request) (*jobman.JobPosting, bool) {
	job, ok := r.Context().Value(app.ContextKeyChosenJob).(*jobman.JobPosting)
	return job, ok
}

func getJobs(w http.ResponseWriter, r *http.Request) {
	user := authUser(r)

	var (
		jobs []*jobman.JobPosting
		err  error
	)

	if user.Role == userman.RoleRecruiter {
		jobs, err = app.Jobs.ListForRecruiter(int(user.ID))
	} else {
		jobs, err = app.Jobs.ListActive()
	}

	if err != nil {
		oapi.ServerError(w, err)
		return
	}

	items := make([]*jobResponse, 0, len(jobs))
	for _, job := range jobs {
		items = append(items, mapJobResponse(job))
	}

	oapi.SendResp(w, items)
}

func getCompanyJobs(w http.ResponseWriter, r *http.Request) {
	user, ok := recruiterUser(r)
	if !ok {
		oapi.Forbidden(w)
		return
	}

	company, ok := chosenCompany(r)
	if !ok {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Company context is required"})
		return
	}

	jobs, err := app.Jobs.ListForRecruiterAndCompany(int(user.ID), company.ID)
	if err != nil {
		oapi.ServerError(w, err)
		return
	}

	items := make([]*jobResponse, 0, len(jobs))
	for _, job := range jobs {
		items = append(items, mapJobResponse(job))
	}

	oapi.SendResp(w, items)
}

func saveJob(w http.ResponseWriter, r *http.Request) {
	user, ok := recruiterUser(r)
	if !ok {
		oapi.Forbidden(w)
		return
	}

	company, hasChosenCompany := chosenCompany(r)
	if r.Method != http.MethodPost && !hasChosenCompany {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Company context is required"})
		return
	}

	var req jobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Location = strings.TrimSpace(req.Location)
	req.EmploymentType = strings.TrimSpace(req.EmploymentType)
	req.Seniority = strings.TrimSpace(req.Seniority)
	req.Status = strings.TrimSpace(strings.ToLower(req.Status))
	req.Description = strings.TrimSpace(req.Description)

	if req.Title == "" || req.Location == "" || req.Description == "" {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Title, location, and description are required"})
		return
	}

	if req.Status != "" && req.Status != jobman.StatusDraft && req.Status != jobman.StatusPosted && req.Status != jobman.StatusClosed {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Status must be draft, posted, or closed"})
		return
	}

	if user.CompanyID == nil && !hasChosenCompany {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Company context is required"})
		return
	}

	job := &jobman.JobPosting{
		PostedBy: uint(user.ID),
		Status:   jobman.StatusDraft,
	}
	if hasChosenCompany {
		job.CompanyID = uint(company.ID)
	} else {
		job.CompanyID = *user.CompanyID
	}

	requirements := make([]jobman.Requirement, 0, len(req.Requirements))
	for _, requirement := range req.Requirements {
		requirement = strings.TrimSpace(requirement)
		if requirement == "" {
			continue
		}
		requirements = append(requirements, jobman.Requirement{Description: requirement})
	}

	job.Title = req.Title
	job.Location = req.Location
	job.AdditionalInfo = req.Description
	job.Type = req.EmploymentType
	job.Level = req.Seniority
	job.Requirements = requirements
	if hasChosenCompany {
		job.CompanyID = uint(company.ID)
	}
	if req.Status != "" {
		job.Status = req.Status
	}

	savedJob, err := app.Jobs.Save(job)
	if err != nil {
		oapi.ServerError(w, err)
		return
	}

	if r.Method == http.MethodPost {
		w.WriteHeader(http.StatusCreated)
	}
	oapi.SendResp(w, mapJobResponse(savedJob))
}

func getJob(w http.ResponseWriter, r *http.Request) {
	job, ok := chosenJob(r)
	if !ok {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Job context is required"})
		return
	}

	oapi.SendResp(w, mapJobResponse(job))
}

func saveCompanyJob(w http.ResponseWriter, r *http.Request) {
	_, ok := recruiterUser(r)
	if !ok {
		oapi.Forbidden(w)
		return
	}

	job, ok := chosenJob(r)
	if !ok {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Job context is required"})
		return
	}

	var req jobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Location = strings.TrimSpace(req.Location)
	req.EmploymentType = strings.TrimSpace(req.EmploymentType)
	req.Seniority = strings.TrimSpace(req.Seniority)
	req.Status = strings.TrimSpace(strings.ToLower(req.Status))
	req.Description = strings.TrimSpace(req.Description)

	if req.Title == "" || req.Location == "" || req.Description == "" {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Title, location, and description are required"})
		return
	}

	if req.Status != "" && req.Status != jobman.StatusDraft && req.Status != jobman.StatusPosted && req.Status != jobman.StatusClosed {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Status must be draft, posted, or closed"})
		return
	}

	requirements := make([]jobman.Requirement, 0, len(req.Requirements))
	for _, requirement := range req.Requirements {
		requirement = strings.TrimSpace(requirement)
		if requirement == "" {
			continue
		}
		requirements = append(requirements, jobman.Requirement{Description: requirement})
	}

	job.Title = req.Title
	job.Location = req.Location
	job.AdditionalInfo = req.Description
	job.Type = req.EmploymentType
	job.Level = req.Seniority
	job.Requirements = requirements
	if req.Status != "" {
		job.Status = req.Status
	}

	savedJob, err := app.Jobs.Save(job)
	if err != nil {
		oapi.ServerError(w, err)
		return
	}

	oapi.SendResp(w, mapJobResponse(savedJob))
}

func deleteJob(w http.ResponseWriter, r *http.Request) {
	_, ok := recruiterUser(r)
	if !ok {
		oapi.Forbidden(w)
		return
	}

	job, ok := chosenJob(r)
	if !ok {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Job context is required"})
		return
	}

	if err := app.Jobs.Delete(job.ID); err != nil {
		oapi.ServerError(w, err)
		return
	}

	oapi.SendResp(w, map[string]string{"message": "Job deleted"})
}
