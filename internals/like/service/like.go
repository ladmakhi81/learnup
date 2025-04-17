package service

import (
	"github.com/ladmakhi81/learnup/db/entities"
	courseService "github.com/ladmakhi81/learnup/internals/course/service"
	dtoreq "github.com/ladmakhi81/learnup/internals/like/dto/req"
	"github.com/ladmakhi81/learnup/internals/like/repo"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
)

type LikeService interface {
	Create(authContext any, dto dtoreq.CreateLikeReq) (*entities.Like, error)
	FetchByCourseID(page, pageSize int, courseId uint) ([]*entities.Like, error)
	FetchCountByCourseID(courseId uint) (int, error)
}

type LikeServiceImpl struct {
	likeRepo       repo.LikeRepo
	userSvc        userService.UserSvc
	translationSvc contracts.Translator
	courseService  courseService.CourseService
}

func NewLikeServiceImpl(
	likeRepo repo.LikeRepo,
	userSvc userService.UserSvc,
	translationSvc contracts.Translator,
	courseService courseService.CourseService,
) *LikeServiceImpl {
	return &LikeServiceImpl{
		likeRepo:       likeRepo,
		userSvc:        userSvc,
		translationSvc: translationSvc,
		courseService:  courseService,
	}
}

func (svc LikeServiceImpl) Create(authContext any, dto dtoreq.CreateLikeReq) (*entities.Like, error) {
	authClaim := authContext.(*types.TokenClaim)
	user, userErr := svc.userSvc.FindById(authClaim.UserID)
	if userErr != nil {
		return nil, userErr
	}
	if user == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate("user.errors.not_found"),
		)
	}
	course, courseErr := svc.courseService.FindById(dto.CourseID)
	if courseErr != nil {
		return nil, courseErr
	}
	if course == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate("course.errors.not_found"),
		)
	}
	likedBefore, likedBeforeErr := svc.likeRepo.FindOne(repo.FindOneLikeOptions{
		UserID:   &user.ID,
		CourseID: &course.ID,
	})
	if likedBeforeErr != nil {
		return nil, types.NewServerError(
			"Error in finding like by user id and course id",
			"LikeServiceImpl.Create",
			likedBeforeErr,
		)
	}
	if likedBefore != nil {
		likedBefore.Type = dto.Type
		updateErr := svc.likeRepo.Update(likedBefore)
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
	if err := svc.likeRepo.Create(like); err != nil {
		return nil, types.NewServerError(
			"Error in creating like entity for course",
			"LikeServiceImpl.Create",
			err,
		)
	}
	return like, nil
}

func (svc LikeServiceImpl) FetchByCourseID(page, pageSize int, courseId uint) ([]*entities.Like, error) {
	likes, likesErr := svc.likeRepo.Fetch(
		repo.FetchLikeOptions{
			PageSize: &pageSize,
			Page:     &page,
			CourseID: &courseId,
			Preloads: []string{"User"},
		},
	)
	if likesErr != nil {
		return nil, types.NewServerError(
			"Error in fetching likes",
			"LikeServiceImpl.Fetch",
			likesErr,
		)
	}
	return likes, nil
}

func (svc LikeServiceImpl) FetchCountByCourseID(courseId uint) (int, error) {
	count, countErr := svc.likeRepo.FetchCount(
		repo.FetchCountLikeOptions{CourseID: &courseId},
	)
	if countErr != nil {
		return 0, types.NewServerError(
			"Error in fetching count of likes",
			"LikeServiceImpl.FetchCount",
			countErr,
		)
	}
	return count, nil
}
