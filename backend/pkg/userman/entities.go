package userman

import (
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/entities"
)

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type User struct {
	entities.Model
	Email       string `json:"email" gorm:"index:idx_user_email"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	ToduID      int    `json:"todu_id"`
	Role        string `json:"role"`
}
