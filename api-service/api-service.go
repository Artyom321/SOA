// @title API Gateway Service
// @version 1.0
// @description API шлюз для доступа к сервисам системы
// @host localhost:8080
// @BasePath /api
// @schemes http
// @securityDefinitions.apikey sessionAuth
// @in cookie
// @name user-session
package main

import (
	"fmt"
	"log"
	"social-network/api-service/handlers"
	"social-network/common/config"

	_ "social-network/docs/api"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg := config.LoadConfig("common/config/config.json")

	handler := handlers.NewHandler("user-service", cfg.UserService.Port)

	r := gin.Default()

	api := r.Group("/api")
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

	log.Printf("API-service running on port %d\n", cfg.APIService.Port)
	r.Run(fmt.Sprintf(":%d", cfg.APIService.Port))
}
