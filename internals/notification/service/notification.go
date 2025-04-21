package service

import (
	"github.com/ladmakhi81/learnup/internals/db"
	notificationEntity "github.com/ladmakhi81/learnup/internals/db/entities"
	"github.com/ladmakhi81/learnup/internals/db/repositories"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
	"time"
)

type NotificationService interface {
	SeenById(id uint) error
	FetchPageable(page, pageSize int) ([]*notificationEntity.Notification, int, error)
}

type NotificationServiceImpl struct {
	repo           *db.Repositories
	translationSvc contracts.Translator
}

func NewNotificationServiceImpl(
	repo *db.Repositories,
	translationSvc contracts.Translator,
) *NotificationServiceImpl {
	return &NotificationServiceImpl{
		repo:           repo,
		translationSvc: translationSvc,
	}
}

func (svc NotificationServiceImpl) SeenById(id uint) error {
	notification, notificationErr := svc.repo.NotificationRepo.GetByID(id)
	if notificationErr != nil {
		return types.NewServerError(
			"Error in fetching single notification with id",
			"NotificationServiceImpl.SeenById",
			notificationErr,
		)
	}
	if notification == nil {
		return types.NewNotFoundError(svc.translationSvc.Translate("notification.errors.not_found"))
	}
	notification.IsSeen = true
	now := time.Now()
	notification.SeenAt = &now
	if err := svc.repo.NotificationRepo.Update(notification); err != nil {
		return types.NewServerError(
			"Error in updating seen status in notification",
			"NotificationServiceImpl.SeenById",
			err,
		)
	}
	return nil
}

func (svc NotificationServiceImpl) FetchPageable(page, pageSize int) ([]*notificationEntity.Notification, int, error) {
	notifications, count, notificationsErr := svc.repo.NotificationRepo.GetPaginated(
		repositories.GetPaginatedOptions{
			Offset:    &page,
			Limit:     &pageSize,
			Relations: []string{"User"},
		},
	)
	if notificationsErr != nil {
		return nil, 0, types.NewServerError(
			"Fetch All Notifications Throw Error",
			"NotificationServiceImpl.FetchPageable",
			notificationsErr,
		)
	}
	return notifications, count, nil
}
