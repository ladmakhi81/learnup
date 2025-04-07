package service

import (
	categoryService "github.com/ladmakhi81/learnup/internals/category/service"
	dtoreq "github.com/ladmakhi81/learnup/internals/course/dto/req"
	"github.com/ladmakhi81/learnup/internals/course/entity"
	"github.com/ladmakhi81/learnup/internals/course/repo"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	"github.com/ladmakhi81/learnup/pkg/translations"
	"github.com/ladmakhi81/learnup/types"
	"time"
)

type CourseService interface {
	Create(authContext any, dto dtoreq.CreateCourseReq) (*entity.Course, error)
	FindByName(name string) (*entity.Course, error)
	IsCourseNameExist(name string) (bool, error)
	GetCourses(page, pageSize int) ([]*entity.Course, error)
	GetCoursesCount() (int, error)
}

type CourseServiceImpl struct {
	courseRepo   repo.CourseRepo
	categorySvc  categoryService.CategoryService
	userSvc      userService.UserSvc
	translateSvc translations.Translator
}

func NewCourseServiceImpl(
	courseRepo repo.CourseRepo,
	translateSvc translations.Translator,
	userSvc userService.UserSvc,
	categorySvc categoryService.CategoryService,
) *CourseServiceImpl {
	return &CourseServiceImpl{
		courseRepo:   courseRepo,
		translateSvc: translateSvc,
		userSvc:      userSvc,
		categorySvc:  categorySvc,
	}
}

func (svc CourseServiceImpl) Create(authContext any, dto dtoreq.CreateCourseReq) (*entity.Course, error) {
	isCourseNameDuplicated, courseNameDuplicatedErr := svc.IsCourseNameExist(dto.Name)
	if courseNameDuplicatedErr != nil {
		return nil, courseNameDuplicatedErr
	}
	if isCourseNameDuplicated {
		return nil, types.NewConflictError(
			svc.translateSvc.Translate("course.errors.name_duplicate"),
		)
	}
	category, categoryErr := svc.categorySvc.FindByID(dto.CategoryID)
	if categoryErr != nil {
		return nil, categoryErr
	}
	if category == nil {
		return nil, types.NewNotFoundError("course.errors.not_found_category")
	}
	teacher, teacherErr := svc.userSvc.FindById(dto.TeacherID)
	if teacherErr != nil {
		return nil, teacherErr
	}
	if teacher == nil {
		return nil, types.NewNotFoundError("course.errors.not_found_teacher")
	}
	authClaim := authContext.(*types.TokenClaim)
	authUser, authUserErr := svc.userSvc.FindById(authClaim.UserID)
	if authUserErr != nil {
		return nil, authUserErr
	}

	verifiedDate := time.Now()
	course := &entity.Course{
		CategoryID:                  &category.ID,
		Name:                        dto.Name,
		AbilityToAddComment:         dto.AbilityToAddComment,
		CanHaveDiscount:             dto.CanHaveDiscount,
		CommentAccessMode:           dto.CommentAccessMode,
		Description:                 dto.Description,
		DiscountFeeAmountPercentage: dto.DiscountFeeAmountPercentage,
		Fee:                         dto.Fee,
		IsPublished:                 true,
		Image:                       dto.Image,
		IntroductionVideo:           dto.IntroductionVideo,
		IsVerifiedByAdmin:           true,
		Level:                       dto.Level,
		MaxDiscountAmount:           dto.MaxDiscountAmount,
		Prerequisite:                dto.Prerequisite,
		Price:                       dto.Price,
		Status:                      entity.CourseStatus_Starting,
		Tags:                        dto.Tags,
		TeacherID:                   &teacher.ID,
		ThumbnailImage:              dto.ThumbnailImage,
		VerifiedDate:                &verifiedDate,
		VerifiedByID:                &authUser.ID,
	}
	if err := svc.courseRepo.Create(course); err != nil {
		return nil, types.NewServerError(
			"Create Course Throw Error",
			"CourseServiceImpl.Create",
			err,
		)
	}
	return course, nil
}

func (svc CourseServiceImpl) FindByName(name string) (*entity.Course, error) {
	course, courseErr := svc.courseRepo.FindByName(name)
	if courseErr != nil {
		return nil, types.NewServerError(
			"Error in finding course by name throw error",
			"CourseServiceImpl.FindByName",
			courseErr,
		)
	}
	if course == nil {
		return nil, nil
	}
	return course, nil
}

func (svc CourseServiceImpl) IsCourseNameExist(name string) (bool, error) {
	course, courseErr := svc.FindByName(name)
	if courseErr != nil {
		return false, courseErr
	}
	if course == nil {
		return false, nil
	}
	return true, nil
}

func (svc CourseServiceImpl) GetCourses(page, pageSize int) ([]*entity.Course, error) {
	courses, coursesErr := svc.courseRepo.GetCourses(page, pageSize)
	if coursesErr != nil {
		return nil, types.NewServerError(
			"Find All Pageable Courses Throw Error",
			"CourseServiceImpl.GetCourses",
			coursesErr,
		)
	}
	return courses, nil
}

func (svc CourseServiceImpl) GetCoursesCount() (int, error) {
	count, countErr := svc.courseRepo.GetCoursesCount()
	if countErr != nil {
		return 0, types.NewServerError(
			"Get Count Of Courses Throw Error",
			"CourseServiceImpl.GetCourses",
			countErr,
		)
	}
	return count, nil
}
