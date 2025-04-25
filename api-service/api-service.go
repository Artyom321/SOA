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

	handler := handlers.NewHandler("user-service", cfg.UserService.Port, "post-service", cfg.PostService.Port)

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

			// Post routes
			posts := authenticated.Group("/posts")
			{
				posts.POST("", handler.CreatePostHandler)
				posts.GET("", handler.ListPostsHandler)
				posts.GET("/:id", handler.GetPostHandler)
				posts.PUT("/:id", handler.UpdatePostHandler)
				posts.DELETE("/:id", handler.DeletePostHandler)
				posts.POST("/:id/view", handler.ViewPostHandler)
				posts.POST("/:id/like", handler.LikePostHandler)
				posts.POST("/:id/comments", handler.AddCommentHandler)
				posts.GET("/:id/comments", handler.GetCommentsHandler)
			}
		}
	}

	log.Printf("API-service running on port %d\n", cfg.APIService.Port)
	r.Run(fmt.Sprintf(":%d", cfg.APIService.Port))
}
