package routes

import (
	"contact/controllers"

	"github.com/gin-gonic/gin"
)

func SetUpUsers(router *gin.Engine) {
	users := router.Group("/api/v1/users/")
	{
		// login users
		users.POST("/login", controllers.Login)
		// Create a new user
		users.POST("/register", controllers.Register)
	}
}
