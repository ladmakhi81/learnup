package repo

import (
	"errors"
	"github.com/ladmakhi81/learnup/db"
	"github.com/ladmakhi81/learnup/db/entities"
	"gorm.io/gorm"
)

type FetchCommentOption struct {
	UserID     *uint
	CourseID   *uint
	Page       *int
	PageSize   *int
	OrderField *string
	Preloads   []string
}

type FetchCountCommentOption struct {
	UserID   *uint
	CourseID *uint
}

type CommentRepo interface {
	Create(comment *entities.Comment) error
	FindById(id uint) (*entities.Comment, error)
	Delete(id uint) error
	Fetch(option FetchCommentOption) ([]*entities.Comment, error)
	FetchCount(option FetchCountCommentOption) (int, error)
}

type CommentRepoImpl struct {
	dbClient *db.Database
}

func NewCommentRepoImpl(dbClient *db.Database) *CommentRepoImpl {
	return &CommentRepoImpl{dbClient: dbClient}
}

func (repo CommentRepoImpl) Create(comment *entities.Comment) error {
	tx := repo.dbClient.Core.Create(comment)
	return tx.Error
}

func (repo CommentRepoImpl) FindById(id uint) (*entities.Comment, error) {
	comment := &entities.Comment{}
	tx := repo.dbClient.Core.Where("id = ?", id).First(comment)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return comment, nil
}

func (repo CommentRepoImpl) Delete(id uint) error {
	tx := repo.dbClient.Core.Where("id = ?", id).Delete(&entities.Comment{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (repo CommentRepoImpl) Fetch(option FetchCommentOption) ([]*entities.Comment, error) {
	var comments []*entities.Comment
	query := repo.dbClient.Core
	if option.UserID != nil {
		query = query.Where("user_id = ?", option.UserID)
	}
	if option.CourseID != nil {
		query = query.Where("course_id = ?", option.CourseID)
	}
	if option.Page != nil && option.PageSize != nil {
		query = query.Offset((*option.Page) * (*option.PageSize)).Limit(*option.PageSize)
	}
	orderField := "created_at desc"
	if option.OrderField != nil {
		orderField = *option.OrderField
	}
	if option.Preloads != nil {
		for _, preload := range option.Preloads {
			query = query.Preload(preload)
		}
	}
	tx := query.Order(orderField).Find(&comments)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return comments, nil
}

func (repo CommentRepoImpl) FetchCount(option FetchCountCommentOption) (int, error) {
	var count int64
	query := repo.dbClient.Core.Model(&entities.Comment{})
	if option.UserID != nil {
		query = query.Where("user_id = ?", option.UserID)
	}
	if option.CourseID != nil {
		query = query.Where("course_id = ?", option.CourseID)
	}
	tx := query.Count(&count)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return int(count), nil
}
