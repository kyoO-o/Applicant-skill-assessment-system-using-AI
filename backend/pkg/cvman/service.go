package cvman

import (
	"errors"
	"log"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("cvprofile.not_found")

type Service struct {
	db       *gorm.DB
	infoLog  *log.Logger
	errorLog *log.Logger
}

func NewService(db *gorm.DB, infoLog, errorLog *log.Logger) *Service {
	return &Service{db: db, infoLog: infoLog, errorLog: errorLog}
}

func (s *Service) GetByUserID(userID int) (*CVProfile, error) {
	var cv CVProfile
	if err := s.db.Where("user_id = ?", userID).First(&cv).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &cv, nil
}

func (s *Service) Save(cv *CVProfile) (*CVProfile, error) {
	if err := s.db.Save(cv).Error; err != nil {
		return nil, err
	}
	return cv, nil
}
