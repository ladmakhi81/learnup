package repositories

import (
	"errors"
	"github.com/ladmakhi81/learnup/db/entities"
	"gorm.io/gorm"
)

type GetByTeacherIDOption struct {
	TeacherID uint
	Page      int
	PageSize  int
}

type CourseRepo interface {
	Repository[entities.Course]
	GetByVideoID(videoID uint) (*entities.Course, error)
	GetByTeacherID(options GetByTeacherIDOption) ([]*entities.Course, int, error)
}
type CourseRepoImpl struct {
	RepositoryImpl[entities.Course]
}

func NewCourseRepo(db *gorm.DB) *CourseRepoImpl {
	return &CourseRepoImpl{
		RepositoryImpl[entities.Course]{
			db: db,
		},
	}
}

func (repo CourseRepoImpl) GetByVideoID(videoID uint) (*entities.Course, error) {
	course := &entities.Course{}
	tx := repo.db.
		Joins("JOIN _videos ON _courses.id = _videos.course_id").
		Where("_videos.id = ?", videoID).
		First(course)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}

	return course, nil
}

func (repo CourseRepoImpl) GetByTeacherID(options GetByTeacherIDOption) ([]*entities.Course, int, error) {
	var courses []*entities.Course
	var count int64

	coursesTx := repo.db.
		Where("teacher_id = ?", options.TeacherID).
		Order("created_at desc").
		Offset(options.Page * options.PageSize).
		Limit(options.PageSize).
		Find(&courses)
	if coursesTx.Error != nil {
		return nil, 0, coursesTx.Error
	}

	countTx := repo.db.Model(&entities.Course{}).
		Where("teacher_id = ?", options.TeacherID).
		Count(&count)
	if countTx.Error != nil {
		return nil, 0, countTx.Error
	}

	return courses, int(count), nil
}
