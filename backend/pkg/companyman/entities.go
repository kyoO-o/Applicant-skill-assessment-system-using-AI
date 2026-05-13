package companyman

import "github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/entities"

type Company struct {
	entities.Model
	Name        string           `json:"name"`
	Description string           `json:"description" gorm:"type:text"`
	RegisterID  string           `json:"register_id"`
	ContactInfo string           `json:"contact_info"`
	City        string           `json:"city"`
	District    string           `json:"district"`
	LocationX   float64          `json:"location_x"`
	LocationY   float64          `json:"location_y"`
	Benefits    []CompanyBenefit `json:"benefits" gorm:"foreignKey:CompanyID"`
}

type CompanyBenefit struct {
	entities.Model
	CompanyID   uint   `json:"company_id"`
	Description string `json:"description"`
}
