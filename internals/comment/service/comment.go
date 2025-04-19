package service

import (
	"github.com/ladmakhi81/learnup/db/entities"
	dtoreq "github.com/ladmakhi81/learnup/internals/comment/dto/req"
	commentRepo "github.com/ladmakhi81/learnup/internals/comment/repo"
	courseService "github.com/ladmakhi81/learnup/internals/course/service"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
)

type CommentService interface {
	Create(authContext any, dto dtoreq.CreateCommentReq) (*entities.Comment, error)
	FindById(id uint) (*entities.Comment, error)
	Delete(id uint) error
	Fetch(page, pageSize int) ([]*entities.Comment, error)
	FetchCount() (int, error)
}

type CommentServiceImpl struct {
	commentRepo    commentRepo.CommentRepo
	userSvc        userService.UserSvc
	courseSvc      courseService.CourseService
	translationSvc contracts.Translator
}

func NewCommentServiceImpl(
	commentRepo commentRepo.CommentRepo,
	userSvc userService.UserSvc,
	courseSvc courseService.CourseService,
	translationSvc contracts.Translator,
) *CommentServiceImpl {
	return &CommentServiceImpl{
		commentRepo:    commentRepo,
		userSvc:        userSvc,
		courseSvc:      courseSvc,
		translationSvc: translationSvc,
	}
}

func (svc CommentServiceImpl) Create(authContext any, dto dtoreq.CreateCommentReq) (*entities.Comment, error) {
	authClaim := authContext.(*types.TokenClaim)
	user, userErr := svc.userSvc.FindById(authClaim.UserID)
	if userErr != nil {
		return nil, userErr
	}
	if user == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate(
				"comment.errors.sender_not_found",
			),
		)
	}
	course, courseErr := svc.courseSvc.FindById(dto.CourseId)
	if courseErr != nil {
		return nil, courseErr
	}
	if course == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate(
				"comment.errors.course_not_found",
			),
		)
	}
	if dto.ParentId != nil {
		parent, parentErr := svc.FindById(*dto.ParentId)
		if parentErr != nil {
			return nil, parentErr
		}
		if parent == nil {
			return nil, types.NewNotFoundError(
				svc.translationSvc.Translate(
					"comment.errors.parent_comment_not_found",
				),
			)
		}
	}
	comment := &entities.Comment{
		Content:         dto.Content,
		UserID:          &user.ID,
		CourseID:        &course.ID,
		ParentCommentId: dto.ParentId,
	}
	if err := svc.commentRepo.Create(comment); err != nil {
		return nil, types.NewServerError(
			"Error in creating comment",
			"CommentServiceImpl.Create",
			err,
		)
	}
	return comment, nil
}

func (svc CommentServiceImpl) FindById(id uint) (*entities.Comment, error) {
	comment, commentErr := svc.commentRepo.FindById(id)
	if commentErr != nil {
		return nil, types.NewServerError(
			"Error in finding comment by id",
			"CommentServiceImpl.FetchById",
			commentErr,
		)
	}
	return comment, nil
}

func (svc CommentServiceImpl) Delete(id uint) error {
	comment, commentErr := svc.FindById(id)
	if commentErr != nil {
		return commentErr
	}
	if comment == nil {
		return types.NewNotFoundError(
			svc.translationSvc.Translate("comment.errors.not_found"),
		)
	}
	if err := svc.commentRepo.Delete(comment.ID); err != nil {
		return types.NewServerError(
			"Error in deleting comment",
			"CommentServiceImpl.Delete",
			err,
		)
	}
	return nil
}

func (svc CommentServiceImpl) Fetch(page, pageSize int) ([]*entities.Comment, error) {
	comments, commentsErr := svc.commentRepo.Fetch(
		commentRepo.FetchCommentOption{
			PageSize: &pageSize,
			Page:     &page,
			Preloads: []string{"User", "Course"},
		},
	)
	if commentsErr != nil {
		return nil, types.NewServerError(
			"Error in fetching comments",
			"CommentServiceImpl.Fetch",
			commentsErr,
		)
	}
	return comments, nil
}

func (svc CommentServiceImpl) FetchCount() (int, error) {
	count, countErr := svc.commentRepo.FetchCount(
		commentRepo.FetchCountCommentOption{},
	)
	if countErr != nil {
		return 0, types.NewServerError(
			"Error in fetching count of comments",
			"CommentServiceImpl.FetchCount",
			countErr,
		)
	}
	return count, nil
}
