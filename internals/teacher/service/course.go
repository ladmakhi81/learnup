package service

import (
	"github.com/ladmakhi81/learnup/internals/db"
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"github.com/ladmakhi81/learnup/internals/db/repositories"
	teacherDtoReq "github.com/ladmakhi81/learnup/internals/teacher/dto/req"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
)

type TeacherCourseService interface {
	Create(authContext any, dto teacherDtoReq.CreateCourseReq) (*entities.Course, error)
	FetchByTeacherId(authContext any, page, pageSize int) ([]*entities.Course, int, error)
}

type TeacherCourseServiceImpl struct {
	unitOfWork     db.UnitOfWork
	translationSvc contracts.Translator
}

func NewTeacherCourseServiceImpl(
	unitOfWork db.UnitOfWork,
	translationSvc contracts.Translator,
) *TeacherCourseServiceImpl {
	return &TeacherCourseServiceImpl{
		unitOfWork:     unitOfWork,
		translationSvc: translationSvc,
	}
}

func (svc TeacherCourseServiceImpl) Create(authContext any, dto teacherDtoReq.CreateCourseReq) (*entities.Course, error) {
	const operationName = "TeacherCourseServiceImpl.Create"
	isDuplicate, err := svc.unitOfWork.CourseRepo().Exist(map[string]any{
		"name": dto.Name,
	})
	if err != nil {
		return nil, types.NewServerError(
			"Error in checking existence of course name",
			operationName,
			err,
		)
	}
	if isDuplicate {
		return nil, types.NewConflictError(
			svc.translationSvc.Translate("course.errors.name_duplicate"),
		)
	}
	category, err := svc.unitOfWork.CategoryRepo().GetByID(dto.CategoryID, nil)
	if err != nil {
		return nil, types.NewServerError(
			"Error in fetching category by id",
			operationName,
			err,
		)
	}
	if category == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate("course.errors.not_found_category"),
		)
	}
	teacherAuth := authContext.(*types.TokenClaim)
	teacher, err := svc.unitOfWork.UserRepo().GetByID(teacherAuth.UserID, nil)
	if err != nil {
		return nil, types.NewServerError(
			"Error in fetching teacher by id",
			operationName,
			err,
		)
	}
	if teacher == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate("course.errors.not_found_teacher"),
		)
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
		return nil, types.NewServerError(
			"Error in creating teacher course",
			"TeacherCourseServiceImpl.Create",
			err,
		)
	}
	// TODO: notification system
	//send email for creating course
	//create notification about creating new course for all admin
	return course, nil
}

func (svc TeacherCourseServiceImpl) FetchByTeacherId(authContext any, page, pageSize int) ([]*entities.Course, int, error) {
	const operationName = "TeacherCourseServiceImpl.FetchByTeacherId"
	authClaim := authContext.(*types.TokenClaim)
	teacher, err := svc.unitOfWork.UserRepo().GetByID(authClaim.UserID, nil)
	if err != nil {
		return nil, 0, types.NewServerError(
			"Error in fetching teacher by id",
			operationName,
			err,
		)
	}
	if teacher == nil {
		return nil, 0, types.NewNotFoundError(
			svc.translationSvc.Translate("user.errors.teacher_not_found"),
		)
	}
	courses, count, err := svc.unitOfWork.CourseRepo().GetPaginated(repositories.GetPaginatedOptions{
		Offset: &page,
		Limit:  &pageSize,
		Conditions: map[string]any{
			"teacher_id": teacher.ID,
		},
	})
	if err != nil {
		return nil, 0, types.NewServerError(
			"Error in fetching courses related to teacher",
			operationName,
			err,
		)
	}
	return courses, count, nil
}
