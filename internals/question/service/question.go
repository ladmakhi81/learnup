package service

import (
	courseError "github.com/ladmakhi81/learnup/internals/course/error"
	"github.com/ladmakhi81/learnup/internals/db"
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"github.com/ladmakhi81/learnup/internals/db/repositories"
	dtoreq "github.com/ladmakhi81/learnup/internals/question/dto/req"
	userError "github.com/ladmakhi81/learnup/internals/user/error"
	videoError "github.com/ladmakhi81/learnup/internals/video/error"
	"github.com/ladmakhi81/learnup/types"
)

type QuestionService interface {
	Create(dto dtoreq.CreateQuestionReq) (*entities.Question, error)
	GetPageable(courseId *uint, page, pageSize int) ([]*entities.Question, int, error)
}

type questionService struct {
	unitOfWork db.UnitOfWork
}

func NewQuestionSvc(unitOfWork db.UnitOfWork) QuestionService {
	return &questionService{unitOfWork: unitOfWork}
}

func (svc questionService) Create(dto dtoreq.CreateQuestionReq) (*entities.Question, error) {
	const operationName = "questionService.Create"
	sender, err := svc.unitOfWork.UserRepo().GetByID(dto.UserID, nil)
	if err != nil {
		return nil, types.NewServerError("Error in fetching sender data", operationName, err)
	}
	if sender == nil {
		return nil, userError.User_NotFound
	}
	course, err := svc.unitOfWork.CourseRepo().GetByID(dto.CourseID, nil)
	if err != nil {
		return nil, types.NewServerError("Error in fetching course by id", operationName, err)
	}
	if course == nil {
		return nil, courseError.Course_NotFound
	}
	question := &entities.Question{
		UserID:   sender.ID,
		CourseID: course.ID,
		Content:  dto.Content,
		Priority: dto.Priority,
	}
	if dto.VideoID != nil {
		video, err := svc.unitOfWork.VideoRepo().GetByID(*dto.VideoID, nil)
		if err != nil {
			return nil, types.NewServerError("Error in fetching video", operationName, err)
		}
		if video == nil {
			return nil, videoError.Video_NotFound
		}
		question.VideoID = &video.ID
	}
	if err := svc.unitOfWork.QuestionRepo().Create(question); err != nil {
		return nil, types.NewServerError("Error in creating question", operationName, err)
	}
	// TODO: notification system
	// send notification for teacher that we have new question
	return question, nil
}

func (svc questionService) GetPageable(courseId *uint, page, pageSize int) ([]*entities.Question, int, error) {
	const operationName = "questionService.GetPageable"
	questions, count, err := svc.unitOfWork.QuestionRepo().GetPaginated(
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
	if err != nil {
		return nil, 0, types.NewServerError("Error in fetching related questions", operationName, err)
	}
	return questions, count, nil
}
