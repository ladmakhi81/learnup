package db

import (
	repositories2 "github.com/ladmakhi81/learnup/internals/db/repositories"
	"gorm.io/gorm"
)

type Repositories struct {
	AnswerRepo       repositories2.AnswerRepo
	CartRepo         repositories2.CartRepo
	CategoryRepo     repositories2.CategoryRepo
	CommentRepo      repositories2.CommentRepo
	CourseRepo       repositories2.CourseRepo
	LikeRepo         repositories2.LikeRepo
	NotificationRepo repositories2.NotificationRepo
	OrderRepo        repositories2.OrderRepo
	OrderItemRepo    repositories2.OrderItemRepo
	PaymentRepo      repositories2.PaymentRepo
	QuestionRepo     repositories2.QuestionRepo
	TransactionRepo  repositories2.TransactionRepo
	UserRepo         repositories2.UserRepo
	VideoRepo        repositories2.VideoRepo
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		AnswerRepo:       repositories2.NewAnswerRepo(db),
		CartRepo:         repositories2.NewCartRepo(db),
		CategoryRepo:     repositories2.NewCategoryRepo(db),
		CommentRepo:      repositories2.NewCommentRepo(db),
		CourseRepo:       repositories2.NewCourseRepo(db),
		OrderRepo:        repositories2.NewOrderRepo(db),
		OrderItemRepo:    repositories2.NewOrderItemRepo(db),
		PaymentRepo:      repositories2.NewPaymentRepo(db),
		QuestionRepo:     repositories2.NewQuestionRepo(db),
		TransactionRepo:  repositories2.NewTransactionRepo(db),
		LikeRepo:         repositories2.NewLikeRepo(db),
		NotificationRepo: repositories2.NewNotificationRepo(db),
		UserRepo:         repositories2.NewUserRepo(db),
		VideoRepo:        repositories2.NewVideoRepo(db),
	}
}
