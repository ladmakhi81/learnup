package repositories

import (
	"github.com/ladmakhi81/learnup/shared/db/entities"
	"gorm.io/gorm"
)

type CourseForumRepo interface {
	Repository[entities.CourseForum]
}

type courseForumRepo struct {
	RepositoryImpl[entities.CourseForum]
}

func NewCourseForumRepo(db *gorm.DB) CourseForumRepo {
	return &courseForumRepo{
		RepositoryImpl[entities.CourseForum]{
			db: db,
		},
	}
}
