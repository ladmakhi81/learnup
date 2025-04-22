package error

import "github.com/ladmakhi81/learnup/types"

var (
	Question_NotFound     = types.NewNotFoundError("question.errors.not_found")
	Question_ClosedStatus = types.NewBadRequestError("question.errors.closed")
)
