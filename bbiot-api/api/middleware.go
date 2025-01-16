package api

import (
	"net/http"

	"github.com/GDGVIT/bbiot-backend/internal/auth"
	"github.com/GDGVIT/bbiot-backend/internal/database"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func IsAdminMiddleware(c *gin.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		role := c.GetString("role")

		if role != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "you are not an admin",
			})
		}

		c.Next()
	}

}

func (s *Server) JWTAuthMiddleware(c *gin.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		authorizationString := c.Request.Header.Get("Authorization")

		if authorizationString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"detail": "Authorization header not present",
			})
		}

		user, err := auth.GetUserFromJWTToken(authorizationString, auth.JWTSecret, database.New(s.db))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"detail": err.Error(),
			})
		}

		c.Set("user", user)
		c.Set("role", user.Role)
		c.Next()
	}
}
