package repositories

import (
	"github.com/ladmakhi81/learnup/shared/db/entities"
	"gorm.io/gorm"
)

type LikeRepo interface {
	Repository[entities.Like]
}

type LikeRepoImpl struct {
	RepositoryImpl[entities.Like]
}

func NewLikeRepo(db *gorm.DB) *LikeRepoImpl {
	return &LikeRepoImpl{
		RepositoryImpl[entities.Like]{
			db: db,
		},
	}
}
