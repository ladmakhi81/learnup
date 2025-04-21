package repositories

import (
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"gorm.io/gorm"
)

type QuestionRepo interface {
	Repository[entities.Question]
}

type QuestionRepoImpl struct {
	RepositoryImpl[entities.Question]
}

func NewQuestionRepo(db *gorm.DB) *QuestionRepoImpl {
	return &QuestionRepoImpl{
		RepositoryImpl[entities.Question]{
			db: db,
		},
	}
}
