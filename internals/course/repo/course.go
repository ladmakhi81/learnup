package repo

import (
	"errors"
	"github.com/ladmakhi81/learnup/db"
	"github.com/ladmakhi81/learnup/internals/course/entity"
	"gorm.io/gorm"
)

type CourseRepo interface {
	Create(course *entity.Course) error
	FindByName(name string) (*entity.Course, error)
}

type CourseRepoImpl struct {
	dbClient *db.Database
}

func NewCourseRepoImpl(dbClient *db.Database) *CourseRepoImpl {
	return &CourseRepoImpl{
		dbClient: dbClient,
	}
}

func (repo CourseRepoImpl) Create(course *entity.Course) error {
	tx := repo.dbClient.Core.Create(course)
	return tx.Error
}

func (repo CourseRepoImpl) FindByName(name string) (*entity.Course, error) {
	course := &entity.Course{}
	tx := repo.dbClient.Core.Where("name = ?", name).First(course)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}
	return course, nil
}
