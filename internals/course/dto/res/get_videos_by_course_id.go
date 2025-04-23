package dtores

import (
	entities2 "github.com/ladmakhi81/learnup/internals/db/entities"
	"time"
)

type verifiedByUser struct {
	ID       uint   `json:"id"`
	FullName string `json:"fullName"`
	Phone    string `json:"phone"`
}

type GetVideoByCourseItemDto struct {
	ID           uint                       `json:"id"`
	CreatedAt    time.Time                  `json:"createdAt"`
	UpdatedAt    time.Time                  `json:"updatedAt"`
	Title        string                     `json:"title"`
	Description  string                     `json:"description"`
	AccessLevel  entities2.VideoAccessLevel `json:"accessLevel"`
	Duration     *string                    `json:"duration"`
	URL          string                     `json:"url"`
	IsPublished  bool                       `json:"isPublished"`
	IsVerified   bool                       `json:"isVerified"`
	VerifiedDate *time.Time                 `json:"verifiedDate"`
	VerifiedBy   *verifiedByUser            `json:"verifiedBy"`
	Status       entities2.VideoStatus      `json:"status"`
}

func MapGetVideoByCourseItemsDto(videos []*entities2.Video) []*GetVideoByCourseItemDto {
	result := make([]*GetVideoByCourseItemDto, len(videos))
	for videoIndex, video := range videos {
		result[videoIndex] = &GetVideoByCourseItemDto{
			URL:          video.URL,
			AccessLevel:  video.AccessLevel,
			Description:  video.Description,
			Title:        video.Title,
			Duration:     video.Duration,
			IsPublished:  video.IsPublished,
			ID:           video.ID,
			Status:       video.Status,
			IsVerified:   video.IsVerified,
			UpdatedAt:    video.UpdatedAt,
			CreatedAt:    video.CreatedAt,
			VerifiedDate: video.VerifiedDate,
		}
		if video.VerifiedBy != nil {
			result[videoIndex].VerifiedBy = &verifiedByUser{
				ID:       video.VerifiedBy.ID,
				Phone:    video.VerifiedBy.Phone,
				FullName: video.VerifiedBy.FullName(),
			}
		}
	}
	return result
}
