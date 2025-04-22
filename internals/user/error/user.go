package error

import "github.com/ladmakhi81/learnup/types"

var (
	User_NotFound        = types.NewNotFoundError("user.errors.not_found")
	User_AdminNotFound   = types.NewNotFoundError("user.errors.admin_not_found")
	User_PhoneDuplicated = types.NewConflictError("user.errors.phone_duplicate")
	User_TeacherNotFound = types.NewNotFoundError("user.errors.teacher_not_found")
)
