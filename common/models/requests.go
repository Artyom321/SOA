package models

import (
	"time"
)

// swagger:model RegisterRequest
type RegisterRequest struct {
	// Логин пользователя (минимум 3 символа)
	// required: true
	// min length: 3
	// max length: 50
	// example: johndoe
	Login string `json:"login" binding:"required,min=3,max=50"`

	// Email пользователя
	// required: true
	// example: john.doe@example.com
	Email string `json:"email" binding:"required,email"`

	// Пароль пользователя (минимум 6 символов)
	// required: true
	// min length: 6
	// example: securePassword123
	Password string `json:"password" binding:"required,min=6,max=100"`
}

// swagger:model LoginRequest
type LoginRequest struct {
	// Логин пользователя
	// required: true
	// example: johndoe
	Login string `json:"login" binding:"required"`

	// Пароль пользователя
	// required: true
	// example: securePassword123
	Password string `json:"password" binding:"required"`
}

// swagger:model ProfileUpdateRequest
type ProfileUpdateRequest struct {
	// Имя пользователя
	// example: John
	Name *string `json:"name" binding:"omitempty,min=2,max=50"`

	// Фамилия пользователя
	// example: Doe
	Surname *string `json:"surname" binding:"omitempty,min=2,max=50"`

	// Email пользователя
	// example: john.doe@example.com
	Email *string `json:"email" binding:"omitempty,email"`

	// Номер телефона пользователя
	// example: +12345678901
	PhoneNumber *string `json:"phoneNumber,omitempty"`

	// Дата рождения пользователя
	// example: 1990-01-01T00:00:00Z
	BirthDate *time.Time `json:"birthDate,omitempty"`
}
