package service

import (
	"github.com/ladmakhi81/learnup/db/entities"
	categoryService "github.com/ladmakhi81/learnup/internals/category/service"
	dtoreq "github.com/ladmakhi81/learnup/internals/course/dto/req"
	"github.com/ladmakhi81/learnup/internals/course/repo"
	notificationReqDto "github.com/ladmakhi81/learnup/internals/notification/dto/req"
	notificationService "github.com/ladmakhi81/learnup/internals/notification/service"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
	"time"
)

type CourseService interface {
	Create(authContext any, dto dtoreq.CreateCourseReq) (*entities.Course, error)
	FindByName(name string) (*entities.Course, error)
	IsCourseNameExist(name string) (bool, error)
	GetCourses(page, pageSize int) ([]*entities.Course, error)
	GetCoursesCount() (int, error)
	FindById(id uint) (*entities.Course, error)
	FindDetailById(id uint) (*entities.Course, error)
	FindByVideoId(id uint) (*entities.Course, error)
	VerifyCourse(authContext any, dto dtoreq.VerifyCourseReq) error
	UpdateIntroductionURL(dto dtoreq.UpdateIntroductionURLReq) error
	CreateCompleteIntroductionVideoNotification(id uint) error
}

type CourseServiceImpl struct {
	courseRepo      repo.CourseRepo
	categorySvc     categoryService.CategoryService
	userSvc         userService.UserSvc
	translateSvc    contracts.Translator
	notificationSvc notificationService.NotificationService
}

func NewCourseServiceImpl(
	courseRepo repo.CourseRepo,
	translateSvc contracts.Translator,
	userSvc userService.UserSvc,
	categorySvc categoryService.CategoryService,
	notificationSvc notificationService.NotificationService,
) *CourseServiceImpl {
	return &CourseServiceImpl{
		courseRepo:      courseRepo,
		translateSvc:    translateSvc,
		userSvc:         userSvc,
		categorySvc:     categorySvc,
		notificationSvc: notificationSvc,
	}
}

func (svc CourseServiceImpl) Create(authContext any, dto dtoreq.CreateCourseReq) (*entities.Course, error) {
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
	course := &entities.Course{
		CategoryID:                  &category.ID,
		Name:                        dto.Name,
		AbilityToAddComment:         dto.AbilityToAddComment,
		CanHaveDiscount:             dto.CanHaveDiscount,
		CommentAccessMode:           dto.CommentAccessMode,
		Description:                 dto.Description,
		DiscountFeeAmountPercentage: dto.DiscountFeeAmountPercentage,
		IsPublished:                 true,
		Image:                       dto.Image,
		IntroductionVideo:           dto.IntroductionVideo,
		IsVerifiedByAdmin:           true,
		Level:                       dto.Level,
		MaxDiscountAmount:           dto.MaxDiscountAmount,
		Prerequisite:                dto.Prerequisite,
		Status:                      entities.CourseStatus_Starting,
		Tags:                        dto.Tags,
		TeacherID:                   &teacher.ID,
		ThumbnailImage:              dto.ThumbnailImage,
		VerifiedDate:                &verifiedDate,
		VerifiedByID:                &authUser.ID,
	}
	if dto.Price == nil {
		course.Price = 0
	} else {
		course.Price = *dto.Price
	}

	if dto.Fee == nil {
		course.Fee = 0
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

func (svc CourseServiceImpl) FindByName(name string) (*entities.Course, error) {
	course, courseErr := svc.courseRepo.FetchByName(name)
	if courseErr != nil {
		return nil, types.NewServerError(
			"Error in finding course by name throw error",
			"CourseServiceImpl.FetchByName",
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

func (svc CourseServiceImpl) GetCourses(page, pageSize int) ([]*entities.Course, error) {
	courses, coursesErr := svc.courseRepo.FetchPage(repo.FetchPageOption{
		PageSize: &pageSize,
		Page:     &page,
		Preloads: []string{
			"Teacher",
			"Category",
			"VerifiedBy",
		},
	})
	if coursesErr != nil {
		return nil, types.NewServerError(
			"Find All Pageable Courses Throw Error",
			"CourseServiceImpl.FetchPage",
			coursesErr,
		)
	}
	return courses, nil
}

func (svc CourseServiceImpl) GetCoursesCount() (int, error) {
	count, countErr := svc.courseRepo.FetchCount(repo.FetchCountOption{})
	if countErr != nil {
		return 0, types.NewServerError(
			"Get Count Of Courses Throw Error",
			"CourseServiceImpl.FetchPage",
			countErr,
		)
	}
	return count, nil
}

func (svc CourseServiceImpl) FindById(id uint) (*entities.Course, error) {
	course, courseErr := svc.courseRepo.FetchById(id)
	if courseErr != nil {
		return nil, types.NewServerError(
			"Find Course By ID Throw Error",
			"CourseServiceImpl.FetchById",
			courseErr,
		)
	}
	return course, nil
}

func (svc CourseServiceImpl) FindDetailById(id uint) (*entities.Course, error) {
	course, courseErr := svc.courseRepo.FetchDetailById(id)
	if courseErr != nil {
		return nil, types.NewServerError(
			"Find Course Detail By ID Throw Error",
			"CourseServiceImpl.FindDetailByID",
			courseErr,
		)
	}
	return course, nil
}

func (svc CourseServiceImpl) FindByVideoId(id uint) (*entities.Course, error) {
	course, courseErr := svc.courseRepo.FetchByVideoId(id)
	if courseErr != nil {
		return nil, types.NewServerError(
			"Error in happen in find course by video id",
			"CourseServiceImpl.FetchByVideoId",
			courseErr,
		)
	}
	return course, nil
}

func (svc CourseServiceImpl) VerifyCourse(authContext any, dto dtoreq.VerifyCourseReq) error {
	course, courseErr := svc.FindById(dto.ID)
	if courseErr != nil {
		return courseErr
	}
	if course == nil {
		return types.NewNotFoundError(
			svc.translateSvc.Translate(
				"course.errors.not_found",
			),
		)
	}
	if course.Status != entities.CourseStatus_InProgress {
		return types.NewBadRequestError(
			svc.translateSvc.Translate("course.errors.unable_to_verify"),
		)
	}
	adminClaim := authContext.(*types.TokenClaim)
	admin, adminErr := svc.userSvc.FindById(adminClaim.UserID)
	if adminErr != nil {
		return adminErr
	}
	if admin == nil {
		return types.NewNotFoundError(
			svc.translateSvc.Translate("user.errors.admin_not_found"),
		)
	}
	if dto.Fee > course.Price || dto.Fee < 0 || dto.Fee > course.Price-course.MaxDiscountAmount {
		return types.NewBadRequestError(
			svc.translateSvc.Translate("course.errors.invalid_fee"),
		)
	}
	if dto.DiscountFeeAmountPercentage > 100 {
		return types.NewBadRequestError(
			svc.translateSvc.Translate("course.errors.invalid_max_discount_percentage"),
		)
	}
	now := time.Now()
	course.Fee = dto.Fee
	if course.CanHaveDiscount {
		course.DiscountFeeAmountPercentage = dto.DiscountFeeAmountPercentage
	}
	course.Status = entities.CourseStatus_Verified
	course.StatusChangedAt = &now
	course.VerifiedByID = &admin.ID
	course.VerifiedDate = &now
	course.IsVerifiedByAdmin = true
	if err := svc.courseRepo.Update(course); err != nil {
		return types.NewServerError(
			"Error in verifying the course by admin",
			"CourseServiceImpl.VerifyCourse",
			err,
		)
	}
	//TODO: notification system
	// create notification for teacher that course verified
	// send email for this notification
	return nil
}

func (svc CourseServiceImpl) UpdateIntroductionURL(dto dtoreq.UpdateIntroductionURLReq) error {
	course, courseErr := svc.FindById(dto.CourseId)
	if courseErr != nil {
		return courseErr
	}
	if course == nil {
		return types.NewNotFoundError(
			svc.translateSvc.Translate("course.errors.not_found"),
		)
	}
	course.IntroductionVideo = dto.URL
	if err := svc.courseRepo.Update(course); err != nil {
		return types.NewServerError(
			"Error in setting introduction video url",
			"CourseServiceImpl.UpdateIntroductionURL",
			err,
		)
	}
	return nil
}

func (svc CourseServiceImpl) CreateCompleteIntroductionVideoNotification(id uint) error {
	course, courseErr := svc.FindById(id)
	if courseErr != nil {
		return courseErr
	}
	if course == nil {
		return types.NewNotFoundError(
			svc.translateSvc.Translate("course.errors.not_found"),
		)
	}
	notificationDto := notificationReqDto.CreateNotificationReq{
		UserID:    *course.TeacherID,
		EventType: entities.NotificationType_CompleteIntroductionCourseVideoUpload,
		Metadata: map[string]any{
			"courseId":   course.ID,
			"courseName": course.Name,
		},
	}
	_, notificationErr := svc.notificationSvc.Create(notificationDto)
	if notificationErr != nil {
		return notificationErr
	}
	return nil
}
