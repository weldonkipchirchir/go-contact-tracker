package routes

import (
	"contact/controllers"
	"contact/db"
	"contact/handlers"
	"contact/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	errorHandler := handlers.NewErrorHandler() // Create an instance of ErrorHandler

	contactHandler := controllers.NewContactCollection(db.GetCollection("contacts"), errorHandler)

	contactGroup := router.Group("/api/v1/contacts")

	contactGroup.Use(middleware.Authentication())
	contactGroup.POST("/contacts", contactHandler.CreateContact)
	contactGroup.GET("/contacts", contactHandler.GetContacts)
	contactGroup.GET("/contacts/:id", contactHandler.GetContactByID)
	contactGroup.PUT("/contacts/:id", contactHandler.UpdateContact)
	contactGroup.DELETE("/contacts/:id", contactHandler.DeleteContact)
}
