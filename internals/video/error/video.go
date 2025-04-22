package error

import "github.com/ladmakhi81/learnup/types"

var (
	Video_NotFound        = types.NewNotFoundError("video.errors.not_found")
	Video_TitleDuplicated = types.NewConflictError("video.errors.title_duplicated")
)
