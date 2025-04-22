package error

import "github.com/ladmakhi81/learnup/types"

var (
	Notification_NotFound = types.NewNotFoundError("notification.errors.not_found")
)
