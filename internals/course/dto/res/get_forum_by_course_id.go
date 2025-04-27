package dtores

import (
	"github.com/ladmakhi81/learnup/shared/db/entities"
	"time"
)

type getForumByCourseIDCourse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type getForumByCourseIDTeacher struct {
	ID       uint   `json:"id"`
	FullName string `json:"name"`
	Phone    string `json:"phone"`
}

type getForumByCourseIDMembers struct {
	ID       uint   `json:"id"`
	FullName string `json:"name"`
}

type GetForumByCourseIDDto struct {
	ID              uint                           `json:"id"`
	CreatedAt       time.Time                      `json:"createdAt"`
	UpdatedAt       time.Time                      `json:"updatedAt"`
	Status          entities.CourseForumStatus     `json:"status"`
	StatusChangedAt *time.Time                     `json:"statusChangedAt"`
	IsPublic        bool                           `json:"isPublic"`
	AccessMode      entities.CourseForumAccessMode `json:"accessMode"`
	Course          getForumByCourseIDCourse       `json:"course"`
	Teacher         getForumByCourseIDTeacher      `json:"teacher"`
	Members         []getForumByCourseIDMembers    `json:"members"`
}

func MapGetForumByCourseIDDto(forum *entities.CourseForum) GetForumByCourseIDDto {
	res := GetForumByCourseIDDto{
		ID:         forum.ID,
		CreatedAt:  forum.CreatedAt,
		UpdatedAt:  forum.UpdatedAt,
		Status:     forum.Status,
		IsPublic:   forum.IsPublic,
		AccessMode: forum.AccessMode,
		Course: getForumByCourseIDCourse{
			ID:          forum.Course.ID,
			Name:        forum.Course.Name,
			Description: forum.Course.Description,
		},
		Teacher: getForumByCourseIDTeacher{
			ID:       forum.Teacher.ID,
			FullName: forum.Teacher.FullName(),
			Phone:    forum.Teacher.Phone,
		},
		Members: make([]getForumByCourseIDMembers, len(forum.Course.Participants)),
	}
	for pIndex, pItem := range forum.Course.Participants {
		res.Members[pIndex] = getForumByCourseIDMembers{
			ID:       pItem.Student.ID,
			FullName: pItem.Student.FullName(),
		}
	}
	return res
}
