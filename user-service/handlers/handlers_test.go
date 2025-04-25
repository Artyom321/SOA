package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"social-network/common/config"
	"social-network/common/models"
)

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	// Используем уникальное имя базы данных для каждого теста
	dbName := fmt.Sprintf("file:memdb%d?mode=memory&cache=shared", time.Now().UnixNano())
	testDB, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	assert.NoError(t, err)

	err = testDB.AutoMigrate(&models.User{})
	assert.NoError(t, err)

	// Получаем базовый sql.DB объект для закрытия соединения
	sqlDB, err := testDB.DB()
	assert.NoError(t, err)

	// Возвращаем функцию для закрытия соединения после теста
	return testDB, func() {
		sqlDB.Close()
	}
}

// Создаем тестовый обработчик без Kafka для тестов
func newTestHandler(db *gorm.DB) *Handler {
	return &Handler{
		DB:            db,
		KafkaProducer: nil,
		Config:        &config.Config{},
	}
}

func TestRegisterHandler(t *testing.T) {
	testDB, cleanup := setupTestDB(t)
	defer cleanup()

	handler := newTestHandler(testDB)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/register", handler.RegisterHandler)

	t.Run("Successful registration", func(t *testing.T) {
		w := httptest.NewRecorder()

		reqBody, _ := json.Marshal(models.RegisterRequest{
			Login:    "newuser",
			Email:    "new@example.com",
			Password: "password123",
		})
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var user models.User
		testDB.Where("login = ?", "newuser").First(&user)
		assert.Equal(t, "newuser", user.Login)
		assert.Equal(t, "new@example.com", user.Email)
	})

	t.Run("Invalid registration data", func(t *testing.T) {
		w := httptest.NewRecorder()

		reqBody, _ := json.Marshal(models.RegisterRequest{
			Login:    "u",
			Email:    "invalid-email",
			Password: "123",
		})
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLoginHandler(t *testing.T) {
	testDB, cleanup := setupTestDB(t)
	defer cleanup()

	handler := newTestHandler(testDB)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	testUser := models.User{
		Login:        "testuser",
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
	}
	testDB.Create(&testUser)

	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Используем одинаковый секрет для store во всех тестах
	store := cookie.NewStore([]byte("test_secret"))
	r.Use(sessions.Sessions("user-session", store))
	r.POST("/login", handler.LoginHandler)

	t.Run("Successful login", func(t *testing.T) {
		w := httptest.NewRecorder()

		reqBody, _ := json.Marshal(models.LoginRequest{
			Login:    "testuser",
			Password: "password123",
		})
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		cookies := w.Result().Cookies()
		var sessionCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "user-session" {
				sessionCookie = cookie
				break
			}
		}
		assert.NotNil(t, sessionCookie, "Session cookie should be set")
	})

	t.Run("Invalid password", func(t *testing.T) {
		w := httptest.NewRecorder()

		reqBody, _ := json.Marshal(models.LoginRequest{
			Login:    "testuser",
			Password: "wrongpassword",
		})
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Non-existent user", func(t *testing.T) {
		w := httptest.NewRecorder()

		reqBody, _ := json.Marshal(models.LoginRequest{
			Login:    "nonexistent",
			Password: "password123",
		})
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestProfileGetHandler(t *testing.T) {
	testDB, cleanup := setupTestDB(t)
	defer cleanup()

	handler := newTestHandler(testDB)

	testUser := models.User{
		ID:        1,
		Login:     "testuser",
		Email:     "test@example.com",
		Name:      "Test",
		Surname:   "User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	testDB.Create(&testUser)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	store := cookie.NewStore([]byte("test_secret"))
	r.Use(sessions.Sessions("user-session", store))

	// Создаем отдельный обработчик для установки сессии
	r.GET("/set-session/:id", func(c *gin.Context) {
		id := c.Param("id")
		userID, _ := strconv.ParseUint(id, 10, 64)

		session := sessions.Default(c)
		session.Set("user_id", userID)
		err := session.Save()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusOK)
	})

	r.GET("/profile", handler.AuthMiddleware(), handler.ProfileGetHandler)

	t.Run("Successful profile retrieval", func(t *testing.T) {
		// Создаем запрос для установки сессии
		w1 := httptest.NewRecorder()
		req1, _ := http.NewRequest("GET", "/set-session/1", nil)
		r.ServeHTTP(w1, req1)

		assert.Equal(t, http.StatusOK, w1.Code)

		// Получаем cookie из ответа
		cookies := w1.Result().Cookies()

		// Создаем запрос на получение профиля с теми же куками
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("GET", "/profile", nil)
		for _, cookie := range cookies {
			req2.AddCookie(cookie)
		}

		r.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusOK, w2.Code)

		var response models.ProfileResponse
		err := json.Unmarshal(w2.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "testuser", response.User.Login)
		assert.Equal(t, "Test", response.User.Name)
		assert.Equal(t, "User", response.User.Surname)
	})

	t.Run("Without session", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/profile", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Non-existent user", func(t *testing.T) {
		// Создаем запрос для установки сессии с несуществующим ID
		w1 := httptest.NewRecorder()
		req1, _ := http.NewRequest("GET", "/set-session/999", nil)
		r.ServeHTTP(w1, req1)

		cookies := w1.Result().Cookies()

		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("GET", "/profile", nil)
		for _, cookie := range cookies {
			req2.AddCookie(cookie)
		}

		r.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusUnauthorized, w2.Code)
	})
}

func TestProfileUpdateHandler(t *testing.T) {
	testDB, cleanup := setupTestDB(t)
	defer cleanup()

	handler := newTestHandler(testDB)

	testUser := models.User{
		ID:        1,
		Login:     "testuser",
		Email:     "test@example.com",
		Name:      "OldName",
		Surname:   "OldSurname",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	testDB.Create(&testUser)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	store := cookie.NewStore([]byte("test_secret"))
	r.Use(sessions.Sessions("user-session", store))

	r.GET("/set-session/:id", func(c *gin.Context) {
		id := c.Param("id")
		userID, _ := strconv.ParseUint(id, 10, 64)

		session := sessions.Default(c)
		session.Set("user_id", userID)
		err := session.Save()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusOK)
	})

	r.PUT("/profile", handler.AuthMiddleware(), handler.ProfileUpdateHandler)

	// Улучшенная функция для получения сессионных cookie
	getSessionCookies := func(userID uint64) []*http.Cookie {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/set-session/"+strconv.FormatUint(userID, 10), nil)
		r.ServeHTTP(w, req)
		return w.Result().Cookies()
	}

	t.Run("Successful profile update", func(t *testing.T) {
		cookies := getSessionCookies(1)

		newName := "NewName"
		newSurname := "NewSurname"
		updatedRequest := models.ProfileUpdateRequest{
			Name:    &newName,
			Surname: &newSurname,
		}

		reqBody, _ := json.Marshal(updatedRequest)

		req, _ := http.NewRequest("PUT", "/profile", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.ProfileUpdateResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "Profile updated successfully", response.Status)
		assert.Equal(t, newName, response.User.Name)
		assert.Equal(t, newSurname, response.User.Surname)

		var updatedUser models.User
		testDB.First(&updatedUser, 1)
		assert.Equal(t, newName, updatedUser.Name)
		assert.Equal(t, newSurname, updatedUser.Surname)
	})

	t.Run("No changes when all fields match", func(t *testing.T) {
		// Обновляем пользователя в базе данных
		exactName := "ExactName"
		exactSurname := "ExactSurname"
		exactEmail := "exact@example.com"
		exactPhone := "+79001234567"
		birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

		testDB.Model(&models.User{}).Where("id = ?", 1).Updates(map[string]interface{}{
			"name":         exactName,
			"surname":      exactSurname,
			"email":        exactEmail,
			"phone_number": exactPhone,
			"birth_date":   birthDate,
		})

		var userInDB models.User
		testDB.First(&userInDB, 1)
		assert.Equal(t, exactName, userInDB.Name)

		// Получаем сессионные куки для пользователя
		cookies := getSessionCookies(1)

		// Создаем запрос на обновление с теми же данными
		updateRequest := models.ProfileUpdateRequest{
			Name:        &exactName,
			Surname:     &exactSurname,
			Email:       &exactEmail,
			PhoneNumber: &exactPhone,
			BirthDate:   &birthDate,
		}

		reqBody, _ := json.Marshal(updateRequest)

		req, _ := http.NewRequest("PUT", "/profile", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.ProfileUpdateResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "No fields were updated", response.Status)

		var userAfterRequest models.User
		testDB.First(&userAfterRequest, 1)
		// Убедимся, что дата обновления не изменилась
		assert.Equal(t, userInDB.UpdatedAt.Unix(), userAfterRequest.UpdatedAt.Unix())
	})
}
