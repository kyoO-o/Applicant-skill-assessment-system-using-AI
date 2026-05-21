package appman

import (
	"errors"
	"log"
	"time"

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

func (s *Service) parseFilter(filter *Filter) *gorm.DB {
	query := s.db
	if filter == nil {
		return query
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	return query
}

func (s *Service) Save(app *Application) (*Application, error) {
	isNew := app.ID == 0
	if err := s.db.Save(app).Error; err != nil {
		return nil, err
	}
	if isNew {
		s.db.Exec("UPDATE job_postings SET apps_count = apps_count + 1 WHERE id = ?", app.JobPostingID)
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

func (s *Service) ListForJob(jobID int, filter *Filter) ([]*Application, error) {
	var apps []*Application
	if err := s.parseFilter(filter).Where("job_posting_id = ?", jobID).
		Order("overall_score DESC, created_at DESC").
		Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps, nil
}

func (s *Service) ListForApplicant(applicantID int, filter *Filter) ([]*Application, error) {
	var apps []*Application
	if err := s.parseFilter(filter).Where("applicant_id = ?", applicantID).
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

func (s *Service) Delete(id int) error {
	return s.db.Delete(&Application{}, id).Error
}

func (s *Service) UpdateStatus(id int, status string) error {
	return s.db.Model(&Application{}).Where("id = ?", id).Update("status", status).Error
}

func (s *Service) FindUpcomingInterviews(from, to time.Time) ([]*Application, error) {
	var apps []*Application
	if err := s.db.Where("interview_at BETWEEN ? AND ?", from, to).Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps, nil
}

func (s *Service) ListScheduledForCompany(companyID uint) ([]*Application, error) {
	var apps []*Application
	if err := s.db.
		Joins("JOIN job_postings ON job_postings.id = applications.job_posting_id").
		Where("job_postings.company_id = ? AND applications.interview_at IS NOT NULL", companyID).
		Order("applications.interview_at ASC").
		Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps, nil
}
