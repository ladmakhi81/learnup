package repo

import (
	"errors"
	"github.com/ladmakhi81/learnup/db"
	videoEntity "github.com/ladmakhi81/learnup/internals/video/entity"
	"gorm.io/gorm"
)

type VideoRepo interface {
	Create(video *videoEntity.Video) error
	FetchByTitle(title string) (*videoEntity.Video, error)
	FetchById(id uint) (*videoEntity.Video, error)
	FetchByCourseId(courseID uint) ([]*videoEntity.Video, error)
	Update(video *videoEntity.Video) error
}

type VideoRepoImpl struct {
	dbClient *db.Database
}

func NewVideoRepoImpl(dbClient *db.Database) *VideoRepoImpl {
	return &VideoRepoImpl{
		dbClient: dbClient,
	}
}

func (repo VideoRepoImpl) Create(video *videoEntity.Video) error {
	tx := repo.dbClient.Core.Create(video)
	return tx.Error
}

func (repo VideoRepoImpl) FetchByTitle(title string) (*videoEntity.Video, error) {
	video := &videoEntity.Video{}
	tx := repo.dbClient.Core.Where("title = ?", title).First(video)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return video, nil
}

func (repo VideoRepoImpl) FetchByCourseId(courseID uint) ([]*videoEntity.Video, error) {
	videos := make([]*videoEntity.Video, 0)
	tx := repo.dbClient.Core.
		Preload("VerifiedBy").
		Where("course_id = ?", courseID).
		Order("created_at desc").
		Find(&videos)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return videos, nil
}

func (repo VideoRepoImpl) Update(video *videoEntity.Video) error {
	tx := repo.dbClient.Core.Updates(video)
	return tx.Error
}

func (repo VideoRepoImpl) FetchById(id uint) (*videoEntity.Video, error) {
	video := &videoEntity.Video{}
	tx := repo.dbClient.Core.
		Where("id = ?", id).
		First(video)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return video, nil
}
