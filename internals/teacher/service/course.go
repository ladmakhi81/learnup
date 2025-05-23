package service

import (
	courseError "github.com/ladmakhi81/learnup/internals/course/error"
	teacherDtoReq "github.com/ladmakhi81/learnup/internals/teacher/dto/req"
	"github.com/ladmakhi81/learnup/shared/db"
	"github.com/ladmakhi81/learnup/shared/db/entities"
	"github.com/ladmakhi81/learnup/shared/db/repositories"
	"github.com/ladmakhi81/learnup/shared/types"
)

type TeacherCourseService interface {
	Create(teacher *entities.User, dto teacherDtoReq.CreateCourseReqDto) (*entities.Course, error)
	FetchByTeacherId(teacher *entities.User, page, pageSize int) ([]*entities.Course, int, error)
}

type teacherCourseService struct {
	unitOfWork db.UnitOfWork
}

func NewTeacherCourseService(unitOfWork db.UnitOfWork) TeacherCourseService {
	return &teacherCourseService{unitOfWork: unitOfWork}
}

func (svc teacherCourseService) Create(teacher *entities.User, dto teacherDtoReq.CreateCourseReqDto) (*entities.Course, error) {
	return db.WithTx(svc.unitOfWork, func(tx db.UnitOfWorkTx) (*entities.Course, error) {
		const operationName = "teacherCourseService.Create"
		isDuplicate, err := tx.CourseRepo().Exist(map[string]any{"name": dto.Name})
		if err != nil {
			return nil, types.NewServerError("Error in checking existence of course name", operationName, err)
		}
		if isDuplicate {
			return nil, courseError.Course_NameDuplicated
		}
		category, err := tx.CategoryRepo().GetByID(dto.CategoryID, nil)
		if err != nil {
			return nil, types.NewServerError("Error in fetching category by id", operationName, err)
		}
		if category == nil {
			return nil, courseError.Course_NotFoundCategory
		}
		course := &entities.Course{
			Name:                dto.Name,
			TeacherID:           &teacher.ID,
			CategoryID:          &category.ID,
			AbilityToAddComment: dto.AbilityToAddComment,
			CanHaveDiscount:     dto.CanHaveDiscount,
			CommentAccessMode:   dto.CommentAccessMode,
			Description:         dto.Description,
			Image:               dto.Image,
			Level:               dto.Level,
			IsPublished:         false,
			IntroductionVideo:   dto.IntroductionVideo,
			MaxDiscountAmount:   dto.MaxDiscountAmount,
			Price:               dto.Price,
			Prerequisite:        dto.Prerequisite,
			Tags:                dto.Tags,
			ThumbnailImage:      dto.ThumbnailImage,
			Status:              entities.CourseStatus_InProgress,
		}
		if err := tx.CourseRepo().Create(course); err != nil {
			return nil, types.NewServerError("Error in creating teacher course", operationName, err)
		}
		forum := &entities.CourseForum{
			TeacherID:  teacher.ID,
			CourseID:   course.ID,
			AccessMode: dto.ForumAccessMode,
			Status:     dto.ForumStatus,
			IsPublic:   dto.IsForumPublic,
		}
		if err := tx.CourseForumRepo().Create(forum); err != nil {
			return nil, types.NewServerError("Error in creating course forum", operationName, err)
		}
		// TODO: notification system
		//send email for creating course
		//create notification about creating new course for all admin
		return course, nil
	})
}

func (svc teacherCourseService) FetchByTeacherId(teacher *entities.User, page, pageSize int) ([]*entities.Course, int, error) {
	const operationName = "teacherCourseService.FetchByTeacherId"
	courses, count, err := svc.unitOfWork.CourseRepo().GetPaginated(repositories.GetPaginatedOptions{
		Offset: &page,
		Limit:  &pageSize,
		Conditions: map[string]any{
			"teacher_id": teacher.ID,
		},
	})
	if err != nil {
		return nil, 0, types.NewServerError("Error in fetching courses related to teacher", operationName, err)
	}
	return courses, count, nil
}
