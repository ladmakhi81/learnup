package repositories

import (
	"github.com/ladmakhi81/learnup/shared/db/entities"
	"gorm.io/gorm"
)

type CourseParticipantRepo interface {
	Create(courseParticipant *entities.CourseParticipant) error
}

type courseParticipantRepo struct {
	db *gorm.DB
}

func NewCourseParticipantRepo(db *gorm.DB) CourseParticipantRepo {
	return &courseParticipantRepo{
		db: db,
	}
}

func (repo courseParticipantRepo) Create(courseParticipant *entities.CourseParticipant) error {
	return repo.db.Create(courseParticipant).Error
}
