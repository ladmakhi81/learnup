package repositories

import (
	"github.com/ladmakhi81/learnup/db/entities"
	"gorm.io/gorm"
)

type VideoRepo interface {
	Repository[entities.Video]
}

type VideoRepoImpl struct {
	RepositoryImpl[entities.Video]
}

func NewVideoRepo(db *gorm.DB) *VideoRepoImpl {
	return &VideoRepoImpl{
		RepositoryImpl[entities.Video]{
			db: db,
		},
	}
}
