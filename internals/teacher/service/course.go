package service

import (
	courseError "github.com/ladmakhi81/learnup/internals/course/error"
	"github.com/ladmakhi81/learnup/internals/db"
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"github.com/ladmakhi81/learnup/internals/db/repositories"
	teacherDtoReq "github.com/ladmakhi81/learnup/internals/teacher/dto/req"
	userError "github.com/ladmakhi81/learnup/internals/user/error"
	"github.com/ladmakhi81/learnup/types"
)

type TeacherCourseService interface {
	Create(authContext any, dto teacherDtoReq.CreateCourseReq) (*entities.Course, error)
	FetchByTeacherId(authContext any, page, pageSize int) ([]*entities.Course, int, error)
}

type teacherCourseService struct {
	unitOfWork db.UnitOfWork
}

func NewTeacherCourseService(unitOfWork db.UnitOfWork) TeacherCourseService {
	return &teacherCourseService{unitOfWork: unitOfWork}
}

func (svc teacherCourseService) Create(authContext any, dto teacherDtoReq.CreateCourseReq) (*entities.Course, error) {
	const operationName = "teacherCourseService.Create"
	isDuplicate, err := svc.unitOfWork.CourseRepo().Exist(map[string]any{"name": dto.Name})
	if err != nil {
		return nil, types.NewServerError("Error in checking existence of course name", operationName, err)
	}
	if isDuplicate {
		return nil, courseError.Course_NameDuplicated
	}
	category, err := svc.unitOfWork.CategoryRepo().GetByID(dto.CategoryID, nil)
	if err != nil {
		return nil, types.NewServerError("Error in fetching category by id", operationName, err)
	}
	if category == nil {
		return nil, courseError.Course_NotFoundCategory
	}
	teacherAuth := authContext.(*types.TokenClaim)
	teacher, err := svc.unitOfWork.UserRepo().GetByID(teacherAuth.UserID, nil)
	if err != nil {
		return nil, types.NewServerError("Error in fetching teacher by id", operationName, err)
	}
	if teacher == nil {
		return nil, courseError.Course_NotFoundTeacher
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
	if err := svc.unitOfWork.CourseRepo().Create(course); err != nil {
		return nil, types.NewServerError("Error in creating teacher course", "teacherCourseService.Create", err)
	}
	// TODO: notification system
	//send email for creating course
	//create notification about creating new course for all admin
	return course, nil
}

func (svc teacherCourseService) FetchByTeacherId(authContext any, page, pageSize int) ([]*entities.Course, int, error) {
	const operationName = "teacherCourseService.FetchByTeacherId"
	authClaim := authContext.(*types.TokenClaim)
	teacher, err := svc.unitOfWork.UserRepo().GetByID(authClaim.UserID, nil)
	if err != nil {
		return nil, 0, types.NewServerError("Error in fetching teacher by id", operationName, err)
	}
	if teacher == nil {
		return nil, 0, userError.User_TeacherNotFound
	}
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
