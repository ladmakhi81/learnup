package service

import (
	courseError "github.com/ladmakhi81/learnup/internals/course/error"
	dtoreq "github.com/ladmakhi81/learnup/internals/like/dto/req"
	"github.com/ladmakhi81/learnup/shared/db"
	"github.com/ladmakhi81/learnup/shared/db/entities"
	"github.com/ladmakhi81/learnup/shared/db/repositories"
	"github.com/ladmakhi81/learnup/shared/types"
)

type LikeService interface {
	Create(user *entities.User, dto dtoreq.CreateLikeReqDto) (*entities.Like, error)
	FetchByCourseID(page, pageSize int, courseId uint) ([]*entities.Like, int, error)
}

type likeService struct {
	unitOfWork db.UnitOfWork
}

func NewLikeSvc(unitOfWork db.UnitOfWork) LikeService {
	return &likeService{unitOfWork: unitOfWork}
}

func (svc likeService) Create(user *entities.User, dto dtoreq.CreateLikeReqDto) (*entities.Like, error) {
	const operationName = "likeService.Create"
	course, err := svc.unitOfWork.CourseRepo().GetByID(dto.CourseID, nil)
	if err != nil {
		return nil, types.NewServerError("Error in fetching course detail", operationName, err)
	}
	if course == nil {
		return nil, courseError.Course_NotFound
	}
	likedBefore, err := svc.unitOfWork.LikeRepo().GetOne(map[string]any{"user_id": user.ID, "course_id": course.ID}, nil)
	if err != nil {
		return nil, types.NewServerError("Error in finding like by user id and course id", operationName, err)
	}
	if likedBefore != nil {
		likedBefore.Type = dto.Type
		if err := svc.unitOfWork.LikeRepo().Update(likedBefore); err != nil {
			return nil, types.NewServerError("Error in updating like type", operationName, err)
		}
		return likedBefore, nil
	}
	like := &entities.Like{
		UserID:   user.ID,
		CourseID: course.ID,
		Type:     dto.Type,
	}
	if err := svc.unitOfWork.LikeRepo().Create(like); err != nil {
		return nil, types.NewServerError("Error in creating like entity for course", operationName, err)
	}
	return like, nil
}

func (svc likeService) FetchByCourseID(page, pageSize int, courseId uint) ([]*entities.Like, int, error) {
	const operationName = "likeService.FetchByCourseID"
	likes, count, err := svc.unitOfWork.LikeRepo().GetPaginated(repositories.GetPaginatedOptions{
		Offset: &page,
		Limit:  &pageSize,
		Conditions: map[string]any{
			"course_id": courseId,
		},
		Relations: []string{"User"},
	})
	if err != nil {
		return nil, 0, types.NewServerError("Error in fetching likes", operationName, err)
	}
	return likes, count, nil
}
