package dtoreq

import "github.com/ladmakhi81/learnup/db/entities"

type CreateQuestionReq struct {
	UserID   uint                      `json:"-"`
	CourseID uint                      `json:"-"`
	Content  string                    `json:"content" validate:"required,min=4"`
	Priority entities.QuestionPriority `json:"priority" validate:"required,oneof=high low"`
	VideoID  *uint                     `json:"videoId" validate:"omitempty,numeric"`
}
