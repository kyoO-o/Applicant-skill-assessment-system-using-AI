package cvman

import (
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/entities"
)

type CVProfile struct {
	entities.Model
	UserID uint `json:"user_id" gorm:"uniqueIndex"`

	// General info
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	DateOfBirth    string `json:"date_of_birth"`
	Gender         string `json:"gender"`
	NationalID     string `json:"national_id"`
	DriverLicenses string `json:"driver_licenses" gorm:"type:text"` // JSON array e.g. ["A","B"]
	MaritalStatus  string `json:"marital_status"`

	// Contact
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Address string `json:"address"`

	// About
	About string `json:"about" gorm:"type:text"`

	// Complex sections stored as JSON text arrays
	WorkExperiences    string `json:"work_experiences" gorm:"type:text"`
	Education          string `json:"education" gorm:"type:text"`
	PersonalSkills     string `json:"personal_skills" gorm:"type:text"`
	ProfessionalSkills string `json:"professional_skills" gorm:"type:text"`
	Languages          string `json:"languages" gorm:"type:text"`
	ComputerSkills     string `json:"computer_skills" gorm:"type:text"`
	ArtSkills          string `json:"art_skills" gorm:"type:text"`
	SportSkills        string `json:"sport_skills" gorm:"type:text"`
	Trainings          string `json:"trainings" gorm:"type:text"`
	Exams              string `json:"exams" gorm:"type:text"`
	Internships        string `json:"internships" gorm:"type:text"`
	Awards             string `json:"awards" gorm:"type:text"`
}
