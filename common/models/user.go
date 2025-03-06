package models

import (
	"time"
)

// swagger:model User
type User struct {
	// ID пользователя
	// example: 1
	ID uint `json:"-" gorm:"primarykey"`

	// Логин пользователя
	// example: johndoe
	Login string `json:"login" gorm:"unique;not null"`

	// Email пользователя
	// example: john.doe@example.com
	Email string `json:"email" gorm:"unique;not null"`

	// Хеш пароля (не отправляется в API ответах)
	PasswordHash string `json:"-" gorm:"not null"`

	// Имя пользователя
	// example: John
	Name string `json:"name,omitempty"`

	// Фамилия пользователя
	// example: Doe
	Surname string `json:"surname,omitempty"`

	// Дата рождения пользователя
	// example: 1990-01-01T00:00:00Z
	BirthDate *time.Time `json:"birthDate,omitempty"`

	// Номер телефона пользователя
	// example: +12345678901
	PhoneNumber string `json:"phoneNumber,omitempty"`

	// Дата создания пользователя
	// example: 2023-01-01T12:00:00Z
	CreatedAt time.Time `json:"createdAt"`

	// Дата последнего обновления пользователя
	// example: 2023-01-02T12:00:00Z
	UpdatedAt time.Time `json:"updatedAt"`
}
