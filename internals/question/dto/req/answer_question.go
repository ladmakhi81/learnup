package dtoreq

type AnswerQuestionReq struct {
	QuestionID uint   `json:"-"`
	Content    string `json:"content" validate:"required"`
}
