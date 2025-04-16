package repo

import (
	"errors"
	"github.com/ladmakhi81/learnup/db"
	"github.com/ladmakhi81/learnup/db/entities"
	"gorm.io/gorm"
)

type FetchPageOption struct {
	PageSize   *int
	Page       *int
	TeacherId  *uint
	OrderField *string
	Preloads   []string
}

type FetchCountOption struct {
	TeacherId *uint
}

type CourseRepo interface {
	Create(course *entities.Course) error
	FetchByName(name string) (*entities.Course, error)
	FetchPage(opt FetchPageOption) ([]*entities.Course, error)
	FetchCount(opt FetchCountOption) (int, error)
	FetchById(id uint) (*entities.Course, error)
	FetchDetailById(id uint) (*entities.Course, error)
	FetchByVideoId(id uint) (*entities.Course, error)
	Update(course *entities.Course) error
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

func (repo CourseRepoImpl) FetchPage(opt FetchPageOption) ([]*entities.Course, error) {
	var courses []*entities.Course
	query := repo.dbClient.Core
	if opt.Preloads != nil && len(opt.Preloads) > 0 {
		for _, preload := range opt.Preloads {
			query = query.Preload(preload)
		}
	}
	if opt.TeacherId != nil {
		query = query.Where("teacher_id = ?", opt.TeacherId)
	}
	if opt.Page != nil && opt.PageSize != nil {
		query = query.Offset((*opt.Page) * (*opt.PageSize)).Limit(*opt.PageSize)
	}
	orderField := "created_at desc"
	if opt.OrderField != nil {
		orderField = *opt.OrderField
	}

	tx := query.Order(orderField).Find(&courses)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return courses, nil
}

func (repo CourseRepoImpl) FetchCount(opt FetchCountOption) (int, error) {
	count := int64(0)
	query := repo.dbClient.Core.Model(&entities.Course{})
	if opt.TeacherId != nil {
		query = query.Where("teacher_id = ?", opt.TeacherId)
	}
	tx := query.Count(&count)
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

func (repo CourseRepoImpl) Update(course *entities.Course) error {
	tx := repo.dbClient.Core.Updates(course)
	return tx.Error
}

func (repo CourseRepoImpl) FetchByTeacherId(teacherId uint, page, pageSize int) ([]*entities.Course, error) {
	var courses []*entities.Course
	tx := repo.dbClient.Core.
		Where("teacher_id = ?", teacherId).
		Order("created_at desc").
		Offset(page * pageSize).
		Limit(pageSize).
		Find(&courses)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return courses, nil
}

func (repo CourseRepoImpl) FetchCountByTeacherId(teacherId int64) (int, error) {
	var count int64
	tx := repo.dbClient.Core.Model(&entities.Course{}).Where("teacher_id = ?", teacherId).Count(&count)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return int(count), nil
}
