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
	Email        string `json:"email" gorm:"index:idx_user_email"`
	ProfileURL   string `json:"profile_url"`
	Name         string `json:"name"`
	PhoneNumber  string `json:"phone_number"`
	ToduID       int    `json:"todu_id"`
	PasswordHash string `json:"-" gorm:"column:password_hash"`
	Role         string `json:"role"`
}
