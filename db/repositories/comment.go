package repositories

import (
	"github.com/ladmakhi81/learnup/db/entities"
	"gorm.io/gorm"
)

type CommentRepo interface {
	Repository[entities.Comment]
}

type CommentRepoImpl struct {
	RepositoryImpl[entities.Comment]
}

func NewCommentRepo(db *gorm.DB) *CommentRepoImpl {
	return &CommentRepoImpl{
		RepositoryImpl[entities.Comment]{
			db: db,
		},
	}
}
