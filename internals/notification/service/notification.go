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
	unitOfWork     db.UnitOfWork
	translationSvc contracts.Translator
}

func NewNotificationServiceImpl(
	unitOfWork db.UnitOfWork,
	translationSvc contracts.Translator,
) *NotificationServiceImpl {
	return &NotificationServiceImpl{
		unitOfWork:     unitOfWork,
		translationSvc: translationSvc,
	}
}

func (svc NotificationServiceImpl) SeenById(id uint) error {
	tx, txErr := svc.unitOfWork.Begin()
	if txErr != nil {
		return txErr
	}
	notification, notificationErr := tx.NotificationRepo().GetByID(id, nil)
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
	if err := tx.NotificationRepo().Update(notification); err != nil {
		return types.NewServerError(
			"Error in updating seen status in notification",
			"NotificationServiceImpl.SeenById",
			err,
		)
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (svc NotificationServiceImpl) FetchPageable(page, pageSize int) ([]*notificationEntity.Notification, int, error) {
	tx, txErr := svc.unitOfWork.Begin()
	if txErr != nil {
		return nil, 0, txErr
	}
	notifications, count, notificationsErr := tx.NotificationRepo().GetPaginated(
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
	if err := tx.Commit(); err != nil {
		return nil, 0, err
	}
	return notifications, count, nil
}
