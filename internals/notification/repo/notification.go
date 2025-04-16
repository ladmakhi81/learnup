package repo

import (
	"errors"
	"github.com/ladmakhi81/learnup/db"
	notificationEntity "github.com/ladmakhi81/learnup/db/entities"
	"gorm.io/gorm"
)

type NotificationRepo interface {
	Create(notification *notificationEntity.Notification) error
	Update(notification *notificationEntity.Notification) error
	FetchById(id uint) (*notificationEntity.Notification, error)
	FetchPageable(page, pageSize int) ([]*notificationEntity.Notification, error)
	FetchCount() (int, error)
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

func (repo NotificationRepoImpl) Update(notification *notificationEntity.Notification) error {
	tx := repo.dbClient.Core.Updates(notification)
	return tx.Error
}

func (repo NotificationRepoImpl) FetchById(id uint) (*notificationEntity.Notification, error) {
	notification := &notificationEntity.Notification{}
	tx := repo.dbClient.Core.Where("id = ?", id).First(&notification)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return notification, nil
}

func (repo NotificationRepoImpl) FetchPageable(page, pageSize int) ([]*notificationEntity.Notification, error) {
	notifications := make([]*notificationEntity.Notification, 0)
	tx := repo.dbClient.Core.
		Preload("User").
		Offset(page * pageSize).
		Limit(pageSize).
		Order("created_at desc").
		Find(&notifications)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return notifications, nil
}

func (repo NotificationRepoImpl) FetchCount() (int, error) {
	var count int64
	tx := repo.dbClient.Core.Model(&notificationEntity.Notification{}).Count(&count)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return int(count), nil
}
