package handlers

import (
	"net/http"
	"social-network/common/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Unauthorized"})
			c.Abort()
			return
		}

		var user models.User
		if err := h.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "User not found"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
