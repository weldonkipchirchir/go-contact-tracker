package controllers

import (
	"contact/models"
	"contact/services"
	"contact/token"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.RegisterUser(&user)
	if err == services.ErrEmailTaken {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
	var loginRequest struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"error":   "Invalid json format",
			"message": err.Error(),
		})
		return
	}

	user, err := services.LoginUser(loginRequest.Email, loginRequest.Password)
	if err == services.ErrUserNotFound || err == services.ErrInvalidCredentials {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	token, refreshToken, err := token.TokenGenerator(user.ID, user.Name, user.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Create a UserResponse without the password
	userResponse := models.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
	}

	res := gin.H{
		"access_token":  token,
		"refresh_token": refreshToken,
		"user":          userResponse,
	}

	// Include the "Login successful" message in the response
	res["message"] = "Login successful"

	c.SetCookie("token", token, int(time.Hour)*30, "/", "", false, true)
	c.SetCookie("refreshToken", token, int(time.Hour)*90, "/", "", false, true)

	// Send the combined JSON response to the client
	c.JSON(http.StatusOK, res)
}

// Logout function
func Logout(c *gin.Context) {
	// Clear cookies
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.SetCookie("refreshToken", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
