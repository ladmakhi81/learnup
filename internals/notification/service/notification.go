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
	const operationName = "NotificationServiceImpl.SeenById"
	notification, err := svc.unitOfWork.NotificationRepo().GetByID(id, nil)
	if err != nil {
		return types.NewServerError(
			"Error in fetching single notification with id",
			operationName,
			err,
		)
	}
	if notification == nil {
		return types.NewNotFoundError(svc.translationSvc.Translate("notification.errors.not_found"))
	}
	notification.IsSeen = true
	now := time.Now()
	notification.SeenAt = &now
	if err := svc.unitOfWork.NotificationRepo().Update(notification); err != nil {
		return types.NewServerError(
			"Error in updating seen status in notification",
			operationName,
			err,
		)
	}
	return nil
}

func (svc NotificationServiceImpl) FetchPageable(page, pageSize int) ([]*notificationEntity.Notification, int, error) {
	const operationName = "NotificationServiceImpl.FetchPageable"
	notifications, count, err := svc.unitOfWork.NotificationRepo().GetPaginated(
		repositories.GetPaginatedOptions{
			Offset:    &page,
			Limit:     &pageSize,
			Relations: []string{"User"},
		},
	)
	if err != nil {
		return nil, 0, types.NewServerError(
			"Fetch All Notifications Throw Error",
			operationName,
			err,
		)
	}
	return notifications, count, nil
}
