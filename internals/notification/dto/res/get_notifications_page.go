package dtores

import (
	"github.com/ladmakhi81/learnup/shared/db/entities"
	"time"
)

type userItem struct {
	ID       uint   `json:"id"`
	FullName string `json:"fullName"`
}

type NotificationPageItemDto struct {
	Type   entities.NotificationType `json:"type"`
	IsSeen bool                      `json:"isSeen"`
	SeenAt   *time.Time                `json:"seenAt"`
	User     *userItem                 `json:"user"`
	Metadata any                       `json:"metadata"`
}

func NewNotificationPageItemsDto(notifications []*entities.Notification) []*NotificationPageItemDto {
	result := make([]*NotificationPageItemDto, len(notifications))
	for index, notification := range notifications {
		result[index] = &NotificationPageItemDto{
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
