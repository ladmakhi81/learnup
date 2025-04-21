package db

import (
	"github.com/ladmakhi81/learnup/db/repositories"
	"gorm.io/gorm"
)

type Repositories struct {
	AnswerRepo       repositories.AnswerRepo
	CartRepo         repositories.CartRepo
	CategoryRepo     repositories.CategoryRepo
	CommentRepo      repositories.CommentRepo
	CourseRepo       repositories.CourseRepo
	LikeRepo         repositories.LikeRepo
	NotificationRepo repositories.NotificationRepo
	OrderRepo        repositories.OrderRepo
	OrderItemRepo    repositories.OrderItemRepo
	PaymentRepo      repositories.PaymentRepo
	QuestionRepo     repositories.QuestionRepo
	TransactionRepo  repositories.TransactionRepo
	UserRepo         repositories.UserRepo
	VideoRepo        repositories.VideoRepo
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		AnswerRepo:       repositories.NewAnswerRepo(db),
		CartRepo:         repositories.NewCartRepo(db),
		CategoryRepo:     repositories.NewCategoryRepo(db),
		CommentRepo:      repositories.NewCommentRepo(db),
		CourseRepo:       repositories.NewCourseRepo(db),
		OrderRepo:        repositories.NewOrderRepo(db),
		OrderItemRepo:    repositories.NewOrderItemRepo(db),
		PaymentRepo:      repositories.NewPaymentRepo(db),
		QuestionRepo:     repositories.NewQuestionRepo(db),
		TransactionRepo:  repositories.NewTransactionRepo(db),
		LikeRepo:         repositories.NewLikeRepo(db),
		NotificationRepo: repositories.NewNotificationRepo(db),
		UserRepo:         repositories.NewUserRepo(db),
		VideoRepo:        repositories.NewVideoRepo(db),
	}
}
