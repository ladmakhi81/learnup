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
	unitOfWork     db.UnitOfWork
	translationSvc contracts.Translator
}

func NewCommentServiceImpl(
	unitOfWork db.UnitOfWork,
	translationSvc contracts.Translator,
) *CommentServiceImpl {
	return &CommentServiceImpl{
		unitOfWork:     unitOfWork,
		translationSvc: translationSvc,
	}
}

func (svc CommentServiceImpl) Create(authContext any, dto dtoreq.CreateCommentReq) (*entities.Comment, error) {
	const operationName = "CommentServiceImpl.Create"
	authClaim := authContext.(*types.TokenClaim)
	user, err := svc.unitOfWork.UserRepo().GetByID(authClaim.UserID, nil)
	if err != nil {
		return nil, types.NewServerError(
			"Error in fetching user logged in information",
			operationName,
			err,
		)
	}
	if user == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate(
				"comment.errors.sender_not_found",
			),
		)
	}
	course, err := svc.unitOfWork.CourseRepo().GetByID(dto.CourseId, nil)
	if err != nil {
		return nil, types.NewServerError(
			"Error in fetching course",
			operationName,
			err,
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
		parent, err := svc.unitOfWork.CommentRepo().GetByID(*dto.ParentId, nil)
		if err != nil {
			return nil, types.NewServerError(
				"Error in fetching comment by parent id",
				operationName,
				err,
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
	if err := svc.unitOfWork.CommentRepo().Create(comment); err != nil {
		return nil, types.NewServerError(
			"Error in creating comment",
			operationName,
			err,
		)
	}
	return comment, nil
}

func (svc CommentServiceImpl) Delete(id uint) error {
	const operationName = "CommentServiceImpl.Delete"
	comment, err := svc.unitOfWork.CommentRepo().GetByID(id, nil)
	if err != nil {
		return types.NewServerError(
			"Error in fetching comment by id",
			operationName,
			err,
		)
	}
	if comment == nil {
		return types.NewNotFoundError(
			svc.translationSvc.Translate("comment.errors.not_found"),
		)
	}
	if err := svc.unitOfWork.CommentRepo().Delete(comment); err != nil {
		return types.NewServerError(
			"Error in deleting comment",
			operationName,
			err,
		)
	}
	return nil
}

func (svc CommentServiceImpl) Fetch(page, pageSize int) ([]*entities.Comment, int, error) {
	const operationName = "CommentServiceImpl.Fetch"
	comments, count, err := svc.unitOfWork.CommentRepo().GetPaginated(repositories.GetPaginatedOptions{
		Offset:    &page,
		Limit:     &pageSize,
		Relations: []string{"User", "Course"},
	})
	if err != nil {
		return nil, 0, types.NewServerError(
			"Error in fetching comments",
			operationName,
			err,
		)
	}
	return comments, count, nil
}
