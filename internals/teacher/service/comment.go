package service

import (
	"github.com/ladmakhi81/learnup/db/entities"
	commentRepo "github.com/ladmakhi81/learnup/internals/comment/repo"
	courseService "github.com/ladmakhi81/learnup/internals/course/service"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
)

type TeacherCommentService interface {
	GetPageableCommentByCourseId(authContext any, courseId uint, page, pageSize int) ([]*entities.Comment, error)
	GetCommentCountByCourseId(authContext any, courseId uint) (int, error)
}

type TeacherCommentServiceImpl struct {
	userSvc        userService.UserSvc
	courseSvc      courseService.CourseService
	commentRepo    commentRepo.CommentRepo
	translationSvc contracts.Translator
}

func NewTeacherCommentServiceImpl(
	userSvc userService.UserSvc,
	courseSvc courseService.CourseService,
	commentRepo commentRepo.CommentRepo,
	translationSvc contracts.Translator,
) *TeacherCommentServiceImpl {
	return &TeacherCommentServiceImpl{
		userSvc:        userSvc,
		courseSvc:      courseSvc,
		commentRepo:    commentRepo,
		translationSvc: translationSvc,
	}
}

func (svc TeacherCommentServiceImpl) GetPageableCommentByCourseId(authContext any, courseId uint, page, pageSize int) ([]*entities.Comment, error) {
	teacherClaim := authContext.(*types.TokenClaim)
	teacher, teacherErr := svc.userSvc.FindById(teacherClaim.UserID)
	if teacherErr != nil {
		return nil, teacherErr
	}
	if teacher == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate("user.errors.teacher_not_found"),
		)
	}
	course, courseErr := svc.courseSvc.FindById(courseId)
	if courseErr != nil {
		return nil, courseErr
	}
	if course == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate("course.errors.not_found"),
		)
	}
	comments, commentsErr := svc.commentRepo.Fetch(commentRepo.FetchCommentOption{
		PageSize: &pageSize,
		Page:     &page,
		Preloads: []string{"User", "Course"},
		UserID:   &teacher.ID,
		CourseID: &course.ID,
	})
	if commentsErr != nil {
		return nil, types.NewServerError(
			"Error in fetching comments",
			"TeacherCommentServiceImpl.GetPageableCommentByCourseId",
			commentsErr,
		)
	}
	return comments, nil
}

func (svc TeacherCommentServiceImpl) GetCommentCountByCourseId(authContext any, courseId uint) (int, error) {
	teacherClaim := authContext.(*types.TokenClaim)
	teacher, teacherErr := svc.userSvc.FindById(teacherClaim.UserID)
	if teacherErr != nil {
		return 0, teacherErr
	}
	if teacher == nil {
		return 0, types.NewNotFoundError(
			svc.translationSvc.Translate("user.errors.teacher_not_found"),
		)
	}
	course, courseErr := svc.courseSvc.FindById(courseId)
	if courseErr != nil {
		return 0, courseErr
	}
	if course == nil {
		return 0, types.NewNotFoundError(
			svc.translationSvc.Translate("course.errors.not_found"),
		)
	}
	count, countErr := svc.commentRepo.FetchCount(commentRepo.FetchCountCommentOption{
		CourseID: &course.ID,
		UserID:   &teacher.ID,
	})
	if countErr != nil {
		return 0, types.NewServerError(
			"Error in fetching count of comments based on course id and user id",
			"TeacherCommentServiceImpl.GetCommentCountByCourseId",
			countErr,
		)
	}
	return count, nil
}
