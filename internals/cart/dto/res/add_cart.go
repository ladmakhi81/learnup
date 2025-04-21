package dtores

import "time"

type AddCartRes struct {
	UserID    uint      `json:"userId"`
	CourseID  uint      `json:"courseId"`
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}
