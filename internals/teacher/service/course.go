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
	tx, txErr := svc.unitOfWork.Begin()
	if txErr != nil {
		return nil, txErr
	}
	isDuplicate, duplicateErr := tx.CourseRepo().Exist(map[string]any{
		"name": dto.Name,
	})
	if duplicateErr != nil {
		return nil, types.NewServerError(
			"Error in checking existence of course name",
			"TeacherCourseServiceImpl.Create",
			duplicateErr,
		)
	}
	if isDuplicate {
		return nil, types.NewConflictError(
			svc.translationSvc.Translate("course.errors.name_duplicate"),
		)
	}
	category, categoryErr := tx.CategoryRepo().GetByID(dto.CategoryID, nil)
	if categoryErr != nil {
		return nil, types.NewServerError(
			"Error in fetching category by id",
			"TeacherCourseServiceImpl.Create",
			categoryErr,
		)
	}
	if category == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate("course.errors.not_found_category"),
		)
	}
	teacherAuth := authContext.(*types.TokenClaim)
	teacher, teacherErr := tx.UserRepo().GetByID(teacherAuth.UserID, nil)
	if teacherErr != nil {
		return nil, types.NewServerError(
			"Error in fetching teacher by id",
			"TeacherCourseServiceImpl.Create",
			teacherErr,
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
	if err := tx.CourseRepo().Create(course); err != nil {
		return nil, types.NewServerError(
			"Error in creating teacher course",
			"TeacherCourseServiceImpl.Create",
			err,
		)
	}
	// TODO: notification system
	//send email for creating course
	//create notification about creating new course for all admin
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return course, nil
}

func (svc TeacherCourseServiceImpl) FetchByTeacherId(authContext any, page, pageSize int) ([]*entities.Course, int, error) {
	tx, txErr := svc.unitOfWork.Begin()
	if txErr != nil {
		return nil, 0, txErr
	}
	authClaim := authContext.(*types.TokenClaim)
	teacher, teacherErr := tx.UserRepo().GetByID(authClaim.UserID, nil)
	if teacherErr != nil {
		return nil, 0, types.NewServerError(
			"Error in fetching teacher by id",
			"TeacherCourseServiceImpl.FetchByTeacherId",
			teacherErr,
		)
	}
	if teacher == nil {
		return nil, 0, types.NewNotFoundError(
			svc.translationSvc.Translate("user.errors.teacher_not_found"),
		)
	}
	courses, count, coursesErr := tx.CourseRepo().GetPaginated(repositories.GetPaginatedOptions{
		Offset: &page,
		Limit:  &pageSize,
		Conditions: map[string]any{
			"teacher_id": teacher.ID,
		},
	})
	if coursesErr != nil {
		return nil, 0, types.NewServerError(
			"Error in fetching courses related to teacher",
			"TeacherCourseServiceImpl.FetchByTeacherId",
			coursesErr,
		)
	}
	if err := tx.Commit(); err != nil {
		return nil, 0, err
	}
	return courses, count, nil
}
