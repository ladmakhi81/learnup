package service

import (
	"github.com/ladmakhi81/learnup/db/entities"
	courseService "github.com/ladmakhi81/learnup/internals/course/service"
	dtoreq "github.com/ladmakhi81/learnup/internals/question/dto/req"
	"github.com/ladmakhi81/learnup/internals/question/repo"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	videoService "github.com/ladmakhi81/learnup/internals/video/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
)

type QuestionService interface {
	Create(dto dtoreq.CreateQuestionReq) (*entities.Question, error)
	GetPageable(courseId *uint, page, pageSize int) ([]*entities.Question, error)
	GetCount(courseId *uint) (int, error)
	FindById(id uint) (*entities.Question, error)
}

type QuestionServiceImpl struct {
	questionRepo   repo.QuestionRepo
	userSvc        userService.UserSvc
	courseSvc      courseService.CourseService
	translationSvc contracts.Translator
	videoSvc       videoService.VideoService
}

func NewQuestionServiceImpl(
	questionRepo repo.QuestionRepo,
	userSvc userService.UserSvc,
	courseSvc courseService.CourseService,
	translationSvc contracts.Translator,
	videoSvc videoService.VideoService,
) *QuestionServiceImpl {
	return &QuestionServiceImpl{
		userSvc:        userSvc,
		courseSvc:      courseSvc,
		translationSvc: translationSvc,
		videoSvc:       videoSvc,
		questionRepo:   questionRepo,
	}
}

func (svc QuestionServiceImpl) Create(dto dtoreq.CreateQuestionReq) (*entities.Question, error) {
	sender, senderErr := svc.userSvc.FindById(dto.UserID)
	if senderErr != nil {
		return nil, senderErr
	}
	if sender == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate("user.errors.not_found"),
		)
	}
	course, courseErr := svc.courseSvc.FindById(dto.CourseID)
	if courseErr != nil {
		return nil, courseErr
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
		video, videoErr := svc.videoSvc.FindById(*dto.VideoID)
		if videoErr != nil {
			return nil, videoErr
		}
		if video == nil {
			return nil, types.NewNotFoundError(
				svc.translationSvc.Translate("video.errors.not_found"),
			)
		}
		question.VideoID = &video.ID
	}
	if err := svc.questionRepo.Create(question); err != nil {
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

func (svc QuestionServiceImpl) GetPageable(courseId *uint, page, pageSize int) ([]*entities.Question, error) {
	questions, questionsErr := svc.questionRepo.Fetch(repo.FetchQuestionOptions{
		PageSize: &pageSize,
		Page:     &page,
		CourseID: courseId,
	})
	if questionsErr != nil {
		return nil, types.NewServerError(
			"Error in fetching related questions",
			"QuestionServiceImpl.GetPageable",
			questionsErr,
		)
	}
	return questions, nil
}

func (svc QuestionServiceImpl) GetCount(courseId *uint) (int, error) {
	count, countErr := svc.questionRepo.FetchCount(repo.FetchCountQuestionOptions{
		CourseID: courseId,
	})
	if countErr != nil {
		return 0, types.NewServerError(
			"Error in fetching count of question",
			"QuestionServiceImpl.GetCount",
			countErr,
		)
	}
	return count, nil
}

func (svc QuestionServiceImpl) FindById(id uint) (*entities.Question, error) {
	question, questionErr := svc.questionRepo.FindOne(
		repo.FetchOneQuestionOptions{ID: &id},
	)
	if questionErr != nil {
		return nil, questionErr
	}
	return question, nil
}
