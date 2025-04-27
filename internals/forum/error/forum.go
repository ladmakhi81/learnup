package error

import "github.com/ladmakhi81/learnup/shared/types"

var (
	Forum_NotFound = types.NewNotFoundError("forum.errors.not_found")
)
