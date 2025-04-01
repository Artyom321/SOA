package handlers

import (
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createUserServiceProxy() *httputil.ReverseProxy {
	director := func(req *http.Request) {
		req.URL.Scheme = h.UserServiceURL.Scheme
		req.URL.Host = h.UserServiceURL.Host

		req.URL.Path = strings.Replace(req.URL.Path, "/api/", "/users/", 1)

		log.Printf("Proxying request to user service: %s", req.URL.String())

		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "API-Gateway")
		}
	}

	return &httputil.ReverseProxy{Director: director}
}

// RegisterHandler обрабатывает запросы на регистрацию
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
	h.ProxyUserRequest(c)
}

// LoginHandler обрабатывает запросы на вход в систему
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
	h.ProxyUserRequest(c)
}

// ProfileGetHandler обрабатывает запросы на получение профиля
// @Summary Получение профиля пользователя
// @Description Возвращает профиль текущего авторизованного пользователя
// @Tags Profile
// @Produce json
// @Security sessionAuth
// @Success 200 {object} models.ProfileResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /profile [get]
func (h *Handler) ProfileGetHandler(c *gin.Context) {
	h.ProxyUserRequest(c)
}

// ProfileUpdateHandler обрабатывает запросы на обновление профиля
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
	h.ProxyUserRequest(c)
}

func (h *Handler) ProxyUserRequest(c *gin.Context) {
	h.UserServiceProxy.ServeHTTP(c.Writer, c.Request)
}
