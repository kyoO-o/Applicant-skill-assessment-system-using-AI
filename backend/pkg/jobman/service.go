package jobman

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

type Service struct {
	DB       *gorm.DB
	infoLog  *log.Logger
	errorLog *log.Logger
}

func NewService(db *gorm.DB, infoLog, errorLog *log.Logger) *Service {
	return &Service{
		DB:       db,
		infoLog:  infoLog,
		errorLog: errorLog,
	}
}

func (s *Service) parseFilter(filter *Filter) *gorm.DB {
	query := s.DB
	if filter == nil {
		return query
	}
	if filter.Keyword != "" {
		query = query.Where("title ILIKE ? || '%%'", filter.Keyword)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	return query
}

func (s *Service) ListForRecruiter(recruiterID int, filter *Filter) ([]*JobPosting, error) {
	var jobs []*JobPosting
	if err := s.parseFilter(filter).
		Preload("Duties").
		Preload("Requirements").
		Preload("Skills").
		Preload("Bonuses").
		Where("posted_by = ?", recruiterID).
		Order("updated_at DESC, created_at DESC").
		Find(&jobs).Error; err != nil {
		return nil, err
	}
	s.fillAppsCounts(jobs)
	return jobs, nil
}

func (s *Service) ListForCompany(companyID int, filter *Filter) ([]*JobPosting, error) {
	var jobs []*JobPosting
	if err := s.parseFilter(filter).
		Preload("Duties").
		Preload("Requirements").
		Preload("Skills").
		Preload("Bonuses").
		Where("company_id = ?", companyID).
		Order("updated_at DESC, created_at DESC").
		Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

func (s *Service) ListForRecruiterAndCompany(recruiterID, companyID int, filter *Filter) ([]*JobPosting, error) {
	var jobs []*JobPosting
	if err := s.parseFilter(filter).
		Preload("Duties").
		Preload("Requirements").
		Preload("Skills").
		Preload("Bonuses").
		Where("posted_by = ? AND company_id = ?", recruiterID, companyID).
		Order("updated_at DESC, created_at DESC").
		Find(&jobs).Error; err != nil {
		return nil, err
	}
	s.fillAppsCounts(jobs)
	return jobs, nil
}

func (s *Service) ListActive(filter *Filter) ([]*JobPosting, error) {
	var jobs []*JobPosting
	if err := s.parseFilter(filter).
		Preload("Duties").
		Preload("Requirements").
		Preload("Skills").
		Preload("Bonuses").
		Where("status = ?", StatusPosted).
		Order("updated_at DESC, created_at DESC").
		Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

func (s *Service) Get(id int) (*JobPosting, error) {
	var job *JobPosting
	if err := s.DB.
		Preload("Duties").
		Preload("Requirements").
		Preload("Skills").
		Preload("Bonuses").
		First(&job, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return job, nil
}

func (s *Service) GetForRecruiter(id, recruiterID int) (*JobPosting, error) {
	var job *JobPosting
	if err := s.DB.
		Preload("Duties").
		Preload("Requirements").
		Preload("Skills").
		Preload("Bonuses").
		Where("id = ? AND posted_by = ?", id, recruiterID).
		First(&job).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return job, nil
}

func (s *Service) Save(job *JobPosting) (*JobPosting, error) {
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Requirements", "Duties", "Skills", "Bonuses").Save(job).Error; err != nil {
			return err
		}

		if err := tx.Where("job_posting_id = ?", job.ID).Delete(new(Duty)).Error; err != nil {
			return err
		}
		if err := tx.Where("job_posting_id = ?", job.ID).Delete(new(Requirement)).Error; err != nil {
			return err
		}
		if err := tx.Where("job_posting_id = ?", job.ID).Delete(new(Skill)).Error; err != nil {
			return err
		}
		if err := tx.Where("job_posting_id = ?", job.ID).Delete(new(Bonus)).Error; err != nil {
			return err
		}

		if len(job.Duties) > 0 {
			duties := make([]Duty, 0, len(job.Duties))
			for _, duty := range job.Duties {
				duties = append(duties, Duty{
					JobPostingID: uint(job.ID),
					Description:  duty.Description,
				})
			}
			if err := tx.Create(&duties).Error; err != nil {
				return err
			}
		}

		if len(job.Requirements) > 0 {
			requirements := make([]Requirement, 0, len(job.Requirements))
			for _, requirement := range job.Requirements {
				requirements = append(requirements, Requirement{
					JobPostingID: uint(job.ID),
					Description:  requirement.Description,
				})
			}
			if err := tx.Create(&requirements).Error; err != nil {
				return err
			}
		}

		if len(job.Skills) > 0 {
			skills := make([]Skill, 0, len(job.Skills))
			for _, skill := range job.Skills {
				skills = append(skills, Skill{
					JobPostingID: uint(job.ID),
					Name:         skill.Name,
				})
			}
			if err := tx.Create(&skills).Error; err != nil {
				return err
			}
		}

		if len(job.Bonuses) > 0 {
			bonuses := make([]Bonus, 0, len(job.Bonuses))
			for _, bonus := range job.Bonuses {
				bonuses = append(bonuses, Bonus{
					JobPostingID: uint(job.ID),
					Description:  bonus.Description,
				})
			}
			if err := tx.Create(&bonuses).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.Get(job.ID)
}

type appCountRow struct {
	JobPostingID int
	Count        int
}

func (s *Service) fillAppsCounts(jobs []*JobPosting) {
	if len(jobs) == 0 {
		return
	}
	ids := make([]int, len(jobs))
	for i, j := range jobs {
		ids[i] = j.ID
	}

	var rows []appCountRow
	s.DB.Raw(
		"SELECT job_posting_id, COUNT(*) as count FROM applications WHERE job_posting_id IN ? AND deleted_at IS NULL GROUP BY job_posting_id",
		ids,
	).Scan(&rows)
	countMap := make(map[int]int, len(rows))
	for _, r := range rows {
		countMap[r.JobPostingID] = r.Count
	}

	weekAgo := time.Now().AddDate(0, 0, -7)
	var weekRows []appCountRow
	s.DB.Raw(
		"SELECT job_posting_id, COUNT(*) as count FROM applications WHERE job_posting_id IN ? AND deleted_at IS NULL AND created_at >= ? GROUP BY job_posting_id",
		ids, weekAgo,
	).Scan(&weekRows)
	weekMap := make(map[int]int, len(weekRows))
	for _, r := range weekRows {
		weekMap[r.JobPostingID] = r.Count
	}

	for _, j := range jobs {
		j.AppsCount = countMap[j.ID]
		j.NewAppsThisWeek = weekMap[j.ID]
	}
}

func (s *Service) Delete(id int) error {
	if err := s.DB.Delete(new(JobPosting), id).Error; err != nil {
		return err
	}
	return nil
}
