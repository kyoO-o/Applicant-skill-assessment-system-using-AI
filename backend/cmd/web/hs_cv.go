package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/cmd/web/app"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/common/oapi"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/aiman"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/cvman"
)

// cvProfileToText converts a saved CV profile into a plain-text representation
// suitable for AI assessment (same role as PDF-extracted text).
func cvProfileToText(cv *cvman.CVProfile) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Name: %s %s\n", cv.FirstName, cv.LastName))
	if cv.DateOfBirth != "" {
		b.WriteString(fmt.Sprintf("Date of Birth: %s\n", cv.DateOfBirth))
	}
	if cv.Gender != "" {
		b.WriteString(fmt.Sprintf("Gender: %s\n", cv.Gender))
	}
	if cv.MaritalStatus != "" {
		b.WriteString(fmt.Sprintf("Marital Status: %s\n", cv.MaritalStatus))
	}
	if cv.Phone != "" {
		b.WriteString(fmt.Sprintf("Phone: %s\n", cv.Phone))
	}
	if cv.Email != "" {
		b.WriteString(fmt.Sprintf("Email: %s\n", cv.Email))
	}
	if cv.Address != "" {
		b.WriteString(fmt.Sprintf("Address: %s\n", cv.Address))
	}

	var licenses []string
	if cv.DriverLicenses != "" {
		json.Unmarshal([]byte(cv.DriverLicenses), &licenses)
	}
	if len(licenses) > 0 {
		b.WriteString(fmt.Sprintf("Driver Licenses: %s\n", strings.Join(licenses, ", ")))
	}

	if cv.About != "" {
		b.WriteString(fmt.Sprintf("\nAbout:\n%s\n", cv.About))
	}

	type workExp struct {
		Company   string `json:"company"`
		Position  string `json:"position"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Current   bool   `json:"current"`
		Desc      string `json:"description"`
	}
	var works []workExp
	if cv.WorkExperiences != "" {
		json.Unmarshal([]byte(cv.WorkExperiences), &works)
	}
	if len(works) > 0 {
		b.WriteString("\nWork Experience:\n")
		for _, w := range works {
			end := w.EndDate
			if w.Current {
				end = "Present"
			}
			b.WriteString(fmt.Sprintf("- %s at %s (%s – %s)\n", w.Position, w.Company, w.StartDate, end))
			if w.Desc != "" {
				b.WriteString(fmt.Sprintf("  %s\n", w.Desc))
			}
		}
	}

	type edu struct {
		School    string `json:"school"`
		Degree    string `json:"degree"`
		Field     string `json:"field"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Current   bool   `json:"current"`
		GPA       string `json:"gpa"`
	}
	var edus []edu
	if cv.Education != "" {
		json.Unmarshal([]byte(cv.Education), &edus)
	}
	if len(edus) > 0 {
		b.WriteString("\nEducation:\n")
		for _, e := range edus {
			end := e.EndDate
			if e.Current {
				end = "Present"
			}
			b.WriteString(fmt.Sprintf("- %s in %s, %s (%s – %s)", e.Degree, e.Field, e.School, e.StartDate, end))
			if e.GPA != "" {
				b.WriteString(fmt.Sprintf(", GPA: %s", e.GPA))
			}
			b.WriteString("\n")
		}
	}

	type lang struct {
		Name  string `json:"name"`
		Level string `json:"level"`
	}
	var langs []lang
	if cv.Languages != "" {
		json.Unmarshal([]byte(cv.Languages), &langs)
	}
	if len(langs) > 0 {
		b.WriteString("\nLanguages:\n")
		for _, l := range langs {
			b.WriteString(fmt.Sprintf("- %s (%s)\n", l.Name, l.Level))
		}
	}

	writeSkills := func(label, raw string) {
		var skills []string
		if raw != "" {
			json.Unmarshal([]byte(raw), &skills)
		}
		if len(skills) > 0 {
			b.WriteString(fmt.Sprintf("\n%s:\n%s\n", label, strings.Join(skills, ", ")))
		}
	}
	writeSkills("Professional Skills", cv.ProfessionalSkills)
	writeSkills("Personal Skills", cv.PersonalSkills)
	writeSkills("Computer Skills", cv.ComputerSkills)
	writeSkills("Art Skills", cv.ArtSkills)
	writeSkills("Sport Skills", cv.SportSkills)

	type training struct {
		Name         string `json:"name"`
		Organization string `json:"organization"`
		Date         string `json:"date"`
		Certificate  string `json:"certificate"`
	}
	var trainings []training
	if cv.Trainings != "" {
		json.Unmarshal([]byte(cv.Trainings), &trainings)
	}
	if len(trainings) > 0 {
		b.WriteString("\nTrainings & Certifications:\n")
		for _, t := range trainings {
			b.WriteString(fmt.Sprintf("- %s, %s (%s)", t.Name, t.Organization, t.Date))
			if t.Certificate != "" {
				b.WriteString(fmt.Sprintf(" — %s", t.Certificate))
			}
			b.WriteString("\n")
		}
	}

	type exam struct {
		Name  string `json:"name"`
		Score string `json:"score"`
		Date  string `json:"date"`
	}
	var exams []exam
	if cv.Exams != "" {
		json.Unmarshal([]byte(cv.Exams), &exams)
	}
	if len(exams) > 0 {
		b.WriteString("\nExams:\n")
		for _, e := range exams {
			b.WriteString(fmt.Sprintf("- %s: %s (%s)\n", e.Name, e.Score, e.Date))
		}
	}

	type internship struct {
		Company   string `json:"company"`
		Position  string `json:"position"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Desc      string `json:"description"`
	}
	var internships []internship
	if cv.Internships != "" {
		json.Unmarshal([]byte(cv.Internships), &internships)
	}
	if len(internships) > 0 {
		b.WriteString("\nInternships:\n")
		for _, i := range internships {
			b.WriteString(fmt.Sprintf("- %s at %s (%s – %s)\n", i.Position, i.Company, i.StartDate, i.EndDate))
			if i.Desc != "" {
				b.WriteString(fmt.Sprintf("  %s\n", i.Desc))
			}
		}
	}

	type award struct {
		Name         string `json:"name"`
		Organization string `json:"organization"`
		Date         string `json:"date"`
		Desc         string `json:"description"`
	}
	var awards []award
	if cv.Awards != "" {
		json.Unmarshal([]byte(cv.Awards), &awards)
	}
	if len(awards) > 0 {
		b.WriteString("\nAwards:\n")
		for _, a := range awards {
			b.WriteString(fmt.Sprintf("- %s, %s (%s)\n", a.Name, a.Organization, a.Date))
			if a.Desc != "" {
				b.WriteString(fmt.Sprintf("  %s\n", a.Desc))
			}
		}
	}

	return b.String()
}

func parseCVField(s string) any {
	if s == "" {
		return []any{}
	}
	var v any
	if json.Unmarshal([]byte(s), &v) == nil {
		return v
	}
	return []any{}
}

func mapCVResponse(cv *cvman.CVProfile) map[string]any {
	return map[string]any{
		"id":                  cv.ID,
		"first_name":          cv.FirstName,
		"last_name":           cv.LastName,
		"date_of_birth":       cv.DateOfBirth,
		"gender":              cv.Gender,
		"national_id":         cv.NationalID,
		"driver_licenses":     parseCVField(cv.DriverLicenses),
		"marital_status":      cv.MaritalStatus,
		"phone":               cv.Phone,
		"email":               cv.Email,
		"address":             cv.Address,
		"about":               cv.About,
		"work_experiences":    parseCVField(cv.WorkExperiences),
		"education":           parseCVField(cv.Education),
		"personal_skills":     parseCVField(cv.PersonalSkills),
		"professional_skills": parseCVField(cv.ProfessionalSkills),
		"languages":           parseCVField(cv.Languages),
		"computer_skills":     parseCVField(cv.ComputerSkills),
		"art_skills":          parseCVField(cv.ArtSkills),
		"sport_skills":        parseCVField(cv.SportSkills),
		"trainings":           parseCVField(cv.Trainings),
		"exams":               parseCVField(cv.Exams),
		"internships":         parseCVField(cv.Internships),
		"awards":              parseCVField(cv.Awards),
	}
}

func emptyCVResponse() map[string]any {
	return map[string]any{
		"id":                  0,
		"first_name":          "",
		"last_name":           "",
		"date_of_birth":       "",
		"gender":              "",
		"national_id":         "",
		"driver_licenses":     []any{},
		"marital_status":      "",
		"phone":               "",
		"email":               "",
		"address":             "",
		"about":               "",
		"work_experiences":    []any{},
		"education":           []any{},
		"personal_skills":     []any{},
		"professional_skills": []any{},
		"languages":           []any{},
		"computer_skills":     []any{},
		"art_skills":          []any{},
		"sport_skills":        []any{},
		"trainings":           []any{},
		"exams":               []any{},
		"internships":         []any{},
		"awards":              []any{},
	}
}

// GET /api/cv-profile
func getCVProfile(w http.ResponseWriter, r *http.Request) {
	user := authUser(r)

	cv, err := app.CVProfiles.GetByUserID(user.ID)
	if err != nil {
		if errors.Is(err, cvman.ErrNotFound) {
			oapi.SendResp(w, emptyCVResponse())
			return
		}
		oapi.ServerError(w, err)
		return
	}

	oapi.SendResp(w, mapCVResponse(cv))
}

// PUT /api/cv-profile
func saveCVProfile(w http.ResponseWriter, r *http.Request) {
	user := authUser(r)

	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Хүсэлт буруу байна"})
		return
	}

	cv, err := app.CVProfiles.GetByUserID(user.ID)
	if err != nil {
		if !errors.Is(err, cvman.ErrNotFound) {
			oapi.ServerError(w, err)
			return
		}
		cv = &cvman.CVProfile{UserID: uint(user.ID)}
	}

	str := func(key string) string {
		if v, ok := body[key]; ok {
			if s, ok := v.(string); ok {
				return s
			}
		}
		return ""
	}
	jsonStr := func(key string) string {
		if v, ok := body[key]; ok {
			b, _ := json.Marshal(v)
			return string(b)
		}
		return ""
	}

	cv.FirstName = str("first_name")
	cv.LastName = str("last_name")
	cv.DateOfBirth = str("date_of_birth")
	cv.Gender = str("gender")
	cv.NationalID = str("national_id")
	cv.DriverLicenses = jsonStr("driver_licenses")
	cv.MaritalStatus = str("marital_status")
	cv.Phone = str("phone")
	cv.Email = str("email")
	cv.Address = str("address")
	cv.About = str("about")
	cv.WorkExperiences = jsonStr("work_experiences")
	cv.Education = jsonStr("education")
	cv.PersonalSkills = jsonStr("personal_skills")
	cv.ProfessionalSkills = jsonStr("professional_skills")
	cv.Languages = jsonStr("languages")
	cv.ComputerSkills = jsonStr("computer_skills")
	cv.ArtSkills = jsonStr("art_skills")
	cv.SportSkills = jsonStr("sport_skills")
	cv.Trainings = jsonStr("trainings")
	cv.Exams = jsonStr("exams")
	cv.Internships = jsonStr("internships")
	cv.Awards = jsonStr("awards")

	saved, err := app.CVProfiles.Save(cv)
	if err != nil {
		oapi.ServerError(w, err)
		return
	}

	oapi.SendResp(w, mapCVResponse(saved))
}

// POST /api/cv-profile/parse — upload PDF, AI parse, return structured data (not saved)
func parseCVProfile(w http.ResponseWriter, r *http.Request) {
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
	if err != nil || len(cvText) < 20 {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "PDF-ээс текст уншихад алдаа гарлаа"})
		return
	}

	jsonStr, err := app.AI.ParseCV(cvText)
	if err != nil {
		oapi.ServerError(w, err)
		return
	}

	var result any
	json.Unmarshal([]byte(jsonStr), &result)
	oapi.SendResp(w, result)
}
