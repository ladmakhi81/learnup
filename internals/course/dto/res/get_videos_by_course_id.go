package dtores

import (
	"github.com/ladmakhi81/learnup/db/entities"
	"time"
)

type verifiedByUser struct {
	ID       uint   `json:"id"`
	FullName string `json:"fullName"`
	Phone    string `json:"phone"`
}

type videosItem struct {
	ID           uint                      `json:"id"`
	CreatedAt    time.Time                 `json:"createdAt"`
	UpdatedAt    time.Time                 `json:"updatedAt"`
	Title        string                    `json:"title"`
	Description  string                    `json:"description"`
	AccessLevel  entities.VideoAccessLevel `json:"accessLevel"`
	Duration     *string                   `json:"duration"`
	URL          string                    `json:"url"`
	IsPublished  bool                      `json:"isPublished"`
	IsVerified   bool                      `json:"isVerified"`
	VerifiedDate *time.Time                `json:"verifiedDate"`
	VerifiedBy   *verifiedByUser           `json:"verifiedBy"`
	Status       entities.VideoStatus      `json:"status"`
}

type GetVideosByCourseIDRes struct {
	Videos   []*videosItem `json:"videos"`
	CourseID uint          `json:"courseId"`
}

func mapper(videos []*entities.Video) []*videosItem {
	result := make([]*videosItem, len(videos))
	for videoIndex, video := range videos {
		result[videoIndex] = &videosItem{
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

func NewGetVideosByCourseIDRes(videos []*entities.Video, courseID uint) GetVideosByCourseIDRes {
	return GetVideosByCourseIDRes{
		Videos:   mapper(videos),
		CourseID: courseID,
	}
}
