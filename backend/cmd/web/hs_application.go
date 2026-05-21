package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/cmd/web/app"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/cmd/web/socket"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/common/oapi"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/aiman"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/appman"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/userman"
)

type applicationResponse struct {
	ID                     int                       `json:"id"`
	JobPostingID           int                       `json:"job_posting_id"`
	JobTitle               string                    `json:"job_title"`
	ApplicantID            int                       `json:"applicant_id"`
	ApplicantName          string                    `json:"applicant_name"`
	ApplicantEmail         string                    `json:"applicant_email"`
	ApplicantProfileURL    string                    `json:"applicant_profile_url"`
	OverallScore           int                       `json:"overall_score"`
	Summary                string                    `json:"summary"`
	MatchedSkills          []aiman.SkillResult       `json:"matched_skills"`
	MissingSkills          []aiman.SkillResult       `json:"missing_skills"`
	Recommendations        []string                  `json:"recommendations"`
	DutyAssessments        []aiman.DutyAssessment    `json:"duty_assessments"`
	RequirementAssessments []aiman.RequirementAssessment `json:"requirement_assessments"`
	Status                 string                    `json:"status"`
	AssessedAt             *time.Time                `json:"assessed_at"`
	InterviewAt            *time.Time                `json:"interview_at"`
	InterviewLocation      string                    `json:"interview_location"`
	InterviewNote          string                    `json:"interview_note"`
	CreatedAt              time.Time                 `json:"created_at"`
}

func mapApplicationResponse(a *appman.Application, jobTitle string) *applicationResponse {
	var matched []aiman.SkillResult
	var missing []aiman.SkillResult
	var recs []string
	var dutyAssessments []aiman.DutyAssessment
	var reqAssessments []aiman.RequirementAssessment

	json.Unmarshal([]byte(a.MatchedSkills), &matched)
	json.Unmarshal([]byte(a.MissingSkills), &missing)
	json.Unmarshal([]byte(a.Recommendations), &recs)
	json.Unmarshal([]byte(a.DutyAssessments), &dutyAssessments)
	json.Unmarshal([]byte(a.RequirementAssessments), &reqAssessments)

	return &applicationResponse{
		ID:                     a.ID,
		JobPostingID:           int(a.JobPostingID),
		JobTitle:               jobTitle,
		ApplicantID:            int(a.ApplicantID),
		ApplicantName:          a.ApplicantName,
		ApplicantEmail:         a.ApplicantEmail,
		OverallScore:           a.OverallScore,
		Summary:                a.Summary,
		MatchedSkills:          matched,
		MissingSkills:          missing,
		Recommendations:        recs,
		DutyAssessments:        dutyAssessments,
		RequirementAssessments: reqAssessments,
		Status:                 a.Status,
		AssessedAt:             a.AssessedAt,
		InterviewAt:            a.InterviewAt,
		InterviewLocation:      a.InterviewLocation,
		InterviewNote:          a.InterviewNote,
		CreatedAt:              a.CreatedAt,
	}
}

// POST /api/jobs/{JobID}/apply  — applicant uploads CV PDF
func applyToJob(w http.ResponseWriter, r *http.Request) {
	user := authUser(r)
	if user.Role == userman.RoleRecruiter {
		oapi.Forbidden(w)
		return
	}

	jobID, err := strconv.Atoi(chi.URLParam(r, "JobID"))
	if err != nil || jobID <= 0 {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Ажлын байрны ID буруу байна"})
		return
	}

	job, err := app.Jobs.Get(jobID)
	if err != nil {
		oapi.CustomError(w, http.StatusNotFound, map[string]string{"message": "Ажлын байр олдсонгүй"})
		return
	}

	// Delete existing application to allow re-apply
	if existing, err := app.Applications.GetForApplicantAndJob(user.ID, jobID); err == nil {
		if delErr := app.Applications.Delete(existing.ID); delErr != nil {
			oapi.ServerError(w, delErr)
			return
		}
		// Remove old CV file if present
		if existing.CVFilePath != "" {
			os.Remove(existing.CVFilePath)
		}
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Файл уншихад алдаа гарлаа"})
		return
	}

	file, header, err := r.FormFile("cv")
	if err != nil {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "CV файл сонгоно уу"})
		return
	}
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

	// Extract text from PDF
	cvText, err := aiman.ExtractTextFromPDF(pdfBytes)
	if err != nil || len(cvText) < 50 {
		cvText = "CV текст уншихад алдаа гарлаа - үнэлгээ хийх боломжгүй"
	}

	// Save PDF to disk
	storagePath := app.Config.StoragePath
	if storagePath == "" {
		storagePath = "./storage"
	}
	cvDir := filepath.Join(storagePath, "cvs")
	os.MkdirAll(cvDir, 0755)
	filename := strconv.Itoa(user.ID) + "_" + strconv.Itoa(jobID) + "_" + strconv.FormatInt(time.Now().Unix(), 10) + ".pdf"
	filePath := filepath.Join(cvDir, filename)
	if err := os.WriteFile(filePath, pdfBytes, 0644); err != nil {
		app.ErrorLog.Printf("warning: failed to save CV file: %v", err)
	}

	// Create application record
	application := &appman.Application{
		JobPostingID: uint(jobID),
		ApplicantID:  uint(user.ID),
		CVFilePath:   filePath,
		CVText:       cvText,
		Status:       appman.StatusPending,
	}

	savedApp, err := app.Applications.Save(application)
	if err != nil {
		oapi.ServerError(w, err)
		return
	}

	socket.NotifyUser(int(job.PostedBy), "Шинэ анкет ирлээ",
		fmt.Sprintf("'%s' ажлын байранд шинэ анкет ирлээ", job.Title), "new_application")

	// Run AI assessment asynchronously
	jobRequirements := make([]string, 0, len(job.Requirements))
	for _, r := range job.Requirements {
		jobRequirements = append(jobRequirements, r.Description)
	}
	jobSkills := make([]string, 0, len(job.Skills))
	for _, s := range job.Skills {
		jobSkills = append(jobSkills, s.Name)
	}
	jobDuties := make([]string, 0, len(job.Duties))
	for _, d := range job.Duties {
		jobDuties = append(jobDuties, d.Description)
	}

	go func(appID int, cvText, jobTitle string, requirements, skills, duties []string) {
		result, err := app.AI.AssessCV(cvText, jobTitle, requirements, skills, duties)
		if err != nil {
			app.ErrorLog.Printf("AI assessment failed for application %d: %v", appID, err)
			return
		}

		matchedJSON, _ := json.Marshal(result.MatchedSkills)
		missingJSON, _ := json.Marshal(result.MissingSkills)
		recsJSON, _ := json.Marshal(result.Recommendations)
		dutyJSON, _ := json.Marshal(result.DutyAssessments)
		reqJSON, _ := json.Marshal(result.RequirementAssessments)
		now := time.Now()

		a, _ := app.Applications.Get(appID)
		if a == nil {
			return
		}
		a.OverallScore = result.OverallScore
		a.Summary = result.Summary
		a.MatchedSkills = string(matchedJSON)
		a.MissingSkills = string(missingJSON)
		a.Recommendations = string(recsJSON)
		a.DutyAssessments = string(dutyJSON)
		a.RequirementAssessments = string(reqJSON)
		a.Status = appman.StatusAssessed
		a.AssessedAt = &now
		app.Applications.Save(a)

		socket.NotifyUser(int(a.ApplicantID), "AI үнэлгээ дууслаа",
			fmt.Sprintf("'%s' ажлын байранд таны анкет %d оноо авлаа", jobTitle, result.OverallScore), "assessment_complete")
	}(savedApp.ID, cvText, job.Title, jobRequirements, jobSkills, jobDuties)

	w.WriteHeader(http.StatusCreated)
	oapi.SendResp(w, map[string]any{
		"message": "Анкет амжилттай илгээгдлээ. AI үнэлгээ хийгдэж байна...",
		"id":      savedApp.ID,
	})
}


// GET /api/applications  — applicant sees their own applications
func listMyApplications(w http.ResponseWriter, r *http.Request) {
	user := authUser(r)
	filter := &appman.Filter{Status: r.URL.Query().Get("status")}
	apps, err := app.Applications.ListForApplicant(user.ID, filter)
	if err != nil {
		oapi.ServerError(w, err)
		return
	}

	items := make([]*applicationResponse, 0, len(apps))
	for _, a := range apps {
		job, _ := app.Jobs.Get(int(a.JobPostingID))
		jobTitle := ""
		if job != nil {
			jobTitle = job.Title
		}
		resp := mapApplicationResponse(a, jobTitle)
		resp.ApplicantName = user.FullName
		resp.ApplicantEmail = user.Email
		items = append(items, resp)
	}

	oapi.SendResp(w, items)
}

// GET /api/applications/{id}  — get single application with full assessment
func getApplication(w http.ResponseWriter, r *http.Request) {
	user := authUser(r)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "ID буруу байна"})
		return
	}

	a, err := app.Applications.Get(id)
	if err != nil {
		if errors.Is(err, appman.ErrNotFound) {
			oapi.CustomError(w, http.StatusNotFound, map[string]string{"message": "Олдсонгүй"})
			return
		}
		oapi.ServerError(w, err)
		return
	}

	// Applicant can only see own; recruiter can see for their jobs
	if user.Role != userman.RoleRecruiter && int(a.ApplicantID) != user.ID {
		oapi.Forbidden(w)
		return
	}

	job, _ := app.Jobs.Get(int(a.JobPostingID))
	jobTitle := ""
	if job != nil {
		jobTitle = job.Title
	}

	applicant, _ := app.Users.Get(int(a.ApplicantID))
	resp := mapApplicationResponse(a, jobTitle)
	if applicant != nil {
		resp.ApplicantName = applicant.FullName
		resp.ApplicantEmail = applicant.Email
	}

	oapi.SendResp(w, resp)
}

// GET /api/jobs/{JobID}/applications  — recruiter sees all applications for a job
func listJobApplications(w http.ResponseWriter, r *http.Request) {
	_, ok := recruiterUser(r)
	if !ok {
		oapi.Forbidden(w)
		return
	}

	jobID, err := strconv.Atoi(chi.URLParam(r, "JobID"))
	if err != nil || jobID <= 0 {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Ажлын байрны ID буруу байна"})
		return
	}

	job, err := app.Jobs.Get(jobID)
	if err != nil {
		oapi.CustomError(w, http.StatusNotFound, map[string]string{"message": "Ажлын байр олдсонгүй"})
		return
	}

	filter := &appman.Filter{Status: r.URL.Query().Get("status")}
	apps, err := app.Applications.ListForJob(jobID, filter)
	if err != nil {
		oapi.ServerError(w, err)
		return
	}

	items := make([]*applicationResponse, 0, len(apps))
	for _, a := range apps {
		resp := mapApplicationResponse(a, job.Title)
		applicant, _ := app.Users.Get(int(a.ApplicantID))
		if applicant != nil {
			resp.ApplicantName = applicant.FullName
			resp.ApplicantEmail = applicant.Email
		}
		items = append(items, resp)
	}

	oapi.SendResp(w, items)
}

// PUT /api/applications/{id}/status  — recruiter updates application status
func updateApplicationStatus(w http.ResponseWriter, r *http.Request) {
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
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Хүсэлт буруу байна"})
		return
	}

	validStatuses := map[string]bool{
		appman.StatusAssessed:    true,
		appman.StatusShortlisted: true,
		appman.StatusRejected:    true,
	}
	if !validStatuses[req.Status] {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Төлөв буруу байна"})
		return
	}

	existing, _ := app.Applications.Get(id)

	if err := app.Applications.UpdateStatus(id, req.Status); err != nil {
		oapi.ServerError(w, err)
		return
	}

	if existing != nil {
		statusLabels := map[string]string{
			appman.StatusShortlisted: "Шалгарсан",
			appman.StatusRejected:    "Татгалзсан",
			appman.StatusAssessed:    "Үнэлэгдсэн",
		}
		label := statusLabels[req.Status]
		socket.NotifyUser(int(existing.ApplicantID), "Анкетын төлөв өөрчлөгдлөө",
			fmt.Sprintf("Таны анкетын төлөв '%s' болов", label), "status_update")
	}

	oapi.SendResp(w, map[string]string{"message": "Төлөв шинэчлэгдлээ"})
}

// PUT /api/applications/{id}/interview  — recruiter sets interview date
func scheduleInterview(w http.ResponseWriter, r *http.Request) {
	recruiter, ok := recruiterUser(r)
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
		InterviewAt       string `json:"interview_at"`
		InterviewLocation string `json:"interview_location"`
		InterviewNote     string `json:"interview_note"`
		GenerateMeet      bool   `json:"generate_meet"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Хүсэлт буруу байна"})
		return
	}

	t, err := time.ParseInLocation("2006-01-02T15:04", req.InterviewAt, app.Location)
	if err != nil {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Огноо буруу байна (2006-01-02T15:04)"})
		return
	}

	a, err := app.Applications.Get(id)
	if err != nil {
		oapi.CustomError(w, http.StatusNotFound, map[string]string{"message": "Олдсонгүй"})
		return
	}

	job, _ := app.Jobs.Get(int(a.JobPostingID))
	jobTitle := ""
	if job != nil {
		jobTitle = job.Title
	}
	applicant, _ := app.Users.Get(int(a.ApplicantID))

	locationToStore := req.InterviewLocation
	gcalDone := false

	// When Google Meet is requested, create the Calendar event synchronously
	// so we can capture and store the Meet link before responding.
	if req.GenerateMeet && recruiter.GoogleRefreshToken != "" {
		desc := ""
		if applicant != nil {
			desc = fmt.Sprintf("Горилогч: %s", applicant.FullName)
		}
		if req.InterviewNote != "" {
			desc += "\n" + req.InterviewNote
		}
		applicantEmail := ""
		if applicant != nil {
			applicantEmail = applicant.Email
		}
		meetLink, gcalErr := createGCalEvent(r.Context(), recruiter, fmt.Sprintf("Ярилцлага: %s", jobTitle), desc, "", applicantEmail, t, true)
		if gcalErr != nil {
			app.ErrorLog.Printf("gcal meet creation failed: %v", gcalErr)
		} else if meetLink != "" {
			locationToStore = meetLink
		}
		gcalDone = true
	}

	a.InterviewAt = &t
	a.InterviewLocation = locationToStore
	a.InterviewNote = req.InterviewNote
	saved, err := app.Applications.Save(a)
	if err != nil {
		oapi.ServerError(w, err)
		return
	}

	resp := mapApplicationResponse(saved, jobTitle)
	if applicant != nil {
		resp.ApplicantName = applicant.FullName
		resp.ApplicantEmail = applicant.Email
	}

	if applicant != nil {
		socket.NotifyUser(int(saved.ApplicantID), "Ярилцлага товлогдлоо",
			fmt.Sprintf("%s — %s", jobTitle, t.Format("2006-01-02 15:04")), "interview_scheduled")
	}

	// Send email invite to applicant and confirmation to recruiter asynchronously
	if applicant != nil {
		go app.Mailer.SendInterviewInviteEmail(
			applicant.Email,
			applicant.FullName,
			recruiter.FullName,
			recruiter.Email,
			jobTitle,
			t,
			locationToStore,
			req.InterviewNote,
		)
		go app.Mailer.SendInterviewConfirmToRecruiter(
			recruiter.Email,
			recruiter.FullName,
			applicant.FullName,
			applicant.Email,
			jobTitle,
			t,
			locationToStore,
			req.InterviewNote,
		)
	}

	// Create Google Calendar event in background (onsite case, not already done above)
	if !gcalDone && recruiter.GoogleRefreshToken != "" {
		applicantEmailCopy := resp.ApplicantEmail
		go func() {
			desc := fmt.Sprintf("Горилогч: %s", resp.ApplicantName)
			if req.InterviewNote != "" {
				desc += "\n" + req.InterviewNote
			}
			if _, err := createGCalEvent(context.Background(), recruiter, fmt.Sprintf("Ярилцлага: %s", jobTitle), desc, req.InterviewLocation, applicantEmailCopy, t, false); err != nil {
				app.ErrorLog.Printf("gcal event creation failed: %v", err)
			}
		}()
	}

	oapi.SendResp(w, resp)
}

// GET /api/interviews  — recruiter sees all scheduled interviews for their company
func listInterviewsHandler(w http.ResponseWriter, r *http.Request) {
	recruiter, ok := recruiterUser(r)
	if !ok {
		oapi.Forbidden(w)
		return
	}

	if recruiter.CompanyID == nil {
		oapi.SendResp(w, []*applicationResponse{})
		return
	}

	apps, err := app.Applications.ListScheduledForCompany(*recruiter.CompanyID)
	if err != nil {
		oapi.ServerError(w, err)
		return
	}

	items := make([]*applicationResponse, 0, len(apps))
	for _, a := range apps {
		job, _ := app.Jobs.Get(int(a.JobPostingID))
		jobTitle := ""
		if job != nil {
			jobTitle = job.Title
		}
		resp := mapApplicationResponse(a, jobTitle)
		applicant, _ := app.Users.Get(int(a.ApplicantID))
		if applicant != nil {
			resp.ApplicantName = applicant.FullName
			resp.ApplicantEmail = applicant.Email
			resp.ApplicantProfileURL = applicant.ProfilePicture
		}
		items = append(items, resp)
	}

	oapi.SendResp(w, items)
}

// POST /api/jobs/{JobID}/analyze  — applicant gets AI assessment preview without saving
func analyzeCV(w http.ResponseWriter, r *http.Request) {
	user := authUser(r)
	if user.Role == userman.RoleRecruiter {
		oapi.Forbidden(w)
		return
	}

	jobID, err := strconv.Atoi(chi.URLParam(r, "JobID"))
	if err != nil || jobID <= 0 {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Ажлын байрны ID буруу байна"})
		return
	}

	job, err := app.Jobs.Get(jobID)
	if err != nil {
		oapi.CustomError(w, http.StatusNotFound, map[string]string{"message": "Ажлын байр олдсонгүй"})
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Файл уншихад алдаа гарлаа"})
		return
	}

	file, header, err := r.FormFile("cv")
	if err != nil {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "CV файл сонгоно уу"})
		return
	}
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

	cvText, err := aiman.ExtractTextFromPDF(pdfBytes)
	if err != nil || len(cvText) < 50 {
		cvText = "CV текст уншихад алдаа гарлаа - үнэлгээ хийх боломжгүй"
	}

	jobRequirements := make([]string, 0, len(job.Requirements))
	for _, r := range job.Requirements {
		jobRequirements = append(jobRequirements, r.Description)
	}
	jobSkills := make([]string, 0, len(job.Skills))
	for _, s := range job.Skills {
		jobSkills = append(jobSkills, s.Name)
	}
	jobDuties := make([]string, 0, len(job.Duties))
	for _, d := range job.Duties {
		jobDuties = append(jobDuties, d.Description)
	}

	result, err := app.AI.AssessCV(cvText, job.Title, jobRequirements, jobSkills, jobDuties)
	if err != nil {
		oapi.ServerError(w, err)
		return
	}

	oapi.SendResp(w, result)
}

// POST /api/jobs/{JobID}/apply-from-profile  — apply using saved CV profile (no file upload)
func applyFromProfile(w http.ResponseWriter, r *http.Request) {
	user := authUser(r)
	if user.Role == userman.RoleRecruiter {
		oapi.Forbidden(w)
		return
	}

	jobID, err := strconv.Atoi(chi.URLParam(r, "JobID"))
	if err != nil || jobID <= 0 {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Ажлын байрны ID буруу байна"})
		return
	}

	job, err := app.Jobs.Get(jobID)
	if err != nil {
		oapi.CustomError(w, http.StatusNotFound, map[string]string{"message": "Ажлын байр олдсонгүй"})
		return
	}

	cv, err := app.CVProfiles.GetByUserID(user.ID)
	if err != nil {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Хадгалагдсан CV олдсонгүй. Эхлээд CV бүрдүүлэгч хэсэгт мэдээллээ оруулна уу."})
		return
	}

	cvText := cvProfileToText(cv)
	if len(cvText) < 50 {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "CV мэдээлэл хангалтгүй байна. CV бүрдүүлэгч хэсэгт мэдээллээ нэмнэ үү."})
		return
	}

	// Delete existing application to allow re-apply
	if existing, err := app.Applications.GetForApplicantAndJob(user.ID, jobID); err == nil {
		if delErr := app.Applications.Delete(existing.ID); delErr != nil {
			oapi.ServerError(w, delErr)
			return
		}
		if existing.CVFilePath != "" {
			os.Remove(existing.CVFilePath)
		}
	}

	application := &appman.Application{
		JobPostingID: uint(jobID),
		ApplicantID:  uint(user.ID),
		CVFilePath:   "",
		CVText:       cvText,
		Status:       appman.StatusPending,
	}

	savedApp, err := app.Applications.Save(application)
	if err != nil {
		oapi.ServerError(w, err)
		return
	}

	jobRequirements := make([]string, 0, len(job.Requirements))
	for _, r := range job.Requirements {
		jobRequirements = append(jobRequirements, r.Description)
	}
	jobSkills := make([]string, 0, len(job.Skills))
	for _, s := range job.Skills {
		jobSkills = append(jobSkills, s.Name)
	}
	jobDuties := make([]string, 0, len(job.Duties))
	for _, d := range job.Duties {
		jobDuties = append(jobDuties, d.Description)
	}

	go func(appID int, cvText, jobTitle string, requirements, skills, duties []string) {
		result, err := app.AI.AssessCV(cvText, jobTitle, requirements, skills, duties)
		if err != nil {
			app.ErrorLog.Printf("AI assessment failed for application %d: %v", appID, err)
			return
		}

		matchedJSON, _ := json.Marshal(result.MatchedSkills)
		missingJSON, _ := json.Marshal(result.MissingSkills)
		recsJSON, _ := json.Marshal(result.Recommendations)
		dutyJSON, _ := json.Marshal(result.DutyAssessments)
		reqJSON, _ := json.Marshal(result.RequirementAssessments)
		now := time.Now()

		a, _ := app.Applications.Get(appID)
		if a == nil {
			return
		}
		a.OverallScore = result.OverallScore
		a.Summary = result.Summary
		a.MatchedSkills = string(matchedJSON)
		a.MissingSkills = string(missingJSON)
		a.Recommendations = string(recsJSON)
		a.DutyAssessments = string(dutyJSON)
		a.RequirementAssessments = string(reqJSON)
		a.Status = appman.StatusAssessed
		a.AssessedAt = &now
		app.Applications.Save(a)

		socket.NotifyUser(int(a.ApplicantID), "AI үнэлгээ дууслаа",
			fmt.Sprintf("'%s' ажлын байранд таны анкет %d оноо авлаа", jobTitle, result.OverallScore), "assessment_complete")
	}(savedApp.ID, cvText, job.Title, jobRequirements, jobSkills, jobDuties)

	w.WriteHeader(http.StatusCreated)
	oapi.SendResp(w, map[string]any{
		"message": "Анкет амжилттай илгээгдлээ. AI үнэлгээ хийгдэж байна...",
		"id":      savedApp.ID,
	})
}

// POST /api/jobs/{JobID}/analyze-from-profile  — AI preview using saved CV profile (no save)
func analyzeFromProfile(w http.ResponseWriter, r *http.Request) {
	user := authUser(r)
	if user.Role == userman.RoleRecruiter {
		oapi.Forbidden(w)
		return
	}

	jobID, err := strconv.Atoi(chi.URLParam(r, "JobID"))
	if err != nil || jobID <= 0 {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Ажлын байрны ID буруу байна"})
		return
	}

	job, err := app.Jobs.Get(jobID)
	if err != nil {
		oapi.CustomError(w, http.StatusNotFound, map[string]string{"message": "Ажлын байр олдсонгүй"})
		return
	}

	cv, err := app.CVProfiles.GetByUserID(user.ID)
	if err != nil {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Хадгалагдсан CV олдсонгүй. Эхлээд CV бүрдүүлэгч хэсэгт мэдээллээ оруулна уу."})
		return
	}

	cvText := cvProfileToText(cv)
	if len(cvText) < 50 {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "CV мэдээлэл хангалтгүй байна. CV бүрдүүлэгч хэсэгт мэдээллээ нэмнэ үү."})
		return
	}

	jobRequirements := make([]string, 0, len(job.Requirements))
	for _, r := range job.Requirements {
		jobRequirements = append(jobRequirements, r.Description)
	}
	jobSkills := make([]string, 0, len(job.Skills))
	for _, s := range job.Skills {
		jobSkills = append(jobSkills, s.Name)
	}
	jobDuties := make([]string, 0, len(job.Duties))
	for _, d := range job.Duties {
		jobDuties = append(jobDuties, d.Description)
	}

	result, err := app.AI.AssessCV(cvText, job.Title, jobRequirements, jobSkills, jobDuties)
	if err != nil {
		oapi.ServerError(w, err)
		return
	}

	oapi.SendResp(w, result)
}
