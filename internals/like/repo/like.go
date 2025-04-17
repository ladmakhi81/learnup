package repo

import (
	"errors"
	"github.com/ladmakhi81/learnup/db"
	"github.com/ladmakhi81/learnup/db/entities"
	"gorm.io/gorm"
)

type FetchLikeOptions struct {
	Page     *int
	PageSize *int
	CourseID *uint
	Preloads []string
}

type FetchCountLikeOptions struct {
	CourseID *uint
}

type FindOneLikeOptions struct {
	ID       *uint
	CourseID *uint
	UserID   *uint
}

type LikeRepo interface {
	Create(like *entities.Like) error
	Update(like *entities.Like) error
	Fetch(options FetchLikeOptions) ([]*entities.Like, error)
	FetchCount(options FetchCountLikeOptions) (int, error)
	FindOne(options FindOneLikeOptions) (*entities.Like, error)
}

type LikeRepoImpl struct {
	dbClient *db.Database
}

func NewLikeRepoImpl(dbClient *db.Database) *LikeRepoImpl {
	return &LikeRepoImpl{
		dbClient: dbClient,
	}
}

func (repo LikeRepoImpl) Create(like *entities.Like) error {
	tx := repo.dbClient.Core.Create(like)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (repo LikeRepoImpl) Update(like *entities.Like) error {
	tx := repo.dbClient.Core.Updates(like)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (repo LikeRepoImpl) Fetch(options FetchLikeOptions) ([]*entities.Like, error) {
	var likes []*entities.Like
	query := repo.dbClient.Core.Where("type != ?", entities.LikeType_None)
	if options.CourseID != nil {
		query = query.Where("course_id = ?", options.CourseID)
	}
	if options.PageSize != nil && options.Page != nil {
		query = query.Offset((*options.Page) * (*options.PageSize)).Limit(*options.PageSize)
	}
	if options.Preloads != nil && len(options.Preloads) > 0 {
		for _, preload := range options.Preloads {
			query = query.Preload(preload)
		}
	}
	tx := query.Order("created_at desc").Find(&likes)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return likes, nil
}

func (repo LikeRepoImpl) FetchCount(options FetchCountLikeOptions) (int, error) {
	var count int64
	query := repo.dbClient.Core.Model(&entities.Like{}).Where("type != ?", entities.LikeType_None)
	if options.CourseID != nil {
		query = query.Where("course_id = ?", options.CourseID)
	}
	tx := query.Count(&count)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return int(count), nil
}

func (repo LikeRepoImpl) FindOne(options FindOneLikeOptions) (*entities.Like, error) {
	var like *entities.Like
	query := repo.dbClient.Core
	if options.ID != nil {
		query = query.Where("id = ?", options.ID)
	}
	if options.CourseID != nil {
		query = query.Where("course_id = ?", options.CourseID)
	}
	if options.UserID != nil {
		query = query.Where("user_id = ?", options.UserID)
	}
	tx := query.First(&like)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return like, nil
}
