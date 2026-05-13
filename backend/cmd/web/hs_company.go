package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/cmd/web/app"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/common/oapi"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/companyman"
	"gorm.io/gorm"
)

type companyRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	RegisterID  string   `json:"register_id"`
	ContactInfo string   `json:"contact_info"`
	City        string   `json:"city"`
	District    string   `json:"district"`
	LocationX   *float64 `json:"location_x"`
	LocationY   *float64 `json:"location_y"`
	Benefits    []string `json:"benefits"`
}

func chosenCompany(r *http.Request) (*companyman.Company, bool) {
	company, ok := r.Context().Value(app.ContextKeyChosenCompany).(*companyman.Company)
	return company, ok
}

func GetCompany(w http.ResponseWriter, r *http.Request) {
	user := authUser(r)
	if user.CompanyID == nil {
		oapi.CustomError(w, http.StatusNotFound, map[string]string{"message": "Company not found"})
		return
	}

	var company companyman.Company
	if err := app.DB.Preload("Benefits").First(&company, *user.CompanyID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			oapi.CustomError(w, http.StatusNotFound, map[string]string{"message": "Company not found"})
			return
		}
		oapi.ServerError(w, err)
		return
	}

	oapi.SendResp(w, &company)
}

func SaveCompany(w http.ResponseWriter, r *http.Request) {
	user, ok := recruiterUser(r)
	if !ok {
		oapi.Forbidden(w)
		return
	}

	var req companyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)
	req.RegisterID = strings.TrimSpace(req.RegisterID)
	req.ContactInfo = strings.TrimSpace(req.ContactInfo)
	req.City = strings.TrimSpace(req.City)
	req.District = strings.TrimSpace(req.District)

	if req.Name == "" {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Company name is required"})
		return
	}

	var company companyman.Company
	isNew := user.CompanyID == nil
	if !isNew {
		if err := app.DB.First(&company, *user.CompanyID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				isNew = true
			} else {
				oapi.ServerError(w, err)
				return
			}
		}
	}

	company.Name = req.Name
	company.Description = req.Description
	company.RegisterID = req.RegisterID
	company.ContactInfo = req.ContactInfo
	company.City = req.City
	company.District = req.District
	if req.LocationX != nil {
		company.LocationX = *req.LocationX
	}
	if req.LocationY != nil {
		company.LocationY = *req.LocationY
	}

	if err := app.DB.Save(&company).Error; err != nil {
		oapi.ServerError(w, err)
		return
	}

	// Replace benefits
	app.DB.Where("company_id = ?", company.ID).Delete(&companyman.CompanyBenefit{})
	for _, desc := range req.Benefits {
		desc = strings.TrimSpace(desc)
		if desc == "" {
			continue
		}
		app.DB.Create(&companyman.CompanyBenefit{CompanyID: uint(company.ID), Description: desc})
	}
	app.DB.Preload("Benefits").First(&company, company.ID)

	if user.CompanyID == nil || *user.CompanyID != uint(company.ID) {
		companyID := uint(company.ID)
		user.CompanyID = &companyID
		if _, err := app.Users.Save(user); err != nil {
			oapi.ServerError(w, err)
			return
		}
	}

	if isNew {
		w.WriteHeader(http.StatusCreated)
	}
	oapi.SendResp(w, &company)
}
