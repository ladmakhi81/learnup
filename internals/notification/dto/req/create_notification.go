package dtoreq

import (
	notificationEntity "github.com/ladmakhi81/learnup/db/entities"
)

type CreateNotificationReq struct {
	UserID    uint                                `json:"user_id"`
	EventType notificationEntity.NotificationType `json:"event_type"`
	Metadata  any                                 `json:"metadata"`
}

func NewCreateNotificationReq(userID uint, eventType notificationEntity.NotificationType, metadata any) CreateNotificationReq {
	return CreateNotificationReq{
		UserID:    userID,
		EventType: eventType,
		Metadata:  metadata,
	}
}
