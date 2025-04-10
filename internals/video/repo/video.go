package repo

import (
	"errors"
	"github.com/ladmakhi81/learnup/db"
	videoEntity "github.com/ladmakhi81/learnup/internals/video/entity"
	"gorm.io/gorm"
)

type VideoRepo interface {
	Create(video *videoEntity.Video) error
	FindByTitle(title string) (*videoEntity.Video, error)
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

func (repo VideoRepoImpl) FindByTitle(title string) (*videoEntity.Video, error) {
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
