package middleware

import (
	"contact/token"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token not provided"})
			c.Abort()
			return
		}

		claims, msg := token.ValidateToken(accessToken)
		if msg != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			c.Abort()
			return
		}

		// If access token is expired, attempt to refresh it
		if claims.ExpiresAt < time.Now().Unix() {
			refreshToken, err := c.Cookie("refreshToken")
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not provided"})
				c.Abort()
				return
			}

			newAccessToken, err := token.UpdateToken(refreshToken)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
				c.Abort()
				return
			}

			// Set the new access token as a cookie
			c.SetCookie("token", newAccessToken, int(time.Hour)*30, "/", "", false, true)
			c.SetCookie("refreshToken", refreshToken, int(time.Hour)*90, "/", "", false, true)

			c.Set("id", claims.Id)
			c.Set("email", claims.Email)
			c.Next()

		} else {
			c.Set("id", claims.Id)
			c.Set("email", claims.Email)
			c.Next()
		}
	}
}
