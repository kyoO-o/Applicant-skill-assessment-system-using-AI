package appman

import (
	"time"

	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/entities"
)

const (
	StatusPending    = "pending"
	StatusAssessed   = "assessed"
	StatusShortlisted = "shortlisted"
	StatusRejected   = "rejected"
)

type SkillResult struct {
	Skill       string `json:"skill"`
	Explanation string `json:"explanation"`
}

type Application struct {
	entities.Model
	JobPostingID    uint      `json:"job_posting_id" gorm:"index"`
	ApplicantID     uint      `json:"applicant_id" gorm:"index"`
	CVFilePath      string    `json:"cv_file_path"`
	CVText          string    `json:"-" gorm:"type:text"`
	OverallScore    int       `json:"overall_score"`
	Summary         string    `json:"summary" gorm:"type:text"`
	MatchedSkills          string `json:"matched_skills" gorm:"type:text"`
	MissingSkills          string `json:"missing_skills" gorm:"type:text"`
	Recommendations        string `json:"recommendations" gorm:"type:text"`
	DutyAssessments        string `json:"duty_assessments" gorm:"type:text"`
	RequirementAssessments string `json:"requirement_assessments" gorm:"type:text"`
	Status           string     `json:"status" gorm:"default:pending"`
	AssessedAt       *time.Time `json:"assessed_at"`
	InterviewAt      *time.Time `json:"interview_at"`
	InterviewLocation string    `json:"interview_location" gorm:"type:text"`
	InterviewNote     string    `json:"interview_note" gorm:"type:text"`

	ApplicantName  string `json:"applicant_name" gorm:"-"`
	ApplicantEmail string `json:"applicant_email" gorm:"-"`
}
