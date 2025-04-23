package error

import (
	"github.com/ladmakhi81/learnup/shared/types"
)

var (
	Comment_SenderNotFound = types.NewNotFoundError("comment.errors.sender_not_found")
	Comment_ParentNotFound = types.NewNotFoundError("comment.errors.parent_comment_not_found")
	Comment_NotFound       = types.NewNotFoundError("comment.errors.not_found")
)
