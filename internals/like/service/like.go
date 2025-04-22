package service

import (
	"github.com/ladmakhi81/learnup/internals/db"
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"github.com/ladmakhi81/learnup/internals/db/repositories"
	dtoreq "github.com/ladmakhi81/learnup/internals/like/dto/req"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
)

type LikeService interface {
	Create(authContext any, dto dtoreq.CreateLikeReq) (*entities.Like, error)
	FetchByCourseID(page, pageSize int, courseId uint) ([]*entities.Like, int, error)
}

type LikeServiceImpl struct {
	unitOfWork     db.UnitOfWork
	translationSvc contracts.Translator
}

func NewLikeServiceImpl(
	unitOfWork db.UnitOfWork,
	translationSvc contracts.Translator,
) *LikeServiceImpl {
	return &LikeServiceImpl{
		unitOfWork:     unitOfWork,
		translationSvc: translationSvc,
	}
}

func (svc LikeServiceImpl) Create(authContext any, dto dtoreq.CreateLikeReq) (*entities.Like, error) {
	const operationName = "LikeServiceImpl.Create"
	authClaim := authContext.(*types.TokenClaim)
	user, err := svc.unitOfWork.UserRepo().GetByID(authClaim.UserID, nil)
	if err != nil {
		return nil, types.NewServerError(
			"Error in fetching logged in user",
			operationName,
			err,
		)
	}
	if user == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate("user.errors.not_found"),
		)
	}
	course, err := svc.unitOfWork.CourseRepo().GetByID(dto.CourseID, nil)
	if err != nil {
		return nil, types.NewServerError(
			"Error in fetching course detail",
			operationName,
			err,
		)
	}
	if course == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate("course.errors.not_found"),
		)
	}
	likedBefore, err := svc.unitOfWork.LikeRepo().GetOne(map[string]any{
		"user_id":   user.ID,
		"course_id": course.ID,
	}, nil)
	if err != nil {
		return nil, types.NewServerError(
			"Error in finding like by user id and course id",
			operationName,
			err,
		)
	}
	if likedBefore != nil {
		likedBefore.Type = dto.Type
		if err := svc.unitOfWork.LikeRepo().Update(likedBefore); err != nil {
			return nil, types.NewServerError(
				"Error in updating like type",
				operationName,
				err,
			)
		}
		return likedBefore, nil
	}
	like := &entities.Like{
		UserID:   user.ID,
		CourseID: course.ID,
		Type:     dto.Type,
	}
	if err := svc.unitOfWork.LikeRepo().Create(like); err != nil {
		return nil, types.NewServerError(
			"Error in creating like entity for course",
			operationName,
			err,
		)
	}
	return like, nil
}

func (svc LikeServiceImpl) FetchByCourseID(page, pageSize int, courseId uint) ([]*entities.Like, int, error) {
	const operationName = "LikeServiceImpl.FetchByCourseID"
	likes, count, err := svc.unitOfWork.LikeRepo().GetPaginated(repositories.GetPaginatedOptions{
		Offset: &page,
		Limit:  &pageSize,
		Conditions: map[string]any{
			"course_id": courseId,
		},
		Relations: []string{"User"},
	})
	if err != nil {
		return nil, 0, types.NewServerError(
			"Error in fetching likes",
			operationName,
			err,
		)
	}
	return likes, count, nil
}
