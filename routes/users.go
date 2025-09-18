package routes

import (
	"database/sql"
	"errors"
	"net/http"
	"rentx/models"
	"rentx/utils"

	"github.com/gin-gonic/gin"
)

func authHandler(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to parse request data."})
		return
	}

	// Check if user already exists
	existingUser := &models.User{Email: user.Email}
	err := existingUser.LoadByEmail()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// User doesn't exist → Signup
			if err := user.Save(); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user."})
				return
			}
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Database error."})
			return
		}
	} else {
		// User exists → Login
		if !utils.ComparePasswords(user.Password, existingUser.Password) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials."})
			return
		}
		user = *existingUser
	}

	// Generate token
	token, err := utils.GenerateToken(user.Id, user.Email, user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Authenticated successfully.",
		"id":      user.Id,
		"name":    user.Name,
		"image":   user.Image,
		"email":   user.Email,
		"phone":   user.Phone,
		"token":   token,
		"role":    user.Role,
	})
}

// func signup(context *gin.Context) {
// 	var user models.User
// 	err := context.ShouldBindJSON(&user)
// 	if err != nil {
// 		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to parse request data."})
// 		return
// 	}
// 	err = user.Save()
// 	if err != nil {
// 		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user."})
// 		return
// 	}
// 	context.JSON(http.StatusCreated, gin.H{"message": "User created."})
// }

// func login(context *gin.Context) {
// 	var user models.User
// 	if err := context.ShouldBindJSON(&user); err != nil {
// 		fmt.Println(err)
// 		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to parse request data."})
// 		return
// 	}

// 	if err := user.ValidateCredentials(); err != nil {
// 		fmt.Println(err)
// 		context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials."})
// 		return
// 	}

// 	token, err := utils.GenerateToken(user.Id, user.Email, user.Role)
// 	if err != nil {
// 		fmt.Println(err)
// 		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token."})
// 		return
// 	}

// 	context.JSON(http.StatusAccepted, gin.H{
// 		"message": "Logged in.",
// 		"id":      user.Id,
// 		"name":    user.Name,
// 		"image":   user.Image,
// 		"email":   user.Email,
// 		"phone":   user.Phone,
// 		"token":   token,
// 		"role":    user.Role,
// 	})
// }

// ValidateCredentials checks email and password
// func (u *User) ValidateCredentials() error {
// 	query := "SELECT id, name, email, phone, password, image, role FROM users WHERE email = ?"
// 	row := db.DB.QueryRow(query, u.Email)

// 	var dbPass string
// 	err := row.Scan(&u.Id, &u.Name, &u.Email, &u.Phone, &dbPass, &u.Image, &u.Role)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return errors.New("user not found")
// 		}
// 		return err
// 	}

// 	if !utils.ComparePasswords(u.Password, dbPass) {
// 		return errors.New("invalid credentials")
// 	}
// 	return nil
// }
