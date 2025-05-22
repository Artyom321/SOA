package models

import "time"

// swagger:model Comment
type Comment struct {
	// ID комментария
	// example: 456
	ID uint64 `json:"id" gorm:"primaryKey;autoIncrement"`

	// ID поста, к которому относится комментарий
	// example: 123
	PostID uint64 `json:"post_id" gorm:"index"`

	// ID пользователя, оставившего комментарий
	// example: 42
	UserID uint64 `json:"user_id"`

	// Текст комментария
	// example: Это очень интересный пост!
	Content string `json:"content"`

	// Дата создания комментария
	// example: 2023-01-15T14:30:45Z
	CreatedAt time.Time `json:"created_at"`
}

// swagger:model AddCommentRequest
type AddCommentRequest struct {
	// Текст комментария
	// required: true
	// example: Отличный пост, спасибо за информацию!
	Content string `json:"content" binding:"required"`
}

// swagger:model CommentResponse
type CommentResponse struct {
	// Данные комментария
	Comment Comment `json:"comment"`
}

// swagger:model CommentListResponse
type CommentListResponse struct {
	// Список комментариев
	Comments []Comment `json:"comments"`

	// Общее количество комментариев
	// example: 25
	TotalCount uint64 `json:"total_count"`
}
