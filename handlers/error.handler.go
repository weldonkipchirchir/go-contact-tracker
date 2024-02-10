package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorHandler struct {
}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

func (h *ErrorHandler) HandleNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"error": "Not Found",
	})
}

func (h *ErrorHandler) HandleBadRequest(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "Bad Request",
	})
}

func (h *ErrorHandler) HandleInternalServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "Internal Server Error",
	})
}

func (h *ErrorHandler) HandleUnauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "Unauthorized",
	})
}

func (e *ErrorHandler) HandleMethodNotAllowed(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{
		"error": "Method Not Allowed",
	})
}
