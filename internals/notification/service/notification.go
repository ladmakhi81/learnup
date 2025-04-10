package service

import (
	dtoreq "github.com/ladmakhi81/learnup/internals/notification/dto/req"
	notificationEntity "github.com/ladmakhi81/learnup/internals/notification/entity"
	"github.com/ladmakhi81/learnup/internals/notification/repo"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	"github.com/ladmakhi81/learnup/types"
)

type NotificationService interface {
	Create(dto dtoreq.CreateNotificationReq) (*notificationEntity.Notification, error)
}

type NotificationServiceImpl struct {
	notificationRepo repo.NotificationRepo
	userSvc          userService.UserSvc
}

func NewNotificationServiceImpl(
	notificationRepo repo.NotificationRepo,
	userSvc userService.UserSvc,
) *NotificationServiceImpl {
	return &NotificationServiceImpl{
		userSvc:          userSvc,
		notificationRepo: notificationRepo,
	}
}

func (svc NotificationServiceImpl) Create(dto dtoreq.CreateNotificationReq) (*notificationEntity.Notification, error) {
	user, userErr := svc.userSvc.FindById(dto.UserID)
	if userErr != nil {
		return nil, userErr
	}
	notification := &notificationEntity.Notification{
		Type:     dto.EventType,
		Metadata: dto.Metadata,
		IsSeen:   false,
		UserID:   &user.ID,
	}
	if err := svc.notificationRepo.Create(notification); err != nil {
		return nil, types.NewServerError(
			"Error in creating notification",
			"NotificationServiceImpl.Create",
			err,
		)
	}
	return notification, nil
}
