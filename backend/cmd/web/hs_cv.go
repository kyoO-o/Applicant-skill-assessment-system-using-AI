package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"path/filepath"

	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/cmd/web/app"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/common/oapi"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/aiman"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/cvman"
)

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
