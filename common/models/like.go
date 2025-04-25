package models

import "time"

// swagger:model Like
type Like struct {
	// ID лайка
	// example: 789
	ID uint64 `json:"id" gorm:"primaryKey;autoIncrement"`

	// ID поста, к которому относится лайк
	// example: 123
	PostID uint64 `json:"post_id" gorm:"index"`

	// ID пользователя, поставившего лайк
	// example: 42
	UserID uint64 `json:"user_id" gorm:"index"`

	// Дата создания лайка
	// example: 2023-01-15T14:32:10Z
	CreatedAt time.Time `json:"created_at"`
}

// swagger:model LikePostResponse
type LikePostResponse struct {
	// Успешность операции лайка
	// example: true
	Success bool `json:"success"`

	// Общее количество лайков у поста
	// example: 42
	TotalLikes int `json:"total_likes"`
}

// swagger:model SuccessResponse
type SuccessResponse struct {
	// Успешность выполнения операции
	// example: true
	Success bool `json:"success"`
}
