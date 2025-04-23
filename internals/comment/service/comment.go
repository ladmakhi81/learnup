package service

import (
	dtoreq "github.com/ladmakhi81/learnup/internals/comment/dto/req"
	commentError "github.com/ladmakhi81/learnup/internals/comment/error"
	courseError "github.com/ladmakhi81/learnup/internals/course/error"
	"github.com/ladmakhi81/learnup/internals/db"
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"github.com/ladmakhi81/learnup/internals/db/repositories"
	"github.com/ladmakhi81/learnup/types"
)

type CommentService interface {
	Create(user *entities.User, dto dtoreq.CreateCommentReqDto) (*entities.Comment, error)
	Delete(id uint) error
	Fetch(page, pageSize int) ([]*entities.Comment, int, error)
}

type commentService struct {
	unitOfWork db.UnitOfWork
}

func NewCommentSvc(unitOfWork db.UnitOfWork) CommentService {
	return &commentService{unitOfWork: unitOfWork}
}

func (svc commentService) Create(user *entities.User, dto dtoreq.CreateCommentReqDto) (*entities.Comment, error) {
	const operationName = "commentService.Create"
	course, err := svc.unitOfWork.CourseRepo().GetByID(dto.CourseId, nil)
	if err != nil {
		return nil, types.NewServerError("Error in fetching course", operationName, err)
	}
	if course == nil {
		return nil, courseError.Course_NotFound
	}
	if dto.ParentId != nil {
		parent, err := svc.unitOfWork.CommentRepo().GetByID(*dto.ParentId, nil)
		if err != nil {
			return nil, types.NewServerError("Error in fetching comment by parent id", operationName, err)
		}
		if parent == nil {
			return nil, commentError.Comment_ParentNotFound
		}
	}
	comment := &entities.Comment{
		Content:         dto.Content,
		UserID:          &user.ID,
		CourseID:        &course.ID,
		ParentCommentId: dto.ParentId,
	}
	if err := svc.unitOfWork.CommentRepo().Create(comment); err != nil {
		return nil, types.NewServerError("Error in creating comment", operationName, err)
	}
	return comment, nil
}

func (svc commentService) Delete(id uint) error {
	const operationName = "commentService.Delete"
	comment, err := svc.unitOfWork.CommentRepo().GetByID(id, nil)
	if err != nil {
		return types.NewServerError("Error in fetching comment by id", operationName, err)
	}
	if comment == nil {
		return commentError.Comment_NotFound
	}
	if err := svc.unitOfWork.CommentRepo().Delete(comment); err != nil {
		return types.NewServerError("Error in deleting comment", operationName, err)
	}
	return nil
}

func (svc commentService) Fetch(page, pageSize int) ([]*entities.Comment, int, error) {
	const operationName = "commentService.Fetch"
	comments, count, err := svc.unitOfWork.CommentRepo().GetPaginated(repositories.GetPaginatedOptions{
		Offset:    &page,
		Limit:     &pageSize,
		Relations: []string{"User", "Course"},
	})
	if err != nil {
		return nil, 0, types.NewServerError("Error in fetching comments", operationName, err)
	}
	return comments, count, nil
}
