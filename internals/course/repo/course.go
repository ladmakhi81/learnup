package repo

import (
	"errors"
	"github.com/ladmakhi81/learnup/db"
	"github.com/ladmakhi81/learnup/internals/course/entity"
	"gorm.io/gorm"
)

type CourseRepo interface {
	Create(course *entity.Course) error
	FetchByName(name string) (*entity.Course, error)
	FetchPage(page, pageSize int) ([]*entity.Course, error)
	FetchCount() (int, error)
	FetchById(id uint) (*entity.Course, error)
	FetchDetailById(id uint) (*entity.Course, error)
	FetchByVideoId(id uint) (*entity.Course, error)
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

func (repo CourseRepoImpl) FetchByName(name string) (*entity.Course, error) {
	course := &entity.Course{}
	tx := repo.dbClient.Core.Where("name = ?", name).First(course)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}
	return course, nil
}

func (repo CourseRepoImpl) FetchPage(page, pageSize int) ([]*entity.Course, error) {
	var courses []*entity.Course
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
	tx := repo.dbClient.Core.Model(&entity.Course{}).Count(&count)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return int(count), nil
}

func (repo CourseRepoImpl) FetchById(id uint) (*entity.Course, error) {
	course := &entity.Course{}
	tx := repo.dbClient.Core.Where("id = ?", id).First(course)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return course, nil
}

func (repo CourseRepoImpl) FetchDetailById(id uint) (*entity.Course, error) {
	course := &entity.Course{}
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

func (repo CourseRepoImpl) FetchByVideoId(id uint) (*entity.Course, error) {
	course := &entity.Course{}
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
