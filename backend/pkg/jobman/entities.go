package jobman

import (
	"time"

	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/entities"
)

const (
	StatusDraft  = "draft"
	StatusPosted = "posted"
	StatusClosed = "closed"
)

type JobPosting struct {
	entities.Model
	Title          string    `json:"title"`
	Location       string    `json:"location"`
	AdditionalInfo string    `json:"additional_info" gorm:"type:text"`
	ContactInfo    string    `json:"contact_info"`
	PostedAt       time.Time `json:"posted_at"`
	Level          string    `json:"level"`      // professional, employee, intern
	Type           string    `json:"type"`       // full-time, part-time, shift, hourly
	City           *string   `json:"city"`       // optional, default to company
	District       *string   `json:"district"`   // optional, default to company
	LocationX      *float64  `json:"location_x"` // optional
	LocationY      *float64  `json:"location_y"` // optional
	SeenCount      int       `json:"seen_count"`
	AppsCount      int       `json:"apps_count"`
	MinSalary      float64   `json:"min_salary"`
	MaxSalary      float64   `json:"max_salary"`
	Status         string    `json:"status" gorm:"default:draft"`
	CompanyID      uint      `json:"company_id"` // Foreign key to Company
	PostedBy       uint      `json:"posted_by"`  // Foreign key to User
	// Relations
	Duties       []Duty        `json:"duties" gorm:"foreignKey:JobPostingID"`
	Requirements []Requirement `json:"requirements" gorm:"foreignKey:JobPostingID"`
	Skills       []Skill       `json:"skills" gorm:"foreignKey:JobPostingID"`
	Bonuses      []Bonus       `json:"bonuses" gorm:"foreignKey:JobPostingID"`
}

type Duty struct {
	entities.Model
	JobPostingID uint   `json:"job_posting_id"`
	Description  string `json:"description"`
}

type Requirement struct {
	entities.Model
	JobPostingID uint   `json:"job_posting_id"`
	Description  string `json:"description"`
}

type Skill struct {
	entities.Model
	JobPostingID uint   `json:"job_posting_id"`
	Name         string `json:"name"`
}

type Bonus struct {
	entities.Model
	JobPostingID uint   `json:"job_posting_id"`
	Description  string `json:"description"`
}
