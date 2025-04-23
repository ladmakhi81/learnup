package dtoreq

import (
	"github.com/ladmakhi81/learnup/internals/db/entities"
)

type CreateQuestionReq struct {
	CourseID uint                      `json:"-"`
	Content  string                    `json:"content" validate:"required,min=4"`
	Priority entities.QuestionPriority `json:"priority" validate:"required,oneof=high low"`
	VideoID  *uint                     `json:"videoId" validate:"omitempty,numeric"`
}
