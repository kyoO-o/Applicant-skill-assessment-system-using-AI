package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/cmd/web/app"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/common/oapi"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/companyman"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/jobman"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/userman"
)

type jobRequest struct {
	Title          string   `json:"title"`
	Location       string   `json:"location"`
	AdditionalInfo string   `json:"additional_info"`
	Description    string   `json:"description"`
	ContactInfo    string   `json:"contact_info"`
	Type           string   `json:"type"`
	EmploymentType string   `json:"employment_type"`
	Level          string   `json:"level"`
	Seniority      string   `json:"seniority"`
	City           *string  `json:"city"`
	District       *string  `json:"district"`
	LocationX      *float64 `json:"location_x"`
	LocationY      *float64 `json:"location_y"`
	MinSalary      float64  `json:"min_salary"`
	MaxSalary      float64  `json:"max_salary"`
	Status         string   `json:"status"`
	Duties         []string `json:"duties"`
	Requirements   []string `json:"requirements"`
	Skills         []string `json:"skills"`
	Bonuses        []string `json:"bonuses"`
}

type jobResponse struct {
	ID                 int       `json:"id"`
	RecruiterID        int       `json:"recruiter_id"`
	CompanyID          int       `json:"company_id"`
	CompanyName        string    `json:"company_name"`
	CompanyLogoURL     string    `json:"company_logo_url"`
	CompanyProfileURL  string    `json:"company_profile_url"`
	Title           string    `json:"title"`
	Location        string    `json:"location"`
	AdditionalInfo  string    `json:"additional_info"`
	Description     string    `json:"description"`
	ContactInfo     string    `json:"contact_info"`
	Type            string    `json:"type"`
	EmploymentType  string    `json:"employment_type"`
	Level           string    `json:"level"`
	Seniority       string    `json:"seniority"`
	City            *string   `json:"city"`
	District        *string   `json:"district"`
	LocationX       *float64  `json:"location_x"`
	LocationY       *float64  `json:"location_y"`
	MinSalary       float64   `json:"min_salary"`
	MaxSalary       float64   `json:"max_salary"`
	Status          string    `json:"status"`
	Duties          []string  `json:"duties"`
	Requirements    []string  `json:"requirements"`
	Skills          []string  `json:"skills"`
	Bonuses         []string  `json:"bonuses"`
	ApplicantsCount int       `json:"applicants_count"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func mapJobResponse(job *jobman.JobPosting) *jobResponse {
	if job == nil {
		return nil
	}

	duties := make([]string, 0, len(job.Duties))
	for _, duty := range job.Duties {
		duty.Description = strings.TrimSpace(duty.Description)
		if duty.Description == "" {
			continue
		}
		duties = append(duties, duty.Description)
	}

	requirements := make([]string, 0, len(job.Requirements))
	for _, requirement := range job.Requirements {
		requirement.Description = strings.TrimSpace(requirement.Description)
		if requirement.Description == "" {
			continue
		}
		requirements = append(requirements, requirement.Description)
	}

	skills := make([]string, 0, len(job.Skills))
	for _, skill := range job.Skills {
		skill.Name = strings.TrimSpace(skill.Name)
		if skill.Name == "" {
			continue
		}
		skills = append(skills, skill.Name)
	}

	bonuses := make([]string, 0, len(job.Bonuses))
	for _, bonus := range job.Bonuses {
		bonus.Description = strings.TrimSpace(bonus.Description)
		if bonus.Description == "" {
			continue
		}
		bonuses = append(bonuses, bonus.Description)
	}

	return &jobResponse{
		ID:              job.ID,
		RecruiterID:     int(job.PostedBy),
		CompanyID:       int(job.CompanyID),
		Title:           job.Title,
		Location:        job.Location,
		AdditionalInfo:  job.AdditionalInfo,
		Description:     job.AdditionalInfo,
		ContactInfo:     job.ContactInfo,
		Type:            job.Type,
		EmploymentType:  job.Type,
		Level:           job.Level,
		Seniority:       job.Level,
		City:            job.City,
		District:        job.District,
		LocationX:       job.LocationX,
		LocationY:       job.LocationY,
		MinSalary:       job.MinSalary,
		MaxSalary:       job.MaxSalary,
		Status:          job.Status,
		Duties:          duties,
		Requirements:    requirements,
		Skills:          skills,
		Bonuses:         bonuses,
		ApplicantsCount: job.AppsCount,
		CreatedAt:       job.CreatedAt,
		UpdatedAt:       job.UpdatedAt,
	}
}

func cleanOptionalText(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func cleanList(values []string) []string {
	items := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		items = append(items, value)
	}
	return items
}

func enrichJobsWithCompany(items []*jobResponse) {
	if len(items) == 0 {
		return
	}
	ids := make([]int, 0, len(items))
	seen := map[int]bool{}
	for _, item := range items {
		if !seen[item.CompanyID] {
			ids = append(ids, item.CompanyID)
			seen[item.CompanyID] = true
		}
	}
	var companies []companyman.Company
	if err := app.DB.Where("id IN ?", ids).Find(&companies).Error; err != nil {
		return
	}
	byID := make(map[int]*companyman.Company, len(companies))
	for i := range companies {
		byID[companies[i].ID] = &companies[i]
	}
	for _, item := range items {
		if co, ok := byID[item.CompanyID]; ok {
			item.CompanyName = co.Name
			item.CompanyLogoURL = co.Logo
			item.CompanyProfileURL = co.ProfileURL
		}
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
	enrichJobsWithCompany(items)

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
	enrichJobsWithCompany(items)

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
	req.AdditionalInfo = strings.TrimSpace(req.AdditionalInfo)
	req.Description = strings.TrimSpace(req.Description)
	req.ContactInfo = strings.TrimSpace(req.ContactInfo)
	req.Type = strings.TrimSpace(req.Type)
	req.EmploymentType = strings.TrimSpace(req.EmploymentType)
	req.Level = strings.TrimSpace(req.Level)
	req.Seniority = strings.TrimSpace(req.Seniority)
	req.Status = strings.TrimSpace(strings.ToLower(req.Status))
	req.City = cleanOptionalText(req.City)
	req.District = cleanOptionalText(req.District)

	if req.Type == "" {
		req.Type = req.EmploymentType
	}
	if req.Level == "" {
		req.Level = req.Seniority
	}
	if req.AdditionalInfo == "" {
		req.AdditionalInfo = req.Description
	}

	if req.Location == "" {
		parts := make([]string, 0, 2)
		if req.District != nil {
			parts = append(parts, *req.District)
		}
		if req.City != nil {
			parts = append(parts, *req.City)
		}
		req.Location = strings.Join(parts, ", ")
	}

	if req.Title == "" || req.Location == "" || req.AdditionalInfo == "" {
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
	for _, requirement := range cleanList(req.Requirements) {
		requirements = append(requirements, jobman.Requirement{Description: requirement})
	}

	duties := make([]jobman.Duty, 0, len(req.Duties))
	for _, duty := range cleanList(req.Duties) {
		duties = append(duties, jobman.Duty{Description: duty})
	}

	skills := make([]jobman.Skill, 0, len(req.Skills))
	for _, skill := range cleanList(req.Skills) {
		skills = append(skills, jobman.Skill{Name: skill})
	}

	bonuses := make([]jobman.Bonus, 0, len(req.Bonuses))
	for _, bonus := range cleanList(req.Bonuses) {
		bonuses = append(bonuses, jobman.Bonus{Description: bonus})
	}

	job.Title = req.Title
	job.Location = req.Location
	job.AdditionalInfo = req.AdditionalInfo
	job.ContactInfo = req.ContactInfo
	job.Type = req.Type
	job.Level = req.Level
	job.City = req.City
	job.District = req.District
	job.LocationX = req.LocationX
	job.LocationY = req.LocationY
	job.MinSalary = req.MinSalary
	job.MaxSalary = req.MaxSalary
	job.Duties = duties
	job.Requirements = requirements
	job.Skills = skills
	job.Bonuses = bonuses
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
	req.AdditionalInfo = strings.TrimSpace(req.AdditionalInfo)
	req.Description = strings.TrimSpace(req.Description)
	req.ContactInfo = strings.TrimSpace(req.ContactInfo)
	req.Type = strings.TrimSpace(req.Type)
	req.EmploymentType = strings.TrimSpace(req.EmploymentType)
	req.Level = strings.TrimSpace(req.Level)
	req.Seniority = strings.TrimSpace(req.Seniority)
	req.Status = strings.TrimSpace(strings.ToLower(req.Status))
	req.City = cleanOptionalText(req.City)
	req.District = cleanOptionalText(req.District)

	if req.Type == "" {
		req.Type = req.EmploymentType
	}
	if req.Level == "" {
		req.Level = req.Seniority
	}
	if req.AdditionalInfo == "" {
		req.AdditionalInfo = req.Description
	}

	if req.Location == "" {
		parts := make([]string, 0, 2)
		if req.District != nil {
			parts = append(parts, *req.District)
		}
		if req.City != nil {
			parts = append(parts, *req.City)
		}
		req.Location = strings.Join(parts, ", ")
	}

	if req.Title == "" || req.Location == "" || req.AdditionalInfo == "" {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Title, location, and description are required"})
		return
	}

	if req.Status != "" && req.Status != jobman.StatusDraft && req.Status != jobman.StatusPosted && req.Status != jobman.StatusClosed {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Status must be draft, posted, or closed"})
		return
	}

	requirements := make([]jobman.Requirement, 0, len(req.Requirements))
	for _, requirement := range cleanList(req.Requirements) {
		requirements = append(requirements, jobman.Requirement{Description: requirement})
	}

	duties := make([]jobman.Duty, 0, len(req.Duties))
	for _, duty := range cleanList(req.Duties) {
		duties = append(duties, jobman.Duty{Description: duty})
	}

	skills := make([]jobman.Skill, 0, len(req.Skills))
	for _, skill := range cleanList(req.Skills) {
		skills = append(skills, jobman.Skill{Name: skill})
	}

	bonuses := make([]jobman.Bonus, 0, len(req.Bonuses))
	for _, bonus := range cleanList(req.Bonuses) {
		bonuses = append(bonuses, jobman.Bonus{Description: bonus})
	}

	job.Title = req.Title
	job.Location = req.Location
	job.AdditionalInfo = req.AdditionalInfo
	job.ContactInfo = req.ContactInfo
	job.Type = req.Type
	job.Level = req.Level
	job.City = req.City
	job.District = req.District
	job.LocationX = req.LocationX
	job.LocationY = req.LocationY
	job.MinSalary = req.MinSalary
	job.MaxSalary = req.MaxSalary
	job.Duties = duties
	job.Requirements = requirements
	job.Skills = skills
	job.Bonuses = bonuses
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
