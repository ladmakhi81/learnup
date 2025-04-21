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
	repo           *db.Repositories
	translationSvc contracts.Translator
}

func NewQuestionServiceImpl(
	repo *db.Repositories,
	translationSvc contracts.Translator,
) *QuestionServiceImpl {
	return &QuestionServiceImpl{
		repo:           repo,
		translationSvc: translationSvc,
	}
}

func (svc QuestionServiceImpl) Create(dto dtoreq.CreateQuestionReq) (*entities.Question, error) {
	sender, senderErr := svc.repo.UserRepo.GetByID(dto.UserID)
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
	course, courseErr := svc.repo.CourseRepo.GetByID(dto.CourseID)
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
		video, videoErr := svc.repo.VideoRepo.GetByID(*dto.VideoID)
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
	if err := svc.repo.QuestionRepo.Create(question); err != nil {
		return nil, types.NewServerError(
			"Error in creating question",
			"QuestionServiceImpl.Create",
			err,
		)
	}
	// TODO: notification system
	// send notification for teacher that we have new question
	return question, nil
}

func (svc QuestionServiceImpl) GetPageable(courseId *uint, page, pageSize int) ([]*entities.Question, int, error) {
	questions, count, questionErr := svc.repo.QuestionRepo.GetPaginated(
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
	return questions, count, nil
}
