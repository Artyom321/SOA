package handlers

import (
	"net/http"
	"social-network/common/models"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

// @Summary Регистрация нового пользователя
// @Description Создает нового пользователя в системе
// @Tags Authentication
// @Accept json
// @Produce json
// @Param body body models.RegisterRequest true "Данные для регистрации"
// @Success 201 {object} models.RegisterResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /register [post]
func (h *Handler) RegisterHandler(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	var existingUser models.User
	if result := h.DB.Where("login = ? OR email = ?", req.Login, req.Email).First(&existingUser); result.Error == nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "User with this login or email already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Internal error"})
		return
	}

	user := models.User{
		Login:        req.Login,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "User creation failed"})
		return
	}

	c.JSON(http.StatusCreated, models.RegisterResponse{
		Status: "Registered successfully",
	})
}

// @Summary Аутентификация пользователя
// @Description Выполняет вход пользователя в систему
// @Tags Authentication
// @Accept json
// @Produce json
// @Param body body models.LoginRequest true "Данные для входа"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /login [post]
func (h *Handler) LoginHandler(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	var user models.User
	if h.DB.Where("login = ?", req.Login).First(&user).Error != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid credentials"})
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to save session"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Message: "Logged in successfully",
	})
}

// @Summary Получение профиля пользователя
// @Description Возвращает профиль текущего авторизованного пользователя
// @Tags Profile
// @Produce json
// @Security sessionAuth
// @Success 200 {object} models.ProfileResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /profile [get]
func (h *Handler) ProfileGetHandler(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to retrieve user data"})
		return
	}

	c.JSON(http.StatusOK, models.ProfileResponse{
		User: user,
	})
}

// @Summary Обновление профиля пользователя
// @Description Обновляет данные профиля текущего пользователя
// @Tags Profile
// @Accept json
// @Produce json
// @Security sessionAuth
// @Param body body models.ProfileUpdateRequest true "Новые данные для профиля"
// @Success 200 {object} models.ProfileUpdateResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /profile [put]
func (h *Handler) ProfileUpdateHandler(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to retrieve user data"})
		return
	}

	var input models.ProfileUpdateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	fieldsUpdated := UpdateUserFields(&user, input)

	if !fieldsUpdated {
		c.JSON(http.StatusOK, models.ProfileUpdateResponse{
			Status: "No fields were updated",
			User:   user,
		})
		return
	}

	user.UpdatedAt = time.Now()

	if err := h.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Could not update user"})
		return
	}

	c.JSON(http.StatusOK, models.ProfileUpdateResponse{
		Status: "Profile updated successfully",
		User:   user,
	})
}
