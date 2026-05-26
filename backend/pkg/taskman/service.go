package taskman

import (
	"errors"
	"log"
	"time"

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

func (s *Service) parseFilter(filter *Filter) *gorm.DB {
	query := s.db
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

func (s *Service) ListForJob(jobID int, filter *Filter) ([]*Task, error) {
	var tasks []*Task
	if err := s.parseFilter(filter).Preload("Submissions").
		Where("job_posting_id = ?", jobID).
		Order("created_at DESC").
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *Service) ListLibrary(filter *Filter) ([]*Task, error) {
	var tasks []*Task
	if err := s.parseFilter(filter).Preload("Submissions").
		Order("created_at DESC").
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *Service) ListForApplicant(applicantID int) ([]*Task, error) {
	var tasks []*Task
	if err := s.db.
		Joins("JOIN applications ON applications.id = tasks.application_id AND applications.applicant_id = ?", applicantID).
		Where("tasks.status IN ?", []string{TaskStatusSent, TaskStatusCompleted, TaskStatusGraded}).
		Order("tasks.created_at DESC").
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return tasks, nil
	}

	taskIDs := make([]int, len(tasks))
	for i, t := range tasks {
		taskIDs[i] = t.ID
	}

	var subs []TaskSubmission
	s.db.Where("task_id IN ? AND applicant_id = ?", taskIDs, applicantID).Find(&subs)
	subMap := make(map[int][]TaskSubmission)
	for _, sub := range subs {
		subMap[int(sub.TaskID)] = append(subMap[int(sub.TaskID)], sub)
	}
	for _, t := range tasks {
		t.Submissions = subMap[t.ID]
	}

	jobIDs := make([]int, 0, len(tasks))
	for _, t := range tasks {
		if t.JobPostingID != nil {
			jobIDs = append(jobIDs, int(*t.JobPostingID))
		}
	}
	if len(jobIDs) > 0 {
		type jobRow struct {
			ID    int
			Title string
		}
		var jobs []jobRow
		s.db.Table("job_postings").Select("id, title").Where("id IN ?", jobIDs).Scan(&jobs)
		jobMap := make(map[int]string)
		for _, j := range jobs {
			jobMap[j.ID] = j.Title
		}
		for _, t := range tasks {
			if t.JobPostingID != nil {
				t.JobTitle = jobMap[int(*t.JobPostingID)]
			}
		}
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

func (s *Service) ListAllSubmissions() ([]*SubmissionView, error) {
	type row struct {
		ID              int
		TaskID          uint
		ApplicantID     uint
		Content         string
		FilePath        string
		Grade           *int
		Feedback        string
		Status          string
		CreatedAt       time.Time
		TaskTitle       string
		TaskDescription string
		ApplicantName   string
		JobPostingID    *uint
		JobTitle        string
	}
	var rows []row
	if err := s.db.
		Table("task_submissions").
		Select("task_submissions.id, task_submissions.task_id, task_submissions.applicant_id, task_submissions.content, task_submissions.file_path, task_submissions.grade, task_submissions.feedback, task_submissions.status, task_submissions.created_at, tasks.title as task_title, tasks.description as task_description, tasks.job_posting_id, users.full_name as applicant_name, COALESCE(job_postings.title, '') as job_title").
		Joins("JOIN tasks ON tasks.id = task_submissions.task_id AND tasks.deleted_at IS NULL").
		Joins("JOIN users ON users.id = task_submissions.applicant_id AND users.deleted_at IS NULL").
		Joins("LEFT JOIN job_postings ON job_postings.id = tasks.job_posting_id AND job_postings.deleted_at IS NULL").
		Where("task_submissions.deleted_at IS NULL").
		Order("task_submissions.created_at DESC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]*SubmissionView, len(rows))
	for i, r := range rows {
		sub := TaskSubmission{
			TaskID:      r.TaskID,
			ApplicantID: r.ApplicantID,
			Content:     r.Content,
			FilePath:    r.FilePath,
			Grade:       r.Grade,
			Feedback:    r.Feedback,
			Status:      r.Status,
		}
		sub.ID = r.ID
		sub.CreatedAt = r.CreatedAt
		result[i] = &SubmissionView{
			TaskSubmission:  sub,
			TaskTitle:       r.TaskTitle,
			TaskDescription: r.TaskDescription,
			ApplicantName:   r.ApplicantName,
			HasFile:         r.FilePath != "",
			JobPostingID:    r.JobPostingID,
			JobTitle:        r.JobTitle,
		}
	}
	return result, nil
}
