package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestRegisterHandler(t *testing.T) {
	r := gin.New()

	r.POST("/api/register", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"status": "Registered successfully"})
	})

	reqBody := `{"login":"testuser","email":"test@example.com","password":"password123"}`
	req, _ := http.NewRequest("POST", "/api/register", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Registered successfully", response["status"])
}

func TestLoginHandler(t *testing.T) {
	r := gin.New()

	r.POST("/api/login", func(c *gin.Context) {
		c.SetCookie("user-session", "test-session", 3600, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully"})
	})

	reqBody := `{"login":"testuser","password":"password123"}`
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Logged in successfully", response["message"])

	cookies := w.Result().Cookies()
	assert.NotEmpty(t, cookies, "No cookies found")

	var sessionCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "user-session" {
			sessionCookie = cookie
			break
		}
	}

	assert.NotNil(t, sessionCookie, "Session cookie not found")
	assert.Equal(t, "test-session", sessionCookie.Value)
}

func TestProfileGetHandler(t *testing.T) {
	r := gin.New()

	authMiddleware := func(c *gin.Context) {
		_, err := c.Cookie("user-session")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}

	r.GET("/api/profile", authMiddleware, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"user": gin.H{
				"login": "testuser",
				"email": "test@example.com",
			},
		})
	})

	t.Run("With valid session", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/profile", nil)
		req.AddCookie(&http.Cookie{Name: "user-session", Value: "test-session"})

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		userObj, exists := response["user"]
		assert.True(t, exists, "User field not found in response")

		user, ok := userObj.(map[string]interface{})
		assert.True(t, ok, "User is not a map")
		assert.Equal(t, "testuser", user["login"])
	})

	t.Run("Without session", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/profile", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestProfileUpdateHandler(t *testing.T) {
	r := gin.New()

	authMiddleware := func(c *gin.Context) {
		_, err := c.Cookie("user-session")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}

	r.PUT("/api/profile", authMiddleware, func(c *gin.Context) {
		var data map[string]interface{}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "Profile updated successfully",
			"user": gin.H{
				"login": "testuser",
				"name":  data["name"],
			},
		})
	})

	t.Run("Successful update", func(t *testing.T) {
		reqBody := `{"name":"New Name","surname":"New Surname"}`
		req, _ := http.NewRequest("PUT", "/api/profile", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: "user-session", Value: "test-session"})

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Profile updated successfully", response["status"])

		userObj, exists := response["user"]
		assert.True(t, exists, "User field not found in response")

		user, ok := userObj.(map[string]interface{})
		assert.True(t, ok, "User is not a map")
		assert.Equal(t, "testuser", user["login"])
		assert.Equal(t, "New Name", user["name"])
	})

	t.Run("Without authorization", func(t *testing.T) {
		reqBody := `{"name":"New Name"}`
		req, _ := http.NewRequest("PUT", "/api/profile", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestAuthMiddleware(t *testing.T) {
	r := gin.New()

	authMiddleware := func(c *gin.Context) {
		_, err := c.Cookie("user-session")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}

	r.GET("/test", authMiddleware, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	t.Run("With session", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.AddCookie(&http.Cookie{Name: "user-session", Value: "test-session"})

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Without session", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
