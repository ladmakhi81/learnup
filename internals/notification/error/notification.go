package error

import (
	"github.com/ladmakhi81/learnup/shared/types"
)

var (
	Notification_NotFound = types.NewNotFoundError("notification.errors.not_found")
)
