package dtoreq

type AnswerQuestionReqDto struct {
	QuestionID uint   `json:"-"`
	Content    string `json:"content" validate:"required"`
}
