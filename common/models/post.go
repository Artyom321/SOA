package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// StringArray is a custom type that implements the necessary interfaces for PostgreSQL arrays
type StringArray []string

// Value implements driver.Valuer interface
func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	// Convert []string to PostgreSQL array format
	var quotedStrings []string
	for _, s := range a {
		quotedStrings = append(quotedStrings, fmt.Sprintf("%q", s))
	}

	return fmt.Sprintf("{%s}", strings.Join(quotedStrings, ",")), nil
}

// Scan implements sql.Scanner interface
func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}

	// Handle string representation from PostgreSQL
	strValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("failed to scan array value: %v", value)
	}

	// Remove the curly braces and split
	strValue = strings.Trim(strValue, "{}")
	if strValue == "" {
		*a = []string{}
		return nil
	}

	// PostgreSQL returns quoted values; we need to parse them
	// This is a simplified version; a full implementation would handle escapes better
	result := []string{}
	for _, item := range strings.Split(strValue, ",") {
		// Remove the quotes around each item
		item = strings.Trim(item, "\"")
		result = append(result, item)
	}

	*a = result
	return nil
}

// RegisterArrayType doesn't need to do anything since the type is already defined
func RegisterArrayType() {
	// This is a no-op function since the StringArray type is already defined in this package
	// We keep it to maintain the API
}

// swagger:model Post
type Post struct {
	// ID поста
	// example: 123
	ID uint `json:"id" gorm:"primaryKey;autoIncrement"`

	// Название поста
	// example: Мой первый пост
	Title string `json:"title" binding:"required"`

	// Описание поста
	// example: Это описание моего первого поста
	Description string `json:"description"`

	// ID создателя поста
	// example: 42
	CreatorID uint `json:"creator_id"`

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
	ID uint `json:"id"`

	// Название поста
	// example: Мой первый пост
	Title string `json:"title"`

	// Описание поста
	// example: Это описание моего первого поста
	Description string `json:"description"`

	// ID создателя поста
	// example: 42
	CreatorID uint `json:"creator_id"`

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
	TotalCount int `json:"total_count"`
}

// swagger:model DeletePostResponse
type DeletePostResponse struct {
	// Результат удаления
	// example: true
	Success bool `json:"success"`
}
