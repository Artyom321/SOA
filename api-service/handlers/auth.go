package handlers

import (
	"log"
	"net/http"
	"net/url"
	"social-network/common/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		client := &http.Client{}

		authCheckURL := url.URL{
			Scheme: h.UserServiceURL.Scheme,
			Host:   h.UserServiceURL.Host,
			Path:   "/users/profile",
		}

		req, err := http.NewRequest(http.MethodGet, authCheckURL.String(), nil)
		if err != nil {
			log.Printf("Error creating auth check request: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError,
				models.ErrorResponse{Error: "Internal server error"})
			return
		}

		for _, cookie := range c.Request.Cookies() {
			req.AddCookie(cookie)
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error checking authentication: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				models.ErrorResponse{Error: "Authentication failed"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("User service returned status: %d", resp.StatusCode)
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				models.ErrorResponse{Error: "Unauthorized"})
			return
		}

		c.Next()
	}
}
