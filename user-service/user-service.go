// @title User Service API
// @version 1.0
// @description API сервис для управления пользователями
// @host localhost:8081
// @BasePath /api
// @schemes http
// @securityDefinitions.apikey sessionAuth
// @in cookie
// @name user-session
package main

import (
	"fmt"
	"log"
	"os"
	"social-network/common/config"
	"social-network/user-service/handlers"

	_ "social-network/docs/user"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	serviceConfig := config.LoadConfig("common/config/config.json")
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))

	db = connectToDB(connStr)

	initDB()

	handler := handlers.NewHandler(db)

	store := cookie.NewStore([]byte(os.Getenv("COOKIE_SECRET")))
	router := gin.Default()
	router.Use(sessions.Sessions("user-session", store))

	api := router.Group("/users")
	{
		api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		api.POST("/register", handler.RegisterHandler)
		api.POST("/login", handler.LoginHandler)

		authenticated := api.Group("/")
		authenticated.Use(handler.AuthMiddleware())
		{
			authenticated.GET("/profile", handler.ProfileGetHandler)
			authenticated.PUT("/profile", handler.ProfileUpdateHandler)
		}
	}

	log.Printf("User-service running on port %d\n", serviceConfig.UserService.Port)
	router.Run(fmt.Sprintf(":%d", serviceConfig.UserService.Port))
}
