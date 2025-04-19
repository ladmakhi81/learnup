package dtoreq

type AnswerQuestionReq struct {
	QuestionID uint   `json:"-"`
	SenderID   uint   `json:"-"`
	Content    string `json:"content" validate:"required"`
}
