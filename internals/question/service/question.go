package service

import (
	"github.com/ladmakhi81/learnup/internals/db"
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"github.com/ladmakhi81/learnup/internals/db/repositories"
	dtoreq "github.com/ladmakhi81/learnup/internals/question/dto/req"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
)

type QuestionService interface {
	Create(dto dtoreq.CreateQuestionReq) (*entities.Question, error)
	GetPageable(courseId *uint, page, pageSize int) ([]*entities.Question, int, error)
}

type QuestionServiceImpl struct {
	unitOfWork     db.UnitOfWork
	translationSvc contracts.Translator
}

func NewQuestionServiceImpl(
	unitOfWork db.UnitOfWork,
	translationSvc contracts.Translator,
) *QuestionServiceImpl {
	return &QuestionServiceImpl{
		unitOfWork:     unitOfWork,
		translationSvc: translationSvc,
	}
}

func (svc QuestionServiceImpl) Create(dto dtoreq.CreateQuestionReq) (*entities.Question, error) {
	tx, txErr := svc.unitOfWork.Begin()
	if txErr != nil {
		return nil, txErr
	}
	sender, senderErr := tx.UserRepo().GetByID(dto.UserID, nil)
	if senderErr != nil {
		return nil, types.NewServerError(
			"Error in fetching sender data",
			"QuestionServiceImpl.Create",
			senderErr,
		)
	}
	if sender == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate("user.errors.not_found"),
		)
	}
	course, courseErr := tx.CourseRepo().GetByID(dto.CourseID, nil)
	if courseErr != nil {
		return nil, types.NewServerError(
			"Error in fetching course by id",
			"QuestionServiceImpl.Create",
			courseErr,
		)
	}
	if course == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate("course.errors.not_found"),
		)
	}
	question := &entities.Question{
		UserID:   sender.ID,
		CourseID: course.ID,
		Content:  dto.Content,
		Priority: dto.Priority,
	}
	if dto.VideoID != nil {
		video, videoErr := tx.VideoRepo().GetByID(*dto.VideoID, nil)
		if videoErr != nil {
			return nil, types.NewServerError(
				"Error in fetching video",
				"QuestionServiceImpl.Create",
				videoErr,
			)
		}
		if video == nil {
			return nil, types.NewNotFoundError(
				svc.translationSvc.Translate("video.errors.not_found"),
			)
		}
		question.VideoID = &video.ID
	}
	if err := tx.QuestionRepo().Create(question); err != nil {
		return nil, types.NewServerError(
			"Error in creating question",
			"QuestionServiceImpl.Create",
			err,
		)
	}
	// TODO: notification system
	// send notification for teacher that we have new question
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return question, nil
}

func (svc QuestionServiceImpl) GetPageable(courseId *uint, page, pageSize int) ([]*entities.Question, int, error) {
	tx, txErr := svc.unitOfWork.Begin()
	if txErr != nil {
		return nil, 0, txErr
	}
	questions, count, questionErr := tx.QuestionRepo().GetPaginated(
		repositories.GetPaginatedOptions{
			Limit:  &pageSize,
			Offset: &page,
			Conditions: map[string]any{
				"course_id": courseId,
			},
			Relations: []string{
				"User",
				"Course",
				"Video",
			},
		})
	if questionErr != nil {
		return nil, 0, types.NewServerError(
			"Error in fetching related questions",
			"QuestionServiceImpl.GetPageable",
			questionErr,
		)
	}
	if err := tx.Commit(); err != nil {
		return nil, 0, err
	}
	return questions, count, nil
}
