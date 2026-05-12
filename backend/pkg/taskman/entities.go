package taskman

import (
	"time"

	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/entities"
)

const (
	TaskStatusDraft     = "draft"
	TaskStatusSent      = "sent"
	TaskStatusCompleted = "completed"
	TaskStatusGraded    = "graded"

	SubStatusSubmitted = "submitted"
	SubStatusGraded    = "graded"
)

type Task struct {
	entities.Model
	JobPostingID  uint       `json:"job_posting_id" gorm:"index"`
	ApplicationID *uint      `json:"application_id" gorm:"index"`
	Title         string     `json:"title"`
	Description   string     `json:"description" gorm:"type:text"`
	DueDate       *time.Time `json:"due_date"`
	Status        string     `json:"status" gorm:"default:draft"`
	CreatedByAI   bool       `json:"created_by_ai"`

	Submissions []TaskSubmission `json:"submissions,omitempty" gorm:"foreignKey:TaskID"`
}

type TaskSubmission struct {
	entities.Model
	TaskID      uint   `json:"task_id" gorm:"index"`
	ApplicantID uint   `json:"applicant_id" gorm:"index"`
	Content     string `json:"content" gorm:"type:text"`
	FilePath    string `json:"file_path"`
	Grade       *int   `json:"grade"`
	Feedback    string `json:"feedback" gorm:"type:text"`
	Status      string `json:"status" gorm:"default:submitted"`
}
