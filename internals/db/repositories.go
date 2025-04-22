package db

import (
	"github.com/ladmakhi81/learnup/internals/db/repositories"
	"gorm.io/gorm"
)

type Repo interface {
	AnswerRepo() repositories.AnswerRepo
	CartRepo() repositories.CartRepo
	CategoryRepo() repositories.CategoryRepo
	CommentRepo() repositories.CommentRepo
	CourseRepo() repositories.CourseRepo
	LikeRepo() repositories.LikeRepo
	NotificationRepo() repositories.NotificationRepo
	OrderRepo() repositories.OrderRepo
	OrderItemRepo() repositories.OrderItemRepo
	PaymentRepo() repositories.PaymentRepo
	QuestionRepo() repositories.QuestionRepo
	TransactionRepo() repositories.TransactionRepo
	UserRepo() repositories.UserRepo
	VideoRepo() repositories.VideoRepo
}

type RepoProvider struct {
	answerRepo       repositories.AnswerRepo
	cartRepo         repositories.CartRepo
	categoryRepo     repositories.CategoryRepo
	commentRepo      repositories.CommentRepo
	courseRepo       repositories.CourseRepo
	likeRepo         repositories.LikeRepo
	notificationRepo repositories.NotificationRepo
	orderRepo        repositories.OrderRepo
	orderItemRepo    repositories.OrderItemRepo
	paymentRepo      repositories.PaymentRepo
	questionRepo     repositories.QuestionRepo
	transactionRepo  repositories.TransactionRepo
	userRepo         repositories.UserRepo
	videoRepo        repositories.VideoRepo
}

func NewRepoProvider(tx *gorm.DB) *RepoProvider {
	return &RepoProvider{
		answerRepo:       repositories.NewAnswerRepo(tx),
		cartRepo:         repositories.NewCartRepo(tx),
		categoryRepo:     repositories.NewCategoryRepo(tx),
		commentRepo:      repositories.NewCommentRepo(tx),
		courseRepo:       repositories.NewCourseRepo(tx),
		likeRepo:         repositories.NewLikeRepo(tx),
		notificationRepo: repositories.NewNotificationRepo(tx),
		orderRepo:        repositories.NewOrderRepo(tx),
		orderItemRepo:    repositories.NewOrderItemRepo(tx),
		paymentRepo:      repositories.NewPaymentRepo(tx),
		questionRepo:     repositories.NewQuestionRepo(tx),
		transactionRepo:  repositories.NewTransactionRepo(tx),
		userRepo:         repositories.NewUserRepo(tx),
		videoRepo:        repositories.NewVideoRepo(tx),
	}
}

func (svc RepoProvider) AnswerRepo() repositories.AnswerRepo {
	return svc.answerRepo
}
func (svc RepoProvider) CartRepo() repositories.CartRepo {
	return svc.cartRepo
}
func (svc RepoProvider) CategoryRepo() repositories.CategoryRepo {
	return svc.categoryRepo
}
func (svc RepoProvider) CommentRepo() repositories.CommentRepo {
	return svc.commentRepo
}
func (svc RepoProvider) CourseRepo() repositories.CourseRepo {
	return svc.courseRepo
}
func (svc RepoProvider) LikeRepo() repositories.LikeRepo {
	return svc.likeRepo
}
func (svc RepoProvider) NotificationRepo() repositories.NotificationRepo {
	return svc.notificationRepo
}
func (svc RepoProvider) OrderRepo() repositories.OrderRepo {
	return svc.orderRepo
}
func (svc RepoProvider) OrderItemRepo() repositories.OrderItemRepo {
	return svc.orderItemRepo
}
func (svc RepoProvider) PaymentRepo() repositories.PaymentRepo {
	return svc.paymentRepo
}
func (svc RepoProvider) QuestionRepo() repositories.QuestionRepo {
	return svc.questionRepo
}
func (svc RepoProvider) TransactionRepo() repositories.TransactionRepo {
	return svc.transactionRepo
}
func (svc RepoProvider) UserRepo() repositories.UserRepo {
	return svc.userRepo
}
func (svc RepoProvider) VideoRepo() repositories.VideoRepo {
	return svc.videoRepo
}
