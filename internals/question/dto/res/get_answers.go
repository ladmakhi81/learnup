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

type GetAnswerItemDto struct {
	ID        uint        `json:"id"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	Content   string      `json:"content"`
	Sender    *senderItem `json:"sender"`
}

func MapGetAnswerItemsDto(answers []*entities.QuestionAnswer) []*GetAnswerItemDto {
	res := make([]*GetAnswerItemDto, len(answers))
	for index, answer := range answers {
		res[index] = &GetAnswerItemDto{
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
