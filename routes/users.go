package routes

import (
	"fmt"
	"net/http"
	"rentx/models"
	"rentx/utils"

	"github.com/gin-gonic/gin"
)

func signup(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to parse request data."})
		return
	}
	err = user.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user."})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "User created."})
}

func login(context *gin.Context) {
	var user models.User
	if err := context.ShouldBindJSON(&user); err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to parse request data."})
		return
	}

	if err := user.ValidateCredentials(); err != nil {
		fmt.Println(err)
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials."})
		return
	}

	token, err := utils.GenerateToken(user.Id, user.Email, user.Role)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token."})
		return
	}

	context.JSON(http.StatusAccepted, gin.H{
		"message": "Logged in.",
		"id":      user.Id,
		"name":    user.Name,
		"image":   user.Image,
		"email":   user.Email,
		"phone":   user.Phone,
		"token":   token,
		"role":    user.Role,
	})
}
