package service

import (
	notificationEntity "github.com/ladmakhi81/learnup/db/entities"
	dtoreq "github.com/ladmakhi81/learnup/internals/notification/dto/req"
	"github.com/ladmakhi81/learnup/internals/notification/repo"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
	"time"
)

type NotificationService interface {
	Create(dto dtoreq.CreateNotificationReq) (*notificationEntity.Notification, error)
	SeenById(id uint) error
	FindById(id uint) (*notificationEntity.Notification, error)
	FetchPageable(page, pageSize int) ([]*notificationEntity.Notification, error)
	FetchCount() (int, error)
}

type NotificationServiceImpl struct {
	notificationRepo repo.NotificationRepo
	userSvc          userService.UserSvc
	translationSvc   contracts.Translator
}

func NewNotificationServiceImpl(
	notificationRepo repo.NotificationRepo,
	userSvc userService.UserSvc,
	translationSvc contracts.Translator,
) *NotificationServiceImpl {
	return &NotificationServiceImpl{
		userSvc:          userSvc,
		notificationRepo: notificationRepo,
		translationSvc:   translationSvc,
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

func (svc NotificationServiceImpl) SeenById(id uint) error {
	notification, notificationErr := svc.FindById(id)
	if notificationErr != nil {
		return notificationErr
	}
	if notification == nil {
		return types.NewNotFoundError(svc.translationSvc.Translate("notification.errors.not_found"))
	}
	notification.IsSeen = true
	now := time.Now()
	notification.SeenAt = &now
	if err := svc.notificationRepo.Update(notification); err != nil {
		return types.NewServerError(
			"Error in updating seen status in notification",
			"NotificationServiceImpl.SeenById",
			err,
		)
	}
	return nil
}

func (svc NotificationServiceImpl) FindById(id uint) (*notificationEntity.Notification, error) {
	notification, notificationErr := svc.notificationRepo.FetchById(id)
	if notificationErr != nil {
		return nil, types.NewServerError(
			"Error in finding notification by id",
			"NotificationServiceImpl.FetchById",
			notificationErr,
		)
	}
	return notification, nil
}

func (svc NotificationServiceImpl) FetchPageable(page, pageSize int) ([]*notificationEntity.Notification, error) {
	notifications, notificationsErr := svc.notificationRepo.FetchPageable(page, pageSize)
	if notificationsErr != nil {
		return nil, types.NewServerError(
			"Fetch All Notifications Throw Error",
			"NotificationServiceImpl.FetchPageable",
			notificationsErr,
		)
	}
	return notifications, nil
}

func (svc NotificationServiceImpl) FetchCount() (int, error) {
	count, countErr := svc.notificationRepo.FetchCount()
	if countErr != nil {
		return 0, types.NewServerError(
			"Fetch Notification Count Throw Error",
			"NotificationServiceImpl.FetchCount",
			countErr,
		)
	}
	return count, nil
}
