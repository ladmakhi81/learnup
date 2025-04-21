package repositories

import (
	"github.com/ladmakhi81/learnup/db/entities"
	"gorm.io/gorm"
)

type NotificationRepo interface {
	Repository[entities.Notification]
}

type NotificationRepoImpl struct {
	RepositoryImpl[entities.Notification]
}

func NewNotificationRepo(db *gorm.DB) *NotificationRepoImpl {
	return &NotificationRepoImpl{
		RepositoryImpl[entities.Notification]{
			db: db,
		},
	}
}
