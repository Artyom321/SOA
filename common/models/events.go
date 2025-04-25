package models

import "time"

// swagger:model UserRegistrationEvent
type UserRegistrationEvent struct {
	// ID зарегистрированного пользователя
	// example: 42
	UserID uint64 `json:"user_id"`

	// Логин пользователя
	// example: johndoe
	Login string `json:"login"`

	// Email пользователя
	// example: john.doe@example.com
	Email string `json:"email"`

	// Дата регистрации пользователя
	// example: 2023-01-15T10:00:00Z
	CreatedAt time.Time `json:"created_at"`
}

// swagger:model PostViewEvent
type PostViewEvent struct {
	// ID пользователя, просмотревшего пост
	// example: 42
	UserID uint64 `json:"user_id"`

	// ID просмотренного поста
	// example: 123
	PostID uint64 `json:"post_id"`

	// Время просмотра поста
	// example: 2023-01-15T14:35:22Z
	ViewedAt time.Time `json:"viewed_at"`
}

// swagger:model PostLikeEvent
type PostLikeEvent struct {
	// ID пользователя, поставившего лайк
	// example: 42
	UserID uint64 `json:"user_id"`

	// ID поста, который получил лайк
	// example: 123
	PostID uint64 `json:"post_id"`

	// Время установки лайка
	// example: 2023-01-15T14:37:45Z
	LikedAt time.Time `json:"liked_at"`
}

// swagger:model PostCommentEvent
type PostCommentEvent struct {
	// ID пользователя, оставившего комментарий
	// example: 42
	UserID uint64 `json:"user_id"`

	// ID поста, к которому оставлен комментарий
	// example: 123
	PostID uint64 `json:"post_id"`

	// ID созданного комментария
	// example: 456
	CommentID uint64 `json:"comment_id"`

	// Текст комментария
	// example: Очень информативный пост!
	CommentText string `json:"comment_text"`

	// Время создания комментария
	// example: 2023-01-15T14:40:12Z
	CreatedAt time.Time `json:"created_at"`
}
