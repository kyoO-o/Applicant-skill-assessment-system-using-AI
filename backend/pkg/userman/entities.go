package userman

import (
	"time"

	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/entities"
)

const (
	RoleAdmin     = "admin"
	RoleUser      = "user"
	RoleRecruiter = "recruiter"
)

type User struct {
	entities.Model
	FirstName           string     `json:"first_name"`
	LastName            string     `json:"last_name"`
	FullName            string     `json:"full_name"`
	Email               string     `json:"email" gorm:"index:idx_user_email"`
	PasswordHash        string     `json:"-" gorm:"column:password_hash"`
	PhoneNumber         string     `json:"phone_number"`
	ToduID              int        `json:"todu_id"`
	Role                string     `json:"role"`
	CompanyID           *uint      `json:"company_id"`
	ProfilePicture      string     `json:"profile_url" gorm:"column:profile_picture"`
	EmailVerified       bool       `json:"email_verified" gorm:"default:false"`
	VerifyCode          string     `json:"-" gorm:"column:verify_code"`
	VerifyCodeExpiry    *time.Time `json:"-" gorm:"column:verify_code_expiry"`
	PendingEmail        string     `json:"-" gorm:"column:pending_email"`
	ResetCode           string     `json:"-" gorm:"column:reset_code"`
	ResetCodeExpiry     *time.Time `json:"-" gorm:"column:reset_code_expiry"`
	GoogleRefreshToken  string     `json:"-" gorm:"column:google_refresh_token"`
}
