package repo

import (
	"errors"
	"github.com/ladmakhi81/learnup/db"
	"github.com/ladmakhi81/learnup/db/entities"
	"gorm.io/gorm"
)

type FetchAnswerOption struct {
	Page       *int
	PageSize   *int
	Preloads   []string
	QuestionID *uint
}

type QuestionAnswerRepo interface {
	Create(answer *entities.QuestionAnswer) error
	Delete(id uint) error
	FetchById(id uint) (*entities.QuestionAnswer, error)
	Fetch(options FetchAnswerOption) ([]*entities.QuestionAnswer, error)
}

type QuestionAnswerRepoImpl struct {
	dbClient *db.Database
}

func NewQuestionAnswerRepoImpl(dbClient *db.Database) *QuestionAnswerRepoImpl {
	return &QuestionAnswerRepoImpl{
		dbClient: dbClient,
	}
}

func (repo QuestionAnswerRepoImpl) Create(answer *entities.QuestionAnswer) error {
	tx := repo.dbClient.Core.Create(&answer)
	return tx.Error
}

func (repo QuestionAnswerRepoImpl) Delete(id uint) error {
	tx := repo.dbClient.Core.Where("id = ?", id).Delete(&entities.QuestionAnswer{})
	return tx.Error
}

func (repo QuestionAnswerRepoImpl) FetchById(id uint) (*entities.QuestionAnswer, error) {
	var answer *entities.QuestionAnswer
	tx := repo.dbClient.Core.Where("id = ?", id).First(&answer)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return answer, nil
}

func (repo QuestionAnswerRepoImpl) Fetch(options FetchAnswerOption) ([]*entities.QuestionAnswer, error) {
	var answers []*entities.QuestionAnswer
	query := repo.dbClient.Core
	if options.QuestionID != nil {
		query = query.Where("question_id = ?", options.QuestionID)
	}
	if options.PageSize != nil && options.Page != nil {
		query = query.Offset((*options.Page) * (*options.PageSize)).Limit(*options.PageSize)
	}
	if options.Preloads != nil {
		for _, preload := range options.Preloads {
			query = query.Preload(preload)
		}
	}
	tx := query.Order("created_at desc").Find(&answers)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return answers, nil
}
