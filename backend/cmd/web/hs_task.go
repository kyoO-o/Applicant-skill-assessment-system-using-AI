package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
	JobPostingID  *int   `json:"job_posting_id"`
	ApplicationID *int   `json:"application_id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	DueDate       string `json:"due_date"`
	DurationDays  *int   `json:"duration_days"`
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
		Title:       strings.TrimSpace(req.Title),
		Description: strings.TrimSpace(req.Description),
		Status:      taskman.TaskStatusDraft,
	}

	if req.JobPostingID != nil && *req.JobPostingID > 0 {
		jid := uint(*req.JobPostingID)
		task.JobPostingID = &jid
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

	if req.DurationDays != nil && *req.DurationDays > 0 {
		task.DurationDays = req.DurationDays
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

// PUT /api/tasks/{id}  — recruiter updates task title/description/due_date
func updateTask(w http.ResponseWriter, r *http.Request) {
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

	var req taskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Хүсэлт буруу байна"})
		return
	}

	if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.Description) == "" {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Гарчиг болон тайлбар шаардлагатай"})
		return
	}

	task.Title = strings.TrimSpace(req.Title)
	task.Description = strings.TrimSpace(req.Description)

	if req.DueDate != "" {
		t, err := time.Parse("2006-01-02", req.DueDate)
		if err == nil {
			task.DueDate = &t
		}
	} else {
		task.DueDate = nil
	}

	if req.DurationDays != nil && *req.DurationDays > 0 {
		task.DurationDays = req.DurationDays
	} else {
		task.DurationDays = nil
	}

	saved, err := app.Tasks.Save(task)
	if err != nil {
		oapi.ServerError(w, err)
		return
	}

	oapi.SendResp(w, saved)
}

// DELETE /api/tasks/{id}  — recruiter deletes a draft task
func deleteTask(w http.ResponseWriter, r *http.Request) {
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

	if task.Status != taskman.TaskStatusDraft {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Зөвхөн ноорог даалгавар устгах боломжтой"})
		return
	}

	if err := app.Tasks.Delete(id); err != nil {
		oapi.ServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

	var body struct {
		ApplicationID *int   `json:"application_id"`
		DueDate       string `json:"due_date"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	if body.ApplicationID != nil && *body.ApplicationID > 0 {
		aid := uint(*body.ApplicationID)
		task.ApplicationID = &aid
	}

	if body.DueDate != "" {
		t, err := time.Parse("2006-01-02", body.DueDate)
		if err == nil {
			task.DueDate = &t
		}
	} else if task.DueDate == nil && task.DurationDays != nil {
		due := time.Now().AddDate(0, 0, *task.DurationDays)
		task.DueDate = &due
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
		filter := &taskman.Filter{
			Keyword: r.URL.Query().Get("keyword"),
			Status:  r.URL.Query().Get("status"),
		}
		jobIDStr := r.URL.Query().Get("job_id")
		if jobIDStr == "" {
			tasks, err := app.Tasks.ListLibrary(filter)
			if err != nil {
				oapi.ServerError(w, err)
				return
			}
			oapi.SendResp(w, tasks)
			return
		}
		jobID, _ := strconv.Atoi(jobIDStr)
		tasks, err := app.Tasks.ListForJob(jobID, filter)
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

// POST /api/tasks/{id}/submit  — applicant submits a task (multipart: content + optional PDF)
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

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Хүсэлт буруу байна"})
		return
	}

	content := strings.TrimSpace(r.FormValue("content"))

	var filePath string
	file, header, fileErr := r.FormFile("file")
	if fileErr == nil {
		defer file.Close()
		if filepath.Ext(header.Filename) != ".pdf" {
			oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Зөвхөн PDF файл зөвшөөрнө"})
			return
		}
		pdfBytes, err := io.ReadAll(file)
		if err != nil {
			oapi.ServerError(w, err)
			return
		}
		storagePath := app.Config.StoragePath
		if storagePath == "" {
			storagePath = "./storage"
		}
		dir := filepath.Join(storagePath, "task_submissions")
		os.MkdirAll(dir, 0755)
		fname := strconv.Itoa(user.ID) + "_" + strconv.Itoa(id) + "_" + strconv.FormatInt(time.Now().Unix(), 10) + ".pdf"
		fp := filepath.Join(dir, fname)
		if writeErr := os.WriteFile(fp, pdfBytes, 0644); writeErr != nil {
			app.ErrorLog.Printf("warning: failed to save submission file: %v", writeErr)
		} else {
			filePath = fp
		}
	}

	if content == "" && filePath == "" {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Агуулга эсвэл PDF файл шаардлагатай"})
		return
	}

	sub := &taskman.TaskSubmission{
		TaskID:      uint(id),
		ApplicantID: uint(user.ID),
		Content:     content,
		FilePath:    filePath,
		Status:      taskman.SubStatusSubmitted,
	}

	saved, err := app.Tasks.SaveSubmission(sub)
	if err != nil {
		oapi.ServerError(w, err)
		return
	}

	task, _ := app.Tasks.Get(id)
	if task != nil {
		task.Status = taskman.TaskStatusCompleted
		app.Tasks.Save(task)

		if recruiterID := recruiterIDForTask(task); recruiterID > 0 {
			socket.NotifyUser(recruiterID, "Шинэ даалгаврын хариу ирлээ",
				fmt.Sprintf("'%s' даалгаврын хариу ирлээ", task.Title), "new_submission")
		}
	}

	w.WriteHeader(http.StatusCreated)
	oapi.SendResp(w, saved)
}

// recruiterIDForTask resolves the recruiter (job owner) for a given task.
func recruiterIDForTask(task *taskman.Task) int {
	if task.JobPostingID != nil {
		if job, err := app.Jobs.Get(int(*task.JobPostingID)); err == nil {
			return int(job.PostedBy)
		}
	}
	if task.ApplicationID != nil {
		if a, err := app.Applications.Get(int(*task.ApplicationID)); err == nil {
			if job, err := app.Jobs.Get(int(a.JobPostingID)); err == nil {
				return int(job.PostedBy)
			}
		}
	}
	return 0
}

// GET /api/tasks/submissions  — recruiter lists all submissions across all tasks
func listAllSubmissions(w http.ResponseWriter, r *http.Request) {
	_, ok := recruiterUser(r)
	if !ok {
		oapi.Forbidden(w)
		return
	}

	subs, err := app.Tasks.ListAllSubmissions()
	if err != nil {
		oapi.ServerError(w, err)
		return
	}

	oapi.SendResp(w, subs)
}

// POST /api/tasks/submissions/{id}/ai-grade  — AI grades a submission
// func aiGradeSubmission(w http.ResponseWriter, r *http.Request) {
// 	_, ok := recruiterUser(r)
// 	if !ok {
// 		oapi.Forbidden(w)
// 		return
// 	}

// 	id, err := strconv.Atoi(chi.URLParam(r, "id"))
// 	if err != nil || id <= 0 {
// 		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "ID буруу байна"})
// 		return
// 	}

// 	sub, err := app.Tasks.GetSubmission(id)
// 	if err != nil {
// 		if errors.Is(err, taskman.ErrNotFound) {
// 			oapi.CustomError(w, http.StatusNotFound, map[string]string{"message": "Илгээлт олдсонгүй"})
// 			return
// 		}
// 		oapi.ServerError(w, err)
// 		return
// 	}

// 	task, err := app.Tasks.Get(int(sub.TaskID))
// 	if err != nil {
// 		oapi.ServerError(w, err)
// 		return
// 	}

// 	content := sub.Content
// 	if content == "" && sub.FilePath != "" {
// 		if pdfBytes, readErr := os.ReadFile(sub.FilePath); readErr == nil {
// 			if extracted, _ := aiman.ExtractTextFromPDF(pdfBytes); len(extracted) > 50 {
// 				content = extracted
// 			}
// 		}
// 	}

// 	if content == "" {
// 		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "AI үнэлэхэд агуулга байхгүй байна"})
// 		return
// 	}

// 	grade, feedback, err := app.AI.GradeTask(task.Title, task.Description, content)
// 	if err != nil {
// 		app.ErrorLog.Printf("AI grading error: %v", err)
// 		oapi.CustomError(w, http.StatusServiceUnavailable, map[string]string{"message": "AI үйлчилгээ түр ажиллахгүй байна"})
// 		return
// 	}

// 	sub.Grade = &grade
// 	sub.Feedback = feedback
// 	sub.Status = taskman.SubStatusGraded

// 	saved, err := app.Tasks.SaveSubmission(sub)
// 	if err != nil {
// 		oapi.ServerError(w, err)
// 		return
// 	}

// 	socket.NotifyUser(int(saved.ApplicantID), "Даалгавар үнэлэгдлээ",
// 		fmt.Sprintf("Таны даалгавар үнэлэгдлээ. Оноо: %d/100", *saved.Grade), "task_graded")

// 	oapi.SendResp(w, saved)
// }

// GET /api/tasks/submissions/{id}/file  — serve submission PDF
func serveSubmissionFile(w http.ResponseWriter, r *http.Request) {
	user := authUser(r)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "ID буруу байна"})
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

	// Applicant can only access their own submission file
	if user.Role != userman.RoleRecruiter && int(sub.ApplicantID) != user.ID {
		oapi.Forbidden(w)
		return
	}

	if sub.FilePath == "" {
		oapi.CustomError(w, http.StatusNotFound, map[string]string{"message": "Файл олдсонгүй"})
		return
	}

	f, err := os.Open(sub.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			oapi.CustomError(w, http.StatusNotFound, map[string]string{"message": "Файл олдсонгүй"})
			return
		}
		oapi.ServerError(w, err)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		oapi.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, filepath.Base(sub.FilePath)))
	w.Header().Set("Content-Length", strconv.FormatInt(fi.Size(), 10))
	// Wrap writer to prevent io.Copy from using the ReadFrom fast-path,
	// which panics in chi v1.5.5 when the sessions middleware wraps the ResponseWriter.
	io.Copy(struct{ io.Writer }{w}, f)
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
