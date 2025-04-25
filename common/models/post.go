package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	var quotedStrings []string
	for _, s := range a {
		quotedStrings = append(quotedStrings, fmt.Sprintf("%q", s))
	}

	return fmt.Sprintf("{%s}", strings.Join(quotedStrings, ",")), nil
}

func (a *StringArray) Scan(value any) error {
	if value == nil {
		*a = nil
		return nil
	}

	strValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("failed to scan array value: %v", value)
	}

	strValue = strings.Trim(strValue, "{}")
	if strValue == "" {
		*a = []string{}
		return nil
	}

	result := []string{}
	for _, item := range strings.Split(strValue, ",") {
		item = strings.Trim(item, "\"")
		result = append(result, item)
	}

	*a = result
	return nil
}

// swagger:model Post
type Post struct {
	// ID поста
	// example: 123
	ID uint64 `json:"id" gorm:"primaryKey;autoIncrement"`

	// Название поста
	// example: Мой первый пост
	Title string `json:"title" binding:"required"`

	// Описание поста
	// example: Это описание моего первого поста
	Description string `json:"description"`

	// ID создателя поста
	// example: 42
	CreatorID uint64 `json:"creator_id"`

	// Дата создания поста
	// example: 2023-01-01T12:00:00Z
	CreatedAt time.Time `json:"created_at"`

	// Дата обновления поста
	// example: 2023-01-02T12:00:00Z
	UpdatedAt time.Time `json:"updated_at"`

	// Флаг приватности (если true, то доступ на просмотр только у создателя)
	// example: false
	IsPrivate bool `json:"is_private"`

	// Список тегов поста
	// example: ["тег1", "тег2"]
	Tags StringArray `json:"tags" gorm:"type:text[]"`
}

// swagger:model CreatePostRequest
type CreatePostRequest struct {
	// Название поста
	// example: Мой первый пост
	// required: true
	Title string `json:"title" binding:"required"`

	// Описание поста
	// example: Это описание моего первого поста
	Description string `json:"description"`

	// Флаг приватности
	// example: false
	IsPrivate bool `json:"is_private"`

	// Список тегов
	// example: ["тег1", "тег2"]
	Tags []string `json:"tags"`
}

// swagger:model UpdatePostRequest
type UpdatePostRequest struct {
	// Название поста
	// example: Обновленное название поста
	Title string `json:"title"`

	// Описание поста
	// example: Обновленное описание поста
	Description string `json:"description"`

	// Флаг приватности
	// example: true
	IsPrivate bool `json:"is_private"`

	// Список тегов
	// example: ["тег1", "тег2", "новыйТег"]
	Tags []string `json:"tags"`
}

// swagger:model PostResponse
type PostResponse struct {
	// ID поста
	// example: 123
	ID uint64 `json:"id"`

	// Название поста
	// example: Мой первый пост
	Title string `json:"title"`

	// Описание поста
	// example: Это описание моего первого поста
	Description string `json:"description"`

	// ID создателя поста
	// example: 42
	CreatorID uint64 `json:"creator_id"`

	// Дата создания поста
	// example: 2023-01-01T12:00:00Z
	CreatedAt time.Time `json:"created_at"`

	// Дата обновления поста
	// example: 2023-01-02T12:00:00Z
	UpdatedAt time.Time `json:"updated_at"`

	// Флаг приватности
	// example: false
	IsPrivate bool `json:"is_private"`

	// Список тегов
	// example: ["тег1", "тег2"]
	Tags []string `json:"tags"`
}

// swagger:model PostListResponse
type PostListResponse struct {
	// Список постов
	Posts []PostResponse `json:"posts"`

	// Общее количество постов
	// example: 42
	TotalCount uint64 `json:"total_count"`
}

// swagger:model DeletePostResponse
type DeletePostResponse struct {
	// Результат удаления
	// example: true
	Success bool `json:"success"`
}
