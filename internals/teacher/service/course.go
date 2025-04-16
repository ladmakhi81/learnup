package service

import (
	"github.com/ladmakhi81/learnup/db/entities"
	categoryService "github.com/ladmakhi81/learnup/internals/category/service"
	courseRepository "github.com/ladmakhi81/learnup/internals/course/repo"
	courseService "github.com/ladmakhi81/learnup/internals/course/service"
	teacherDtoReq "github.com/ladmakhi81/learnup/internals/teacher/dto/req"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
)

type TeacherCourseService interface {
	Create(authContext any, dto teacherDtoReq.CreateCourseReq) (*entities.Course, error)
	FetchByTeacherId(authContext any, page, pageSize int) ([]*entities.Course, error)
	FetchCountByTeacherId(authContext any) (int, error)
}

type TeacherCourseServiceImpl struct {
	courseSvc      courseService.CourseService
	categorySvc    categoryService.CategoryService
	userSvc        userService.UserSvc
	courseRepo     courseRepository.CourseRepo
	translationSvc contracts.Translator
}

func NewTeacherCourseServiceImpl(
	courseSvc courseService.CourseService,
	categorySvc categoryService.CategoryService,
	userSvc userService.UserSvc,
	courseRepo courseRepository.CourseRepo,
	translationSvc contracts.Translator,
) *TeacherCourseServiceImpl {
	return &TeacherCourseServiceImpl{
		courseSvc:      courseSvc,
		categorySvc:    categorySvc,
		userSvc:        userSvc,
		courseRepo:     courseRepo,
		translationSvc: translationSvc,
	}
}

func (svc TeacherCourseServiceImpl) Create(authContext any, dto teacherDtoReq.CreateCourseReq) (*entities.Course, error) {
	isDuplicate, duplicateErr := svc.courseSvc.IsCourseNameExist(dto.Name)
	if duplicateErr != nil {
		return nil, duplicateErr
	}
	if isDuplicate {
		return nil, types.NewConflictError(
			svc.translationSvc.Translate("course.errors.name_duplicate"),
		)
	}
	category, categoryErr := svc.categorySvc.FindByID(dto.CategoryID)
	if categoryErr != nil {
		return nil, categoryErr
	}
	if category == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate("course.errors.not_found_category"),
		)
	}
	teacherAuth := authContext.(*types.TokenClaim)
	teacher, teacherErr := svc.userSvc.FindById(teacherAuth.UserID)
	if teacherErr != nil {
		return nil, teacherErr
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
	if err := svc.courseRepo.Create(course); err != nil {
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

func (svc TeacherCourseServiceImpl) FetchByTeacherId(authContext any, page, pageSize int) ([]*entities.Course, error) {
	authClaim := authContext.(*types.TokenClaim)
	teacher, teacherErr := svc.userSvc.FindById(authClaim.UserID)
	if teacherErr != nil {
		return nil, teacherErr
	}
	if teacher == nil {
		return nil, types.NewNotFoundError(
			svc.translationSvc.Translate("user.errors.teacher_not_found"),
		)
	}
	courses, coursesErr := svc.courseRepo.FetchPage(courseRepository.FetchPageOption{
		Page:      &page,
		PageSize:  &pageSize,
		TeacherId: &teacher.ID,
	})
	if coursesErr != nil {
		return nil, types.NewServerError(
			"Error in fetching courses related to teacher",
			"TeacherCourseServiceImpl.FetchByTeacherId",
			coursesErr,
		)
	}
	return courses, nil
}

func (svc TeacherCourseServiceImpl) FetchCountByTeacherId(authContext any) (int, error) {
	authClaim := authContext.(*types.TokenClaim)
	teacher, teacherErr := svc.userSvc.FindById(authClaim.UserID)
	if teacherErr != nil {
		return 0, teacherErr
	}
	if teacher == nil {
		return 0, types.NewNotFoundError(
			svc.translationSvc.Translate("user.errors.teacher_not_found"),
		)
	}
	count, countErr := svc.courseRepo.FetchCount(courseRepository.FetchCountOption{
		TeacherId: &teacher.ID,
	})
	if countErr != nil {
		return 0, types.NewServerError(
			"Error in return count of course by teacher id",
			"FetchCountByTeacherId",
			countErr,
		)
	}
	return count, nil
}
