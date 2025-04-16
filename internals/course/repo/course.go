package repo

import (
	"errors"
	"github.com/ladmakhi81/learnup/db"
	"github.com/ladmakhi81/learnup/db/entities"
	"gorm.io/gorm"
)

type CourseRepo interface {
	Create(course *entities.Course) error
	FetchByName(name string) (*entities.Course, error)
	FetchPage(page, pageSize int) ([]*entities.Course, error)
	FetchCount() (int, error)
	FetchById(id uint) (*entities.Course, error)
	FetchDetailById(id uint) (*entities.Course, error)
	FetchByVideoId(id uint) (*entities.Course, error)
}

type CourseRepoImpl struct {
	dbClient *db.Database
}

func NewCourseRepoImpl(dbClient *db.Database) *CourseRepoImpl {
	return &CourseRepoImpl{
		dbClient: dbClient,
	}
}

func (repo CourseRepoImpl) Create(course *entities.Course) error {
	tx := repo.dbClient.Core.Create(course)
	return tx.Error
}

func (repo CourseRepoImpl) FetchByName(name string) (*entities.Course, error) {
	course := &entities.Course{}
	tx := repo.dbClient.Core.Where("name = ?", name).First(course)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}
	return course, nil
}

func (repo CourseRepoImpl) FetchPage(page, pageSize int) ([]*entities.Course, error) {
	var courses []*entities.Course
	tx := repo.dbClient.Core.
		Preload("Teacher").
		Preload("Category").
		Preload("VerifiedBy").
		Offset(page * pageSize).
		Limit(pageSize).
		Order("created_at desc").
		Find(&courses)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return courses, nil
}

func (repo CourseRepoImpl) FetchCount() (int, error) {
	count := int64(0)
	tx := repo.dbClient.Core.Model(&entities.Course{}).Count(&count)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return int(count), nil
}

func (repo CourseRepoImpl) FetchById(id uint) (*entities.Course, error) {
	course := &entities.Course{}
	tx := repo.dbClient.Core.Where("id = ?", id).First(course)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return course, nil
}

func (repo CourseRepoImpl) FetchDetailById(id uint) (*entities.Course, error) {
	course := &entities.Course{}
	tx := repo.dbClient.Core.
		Where("id = ?", id).
		Preload("Teacher").
		Preload("Category").
		Preload("VerifiedBy").
		First(course)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return course, nil
}

func (repo CourseRepoImpl) FetchByVideoId(id uint) (*entities.Course, error) {
	course := &entities.Course{}
	tx := repo.dbClient.Core.
		Joins("JOIN _videos ON _courses.id = _videos.course_id").
		Where("_videos.id = ?", id).
		First(course)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}

	return course, nil
}
