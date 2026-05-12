package appman

import (
	"errors"
	"log"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("application.not_found")

type Service struct {
	db       *gorm.DB
	infoLog  *log.Logger
	errorLog *log.Logger
}

func NewService(db *gorm.DB, infoLog, errorLog *log.Logger) *Service {
	return &Service{db: db, infoLog: infoLog, errorLog: errorLog}
}

func (s *Service) Save(app *Application) (*Application, error) {
	if err := s.db.Save(app).Error; err != nil {
		return nil, err
	}
	return app, nil
}

func (s *Service) Get(id int) (*Application, error) {
	var app Application
	if err := s.db.First(&app, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &app, nil
}

func (s *Service) ListForJob(jobID int) ([]*Application, error) {
	var apps []*Application
	if err := s.db.Where("job_posting_id = ?", jobID).
		Order("overall_score DESC, created_at DESC").
		Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps, nil
}

func (s *Service) ListForApplicant(applicantID int) ([]*Application, error) {
	var apps []*Application
	if err := s.db.Where("applicant_id = ?", applicantID).
		Order("created_at DESC").
		Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps, nil
}

func (s *Service) GetForApplicantAndJob(applicantID, jobID int) (*Application, error) {
	var app Application
	if err := s.db.Where("applicant_id = ? AND job_posting_id = ?", applicantID, jobID).
		First(&app).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &app, nil
}

func (s *Service) UpdateStatus(id int, status string) error {
	return s.db.Model(&Application{}).Where("id = ?", id).Update("status", status).Error
}
