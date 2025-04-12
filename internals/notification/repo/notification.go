package repo

import (
	"github.com/ladmakhi81/learnup/db"
	notificationEntity "github.com/ladmakhi81/learnup/internals/notification/entity"
)

type NotificationRepo interface {
	Create(notification *notificationEntity.Notification) error
}

type NotificationRepoImpl struct {
	dbClient *db.Database
}

func NewNotificationRepoImpl(dbClient *db.Database) *NotificationRepoImpl {
	return &NotificationRepoImpl{
		dbClient: dbClient,
	}
}

func (repo NotificationRepoImpl) Create(notification *notificationEntity.Notification) error {
	tx := repo.dbClient.Core.Create(notification)
	return tx.Error
}
