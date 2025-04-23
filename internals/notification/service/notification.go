package service

import (
	"github.com/ladmakhi81/learnup/internals/db"
	notificationEntity "github.com/ladmakhi81/learnup/internals/db/entities"
	"github.com/ladmakhi81/learnup/internals/db/repositories"
	notificationError "github.com/ladmakhi81/learnup/internals/notification/error"
	"github.com/ladmakhi81/learnup/types"
	"github.com/ladmakhi81/learnup/utils"
)

type NotificationService interface {
	SeenById(id uint) error
	FetchPageable(page, pageSize int) ([]*notificationEntity.Notification, int, error)
}

type notificationService struct {
	unitOfWork db.UnitOfWork
}

func NewNotificationSvc(unitOfWork db.UnitOfWork) NotificationService {
	return &notificationService{unitOfWork: unitOfWork}
}

func (svc notificationService) SeenById(id uint) error {
	const operationName = "notificationService.SeenById"
	notification, err := svc.unitOfWork.NotificationRepo().GetByID(id, nil)
	if err != nil {
		return types.NewServerError("Error in fetching single notification with id", operationName, err)
	}
	if notification == nil {
		return notificationError.Notification_NotFound
	}
	notification.IsSeen = true
	notification.SeenAt = utils.Now()
	if err := svc.unitOfWork.NotificationRepo().Update(notification); err != nil {
		return types.NewServerError("Error in updating seen status in notification", operationName, err)
	}
	return nil
}

func (svc notificationService) FetchPageable(page, pageSize int) ([]*notificationEntity.Notification, int, error) {
	const operationName = "notificationService.FetchPageable"
	notifications, count, err := svc.unitOfWork.NotificationRepo().GetPaginated(
		repositories.GetPaginatedOptions{
			Offset:    &page,
			Limit:     &pageSize,
			Relations: []string{"User"},
		},
	)
	if err != nil {
		return nil, 0, types.NewServerError("Fetch All Notifications Throw Error", operationName, err)
	}
	return notifications, count, nil
}
