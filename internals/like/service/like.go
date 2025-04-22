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
	tx, txErr := svc.unitOfWork.Begin()
	if txErr != nil {
		return nil, txErr
	}
	authClaim := authContext.(*types.TokenClaim)
	user, userErr := tx.UserRepo().GetByID(authClaim.UserID, nil)
	if userErr != nil {
		return nil, types.NewServerError(
			"Error in fetching logged in user",
			"LikeServiceImpl.Create",
			userErr,
		)
	}
	if user == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate("user.errors.not_found"),
		)
	}
	course, courseErr := tx.CourseRepo().GetByID(dto.CourseID, nil)
	if courseErr != nil {
		return nil, types.NewServerError(
			"Error in fetching course detail",
			"LikeServiceImpl.Create",
			courseErr,
		)
	}
	if course == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate("course.errors.not_found"),
		)
	}
	likedBefore, likedBeforeErr := tx.LikeRepo().GetOne(map[string]any{
		"user_id":   user.ID,
		"course_id": course.ID,
	}, nil)
	if likedBeforeErr != nil {
		return nil, types.NewServerError(
			"Error in finding like by user id and course id",
			"LikeServiceImpl.Create",
			likedBeforeErr,
		)
	}
	if likedBefore != nil {
		likedBefore.Type = dto.Type
		updateErr := tx.LikeRepo().Update(likedBefore)
		if updateErr != nil {
			return nil, types.NewServerError(
				"Error in updating like type",
				"LikeServiceImpl.Create",
				updateErr,
			)
		}
		return likedBefore, nil
	}
	like := &entities.Like{
		UserID:   user.ID,
		CourseID: course.ID,
		Type:     dto.Type,
	}
	if err := tx.LikeRepo().Create(like); err != nil {
		return nil, types.NewServerError(
			"Error in creating like entity for course",
			"LikeServiceImpl.Create",
			err,
		)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return like, nil
}

func (svc LikeServiceImpl) FetchByCourseID(page, pageSize int, courseId uint) ([]*entities.Like, int, error) {
	tx, txErr := svc.unitOfWork.Begin()
	if txErr != nil {
		return nil, 0, txErr
	}
	likes, count, likesErr := tx.LikeRepo().GetPaginated(repositories.GetPaginatedOptions{
		Offset: &page,
		Limit:  &pageSize,
		Conditions: map[string]any{
			"course_id": courseId,
		},
		Relations: []string{"User"},
	})
	if likesErr != nil {
		return nil, 0, types.NewServerError(
			"Error in fetching likes",
			"LikeServiceImpl.Fetch",
			likesErr,
		)
	}
	if err := tx.Commit(); err != nil {
		return nil, 0, err
	}
	return likes, count, nil
}
