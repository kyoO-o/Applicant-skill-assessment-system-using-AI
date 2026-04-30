package userman

import (
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/entities"
)

const (
	RoleAdmin     = "admin"
	RoleUser      = "user"
	RoleRecruiter = "recruiter"
)

type User struct {
	entities.Model
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	FullName     string `json:"full_name"`
	Email        string `json:"email" gorm:"index:idx_user_email"`
	PasswordHash string `json:"-" gorm:"column:password_hash"`
	PhoneNumber  string `json:"phone_number"`
	ToduID       int    `json:"todu_id"`
	Role         string `json:"role"`
	CompanyID    *uint  `json:"company_id"` // Foreign key to Company, nullable
}
