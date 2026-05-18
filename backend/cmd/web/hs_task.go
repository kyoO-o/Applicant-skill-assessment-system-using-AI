package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/cmd/web/app"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/cmd/web/socket"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/common/oapi"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/taskman"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/userman"
)

type taskRequest struct {
	JobPostingID  int    `json:"job_posting_id"`
	ApplicationID *int   `json:"application_id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	DueDate       string `json:"due_date"`
}

type generateTaskRequest struct {
	JobPostingID int `json:"job_posting_id"`
}

// POST /api/tasks  — recruiter creates task manually
func createTask(w http.ResponseWriter, r *http.Request) {
	_, ok := recruiterUser(r)
	if !ok {
		oapi.Forbidden(w)
		return
	}

	var req taskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Хүсэлт буруу байна"})
		return
	}

	if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.Description) == "" {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Гарчиг болон тайлбар шаардлагатай"})
		return
	}

	task := &taskman.Task{
		JobPostingID: uint(req.JobPostingID),
		Title:        strings.TrimSpace(req.Title),
		Description:  strings.TrimSpace(req.Description),
		Status:       taskman.TaskStatusDraft,
	}

	if req.ApplicationID != nil {
		aid := uint(*req.ApplicationID)
		task.ApplicationID = &aid
	}

	if req.DueDate != "" {
		t, err := time.Parse("2006-01-02", req.DueDate)
		if err == nil {
			task.DueDate = &t
		}
	}

	saved, err := app.Tasks.Save(task)
	if err != nil {
		oapi.ServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	oapi.SendResp(w, saved)
}

// POST /api/tasks/generate  — AI generates a task based on job requirements
func generateTask(w http.ResponseWriter, r *http.Request) {
	_, ok := recruiterUser(r)
	if !ok {
		oapi.Forbidden(w)
		return
	}

	var req generateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Хүсэлт буруу байна"})
		return
	}

	job, err := app.Jobs.Get(req.JobPostingID)
	if err != nil {
		oapi.CustomError(w, http.StatusNotFound, map[string]string{"message": "Ажлын байр олдсонгүй"})
		return
	}

	requirements := make([]string, 0, len(job.Requirements))
	for _, r := range job.Requirements {
		requirements = append(requirements, r.Description)
	}
	skills := make([]string, 0, len(job.Skills))
	for _, s := range job.Skills {
		skills = append(skills, s.Name)
	}

	title, description, err := app.AI.GenerateTask(job.Title, requirements, skills)
	if err != nil {
		app.ErrorLog.Printf("AI task generation error: %v", err)
		oapi.CustomError(w, http.StatusServiceUnavailable, map[string]string{"message": "AI үйлчилгээ түр ажиллахгүй байна"})
		return
	}

	oapi.SendResp(w, map[string]string{
		"title":       title,
		"description": description,
	})
}

// PUT /api/tasks/{id}/send  — recruiter sends a task to applicant
func sendTask(w http.ResponseWriter, r *http.Request) {
	_, ok := recruiterUser(r)
	if !ok {
		oapi.Forbidden(w)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "ID буруу байна"})
		return
	}

	task, err := app.Tasks.Get(id)
	if err != nil {
		if errors.Is(err, taskman.ErrNotFound) {
			oapi.CustomError(w, http.StatusNotFound, map[string]string{"message": "Даалгавар олдсонгүй"})
			return
		}
		oapi.ServerError(w, err)
		return
	}

	task.Status = taskman.TaskStatusSent
	saved, err := app.Tasks.Save(task)
	if err != nil {
		oapi.ServerError(w, err)
		return
	}

	if saved.ApplicationID != nil {
		if a, err := app.Applications.Get(int(*saved.ApplicationID)); err == nil {
			socket.NotifyUser(int(a.ApplicantID), "Шинэ даалгавар ирлээ", saved.Title, "task_sent")
		}
	}

	oapi.SendResp(w, saved)
}

// GET /api/tasks  — list tasks (recruiter: all for their jobs; applicant: their assigned tasks)
func listTasks(w http.ResponseWriter, r *http.Request) {
	user := authUser(r)

	if user.Role == userman.RoleRecruiter {
		// Recruiter sees tasks for jobs they own
		jobIDStr := r.URL.Query().Get("job_id")
		if jobIDStr == "" {
			oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "job_id шаардлагатай"})
			return
		}
		jobID, _ := strconv.Atoi(jobIDStr)
		tasks, err := app.Tasks.ListForJob(jobID)
		if err != nil {
			oapi.ServerError(w, err)
			return
		}
		oapi.SendResp(w, tasks)
		return
	}

	tasks, err := app.Tasks.ListForApplicant(user.ID)
	if err != nil {
		oapi.ServerError(w, err)
		return
	}
	oapi.SendResp(w, tasks)
}

// GET /api/tasks/{id}
func getTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "ID буруу байна"})
		return
	}

	task, err := app.Tasks.Get(id)
	if err != nil {
		if errors.Is(err, taskman.ErrNotFound) {
			oapi.CustomError(w, http.StatusNotFound, map[string]string{"message": "Даалгавар олдсонгүй"})
			return
		}
		oapi.ServerError(w, err)
		return
	}

	oapi.SendResp(w, task)
}

// POST /api/tasks/{id}/submit  — applicant submits a task
func submitTask(w http.ResponseWriter, r *http.Request) {
	user := authUser(r)
	if user.Role == userman.RoleRecruiter {
		oapi.Forbidden(w)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "ID буруу байна"})
		return
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Хүсэлт буруу байна"})
		return
	}

	if strings.TrimSpace(req.Content) == "" {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Гүйцэтгэлийн агуулга шаардлагатай"})
		return
	}

	sub := &taskman.TaskSubmission{
		TaskID:      uint(id),
		ApplicantID: uint(user.ID),
		Content:     strings.TrimSpace(req.Content),
		Status:      taskman.SubStatusSubmitted,
	}

	saved, err := app.Tasks.SaveSubmission(sub)
	if err != nil {
		oapi.ServerError(w, err)
		return
	}

	// Update task status
	task, _ := app.Tasks.Get(id)
	if task != nil {
		task.Status = taskman.TaskStatusCompleted
		app.Tasks.Save(task)
	}

	w.WriteHeader(http.StatusCreated)
	oapi.SendResp(w, saved)
}

// PUT /api/tasks/submissions/{id}/grade  — recruiter grades a submission
func gradeSubmission(w http.ResponseWriter, r *http.Request) {
	_, ok := recruiterUser(r)
	if !ok {
		oapi.Forbidden(w)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "ID буруу байна"})
		return
	}

	var req struct {
		Grade    int    `json:"grade"`
		Feedback string `json:"feedback"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Хүсэлт буруу байна"})
		return
	}

	sub, err := app.Tasks.GetSubmission(id)
	if err != nil {
		if errors.Is(err, taskman.ErrNotFound) {
			oapi.CustomError(w, http.StatusNotFound, map[string]string{"message": "Илгээлт олдсонгүй"})
			return
		}
		oapi.ServerError(w, err)
		return
	}

	sub.Grade = &req.Grade
	sub.Feedback = req.Feedback
	sub.Status = taskman.SubStatusGraded

	saved, err := app.Tasks.SaveSubmission(sub)
	if err != nil {
		oapi.ServerError(w, err)
		return
	}

	socket.NotifyUser(int(saved.ApplicantID), "Даалгавар үнэлэгдлээ",
		fmt.Sprintf("Таны даалгавар үнэлэгдлээ. Оноо: %d/100", *saved.Grade), "task_graded")

	oapi.SendResp(w, saved)
}
