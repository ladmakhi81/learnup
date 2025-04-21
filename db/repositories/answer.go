package repositories

import (
	"github.com/ladmakhi81/learnup/db/entities"
	"gorm.io/gorm"
)

type AnswerRepo interface {
	Repository[entities.QuestionAnswer]
}

type AnswerRepoImpl struct {
	RepositoryImpl[entities.QuestionAnswer]
}

func NewAnswerRepo(db *gorm.DB) *AnswerRepoImpl {
	return &AnswerRepoImpl{
		RepositoryImpl[entities.QuestionAnswer]{
			db: db,
		},
	}
}
