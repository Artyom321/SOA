package models

// swagger:model ErrorResponse
type ErrorResponse struct {
	// Сообщение об ошибке
	// example: Invalid credentials
	Error string `json:"error"`
}

// swagger:model RegisterResponse
type RegisterResponse struct {
	// Статус операции
	// example: user created
	Status string `json:"status"`
}

// swagger:model LoginResponse
type LoginResponse struct {
	// Сообщение о результате
	// example: Logged in successfully
	Message string `json:"message"`
}

// swagger:model ProfileResponse
type ProfileResponse struct {
	// Данные пользователя
	User User `json:"user"`
}

// swagger:model ProfileUpdateResponse
type ProfileUpdateResponse struct {
	// Статус операции
	// example: Profile updated successfully
	Status string `json:"status"`

	// Обновленные данные пользователя
	User User `json:"user"`
}
