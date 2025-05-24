package models

import "time"

// swagger:model PostStatsResponse
type PostStatsResponse struct {
	// ID поста
	// example: 123
	PostID uint64 `json:"post_id"`

	// Количество просмотров
	// example: 1500
	ViewsCount uint64 `json:"views_count"`

	// Количество лайков
	// example: 320
	LikesCount uint64 `json:"likes_count"`

	// Количество комментариев
	// example: 45
	CommentsCount uint64 `json:"comments_count"`
}

// swagger:model TimelineItem
type TimelineItem struct {
	// Дата (день)
	// example: 2023-09-01T00:00:00Z
	Date time.Time `json:"date"`

	// Количество событий
	// example: 42
	Count uint64 `json:"count"`
}

// swagger:model TimelineResponse
type TimelineResponse struct {
	// Временная динамика по дням
	Items []TimelineItem `json:"items"`
}

// swagger:model MetricType
type MetricType string

const (
	// MetricTypeViews для просмотров
	MetricTypeViews MetricType = "views"
	// MetricTypeLikes для лайков
	MetricTypeLikes MetricType = "likes"
	// MetricTypeComments для комментариев
	MetricTypeComments MetricType = "comments"
)

// swagger:model TopPostItem
type TopPostItem struct {
	// ID поста
	// example: 123
	PostID uint64 `json:"post_id"`

	// Количество (просмотров/лайков/комментариев)
	// example: 1500
	Count uint64 `json:"count"`
}

// swagger:model TopPostsResponse
type TopPostsResponse struct {
	// Список топ постов
	Posts []TopPostItem `json:"posts"`
}

// swagger:model TopUserItem
type TopUserItem struct {
	// ID пользователя
	// example: 42
	UserID uint64 `json:"user_id"`

	// Количество (просмотров/лайков/комментариев)
	// example: 1500
	Count uint64 `json:"count"`
}

// swagger:model TopUsersResponse
type TopUsersResponse struct {
	// Список топ пользователей
	Users []TopUserItem `json:"users"`
}
