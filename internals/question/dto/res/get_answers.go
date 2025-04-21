package dtores

import (
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"time"
)

type senderItem struct {
	ID       uint   `json:"id"`
	FullName string `json:"fullName"`
	Phone    string `json:"phone"`
}

type GetAnswersRes struct {
	ID        uint        `json:"id"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	Content   string      `json:"content"`
	Sender    *senderItem `json:"sender"`
}

func MapGetAnswersRes(answers []*entities.QuestionAnswer) []*GetAnswersRes {
	res := make([]*GetAnswersRes, len(answers))
	for index, answer := range answers {
		res[index] = &GetAnswersRes{
			ID:        answer.ID,
			CreatedAt: answer.CreatedAt,
			UpdatedAt: answer.UpdatedAt,
			Content:   answer.Content,
			Sender: &senderItem{
				ID:       answer.SenderID,
				FullName: answer.Sender.FullName(),
				Phone:    answer.Sender.Phone,
			},
		}
	}
	return res
}
