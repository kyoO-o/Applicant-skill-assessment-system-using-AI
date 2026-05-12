package taskman

import (
	"errors"
	"log"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("task.not_found")

type Service struct {
	db       *gorm.DB
	infoLog  *log.Logger
	errorLog *log.Logger
}

func NewService(db *gorm.DB, infoLog, errorLog *log.Logger) *Service {
	return &Service{db: db, infoLog: infoLog, errorLog: errorLog}
}

func (s *Service) Save(task *Task) (*Task, error) {
	if err := s.db.Save(task).Error; err != nil {
		return nil, err
	}
	return task, nil
}

func (s *Service) Get(id int) (*Task, error) {
	var task Task
	if err := s.db.Preload("Submissions").First(&task, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &task, nil
}

func (s *Service) ListForJob(jobID int) ([]*Task, error) {
	var tasks []*Task
	if err := s.db.Preload("Submissions").
		Where("job_posting_id = ?", jobID).
		Order("created_at DESC").
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *Service) ListForApplicant(applicantID int) ([]*Task, error) {
	var tasks []*Task
	if err := s.db.
		Joins("JOIN task_submissions ON task_submissions.task_id = tasks.id AND task_submissions.applicant_id = ?", applicantID).
		Where("tasks.status = ?", TaskStatusSent).
		Order("tasks.created_at DESC").
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *Service) GetSentToApplicant(taskID, applicantID int) (*Task, error) {
	var task Task
	if err := s.db.First(&task, "id = ? AND status = ?", taskID, TaskStatusSent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &task, nil
}

func (s *Service) Delete(id int) error {
	return s.db.Delete(&Task{}, id).Error
}

func (s *Service) SaveSubmission(sub *TaskSubmission) (*TaskSubmission, error) {
	if err := s.db.Save(sub).Error; err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *Service) GetSubmission(id int) (*TaskSubmission, error) {
	var sub TaskSubmission
	if err := s.db.First(&sub, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &sub, nil
}

func (s *Service) ListSubmissionsForTask(taskID int) ([]*TaskSubmission, error) {
	var subs []*TaskSubmission
	if err := s.db.Where("task_id = ?", taskID).Order("created_at DESC").Find(&subs).Error; err != nil {
		return nil, err
	}
	return subs, nil
}
