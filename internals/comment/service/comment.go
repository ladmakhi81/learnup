package service

import (
	dtoreq "github.com/ladmakhi81/learnup/internals/comment/dto/req"
	"github.com/ladmakhi81/learnup/internals/db"
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"github.com/ladmakhi81/learnup/internals/db/repositories"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
)

type CommentService interface {
	Create(authContext any, dto dtoreq.CreateCommentReq) (*entities.Comment, error)
	Delete(id uint) error
	Fetch(page, pageSize int) ([]*entities.Comment, int, error)
}

type CommentServiceImpl struct {
	repo           *db.Repositories
	translationSvc contracts.Translator
}

func NewCommentServiceImpl(
	repo *db.Repositories,
	translationSvc contracts.Translator,
) *CommentServiceImpl {
	return &CommentServiceImpl{
		repo:           repo,
		translationSvc: translationSvc,
	}
}

func (svc CommentServiceImpl) Create(authContext any, dto dtoreq.CreateCommentReq) (*entities.Comment, error) {
	authClaim := authContext.(*types.TokenClaim)
	user, userErr := svc.repo.UserRepo.GetByID(authClaim.UserID, nil)
	if userErr != nil {
		return nil, types.NewServerError(
			"Error in fetching user logged in information",
			"CommentServiceImpl.Create",
			userErr,
		)
	}
	if user == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate(
				"comment.errors.sender_not_found",
			),
		)
	}
	course, courseErr := svc.repo.CourseRepo.GetByID(dto.CourseId, nil)
	if courseErr != nil {
		return nil, types.NewServerError(
			"Error in fetching course",
			"CommentServiceImpl.Create",
			courseErr,
		)
	}
	if course == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate(
				"comment.errors.course_not_found",
			),
		)
	}
	if dto.ParentId != nil {
		parent, parentErr := svc.repo.CommentRepo.GetByID(*dto.ParentId, nil)
		if parentErr != nil {
			return nil, types.NewServerError(
				"Error in fetching comment by parent id",
				"CommentServiceImpl.Create",
				parentErr,
			)
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
	if err := svc.repo.CommentRepo.Create(comment); err != nil {
		return nil, types.NewServerError(
			"Error in creating comment",
			"CommentServiceImpl.Create",
			err,
		)
	}
	return comment, nil
}

func (svc CommentServiceImpl) Delete(id uint) error {
	comment, commentErr := svc.repo.CommentRepo.GetByID(id, nil)
	if commentErr != nil {
		return types.NewServerError(
			"Error in fetching comment by id",
			"CommentServiceImpl.Delete",
			commentErr,
		)
	}
	if comment == nil {
		return types.NewNotFoundError(
			svc.translationSvc.Translate("comment.errors.not_found"),
		)
	}
	if err := svc.repo.CommentRepo.Delete(comment); err != nil {
		return types.NewServerError(
			"Error in deleting comment",
			"CommentServiceImpl.Delete",
			err,
		)
	}
	return nil
}

func (svc CommentServiceImpl) Fetch(page, pageSize int) ([]*entities.Comment, int, error) {
	comments, count, commentsErr := svc.repo.CommentRepo.GetPaginated(repositories.GetPaginatedOptions{
		Offset:    &page,
		Limit:     &pageSize,
		Relations: []string{"User", "Course"},
	})
	if commentsErr != nil {
		return nil, 0, types.NewServerError(
			"Error in fetching comments",
			"CommentServiceImpl.Fetch",
			commentsErr,
		)
	}
	return comments, count, nil
}
