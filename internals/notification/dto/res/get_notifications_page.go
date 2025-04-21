package dtores

import (
	entities2 "github.com/ladmakhi81/learnup/internals/db/entities"
	"time"
)

type userItem struct {
	ID       uint   `json:"id"`
	FullName string `json:"fullName"`
}

type NotificationPageItem struct {
	Type     entities2.NotificationType `json:"type"`
	IsSeen   bool                       `json:"isSeen"`
	SeenAt   *time.Time                 `json:"seenAt"`
	User     *userItem                  `json:"user"`
	Metadata any                        `json:"metadata"`
}

func NewNotificationPageItems(notifications []*entities2.Notification) []*NotificationPageItem {
	result := make([]*NotificationPageItem, len(notifications))
	for index, notification := range notifications {
		result[index] = &NotificationPageItem{
			Type:     notification.Type,
			IsSeen:   notification.IsSeen,
			Metadata: notification.Metadata,
			User: &userItem{
				ID:       notification.User.ID,
				FullName: notification.User.FullName(),
			},
			SeenAt: notification.SeenAt,
		}
	}
	return result
}
