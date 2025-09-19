package routes

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"rentx/models"
	"rentx/utils"

	"github.com/gin-gonic/gin"
)

func refreshTokenHandler(ctx *gin.Context) {
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Refresh token required"})
		return
	}

	// Fetch refresh token record
	rt, err := models.GetRefreshToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired refresh token"})
		return
	}

	// Fetch associated user
	user, err := models.GetUserByID(rt.UserId)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
		return
	}

	// Generate new access token
	token, err := utils.GenerateToken(user.Id, user.Email, user.Phone, user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token"})
		return
	}

	// Return full user info with tokens
	ctx.JSON(http.StatusOK, gin.H{
		"message":      "Authenticated successfully.",
		"id":           user.Id,
		"name":         user.Name,
		"image":        user.Image,
		"email":        user.Email,
		"phone":        user.Phone,
		"token":        token,
		"refreshToken": rt.Token,
		"role":         user.Role,
	})
}

func emailAuthHandler(ctx *gin.Context) {
	var incomingUser models.User
	if err := ctx.ShouldBindJSON(&incomingUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to parse request data."})
		return
	}

	// === STEP 1: Send OTP (future) ===
	// TODO(OTP): In future — if "otp" field not present in request,
	// generate OTP, send to email, and return "OTP sent" response here.
	// return
	// otp := ctx.PostForm("otp") // or from JSON depending on your request
	// if otp == "" {
	// 	// === STEP 1: Send OTP ===
	// 	_, err := utils.SendEmailOTP(incomingUser.Email)
	// 	if err != nil {
	// 		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to send OTP."})
	// 		return
	// 	}
	// 	ctx.JSON(http.StatusOK, gin.H{"message": "OTP sent to email."})
	// 	return
	// }

	// === STEP 2: Verify OTP (future) ===
	// TODO(OTP): In future — if "otp" field present, verify it here.
	// If OTP invalid → return 401 Unauthorized.
	// If OTP valid → proceed to create or login user below.
	// if err := utils.VerifyEmailOTP(incomingUser.Email, otp); err != nil {
	// 	ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired OTP."})
	// 	return
	// }

	// === STEP 3: Existing logic (create or login user) ===
	existingUser := &models.User{Email: incomingUser.Email}
	err := existingUser.LoadByEmail()

	var user models.User

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// User doesn't exist → Signup
			if err := incomingUser.Save(); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user."})
				return
			}
			user = incomingUser
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Database error."})
			return
		}
	} else {
		// User exists → Login
		if !utils.ComparePasswords(incomingUser.Password, existingUser.Password) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials."})
			return
		}
		user = *existingUser
	}

	// === STEP 4: Generate token ===
	token, err := utils.GenerateToken(user.Id, user.Email, "", user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token."})
		return
	}

	// === STEP 5: Generate refresh token ===
	refreshToken, err := models.NewRefreshToken(user.Id, 30) // 30 days valid
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create refresh token."})
		return
	}
	if err := refreshToken.Save(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save refresh token."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":      "Authenticated successfully.",
		"id":           user.Id,
		"name":         user.Name,
		"image":        user.Image,
		"email":        user.Email,
		"phone":        user.Phone,
		"token":        token,              // short-lived JWT
		"refreshToken": refreshToken.Token, // long-lived token
		"role":         user.Role,
	})
}

func phoneAuthHandler(ctx *gin.Context) {
	var incomingUser models.User
	if err := ctx.ShouldBindJSON(&incomingUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to parse request data."})
		return
	}
	// === STEP 1: Send OTP ===
	// TODO: In future — if "otp" field not present in request,
	// generate OTP, send to phone, and return "OTP sent" response here.
	// return
	// otp := ctx.PostForm("otp") // or from JSON depending on your request
	// if otp == "" {
	// 	// === STEP 1: Send OTP ===
	// 	_, err := utils.SendOTP(incomingUser.Phone)
	// 	if err != nil {
	// 		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to send OTP."})
	// 		return
	// 	}
	// 	ctx.JSON(http.StatusOK, gin.H{"message": "OTP sent to phone."})
	// 	return
	// }

	// === STEP 2: Verify OTP ===
	// === STEP 2: Verify OTP ===
	// TODO: In future — if "otp" field present, verify it here.
	// If OTP invalid → return 401 Unauthorized.
	// If OTP valid → proceed to create or login user below.
	// if err := utils.VerifyOTP(incomingUser.Phone, otp); err != nil {
	// 	ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired OTP."})
	// 	return
	// }

	// === STEP 3: Existing logic (create or login user) ===
	existingUser := &models.User{Phone: incomingUser.Phone}
	err := existingUser.LoadByPhone()

	var user models.User

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if err := incomingUser.Save(); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user."})
				return
			}
			user = incomingUser
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Database error."})
			return
		}
	} else {
		user = *existingUser
	}

	// === STEP 4: Generate token ===
	token, err := utils.GenerateToken(user.Id, "", user.Phone, user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token."})
		return
	}

	// === STEP 5: Generate refresh token ===
	refreshToken, err := models.NewRefreshToken(user.Id, 30) // 30 days valid
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create refresh token."})
		return
	}
	if err := refreshToken.Save(); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save refresh token."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":      "Authenticated successfully.",
		"id":           user.Id,
		"name":         user.Name,
		"image":        user.Image,
		"email":        user.Email,
		"phone":        user.Phone,
		"token":        token,              // short-lived JWT
		"refreshToken": refreshToken.Token, // long-lived token
		"role":         user.Role,
	})
}

func oauthAuthHandler(ctx *gin.Context) {
	var incomingUser struct {
		Email string `json:"email"`
	}
	if err := ctx.ShouldBindJSON(&incomingUser); err != nil || incomingUser.Email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Email is required"})
		return
	}

	// Step 1: Check if user exists
	existingUser := &models.User{Email: incomingUser.Email}
	err := existingUser.LoadByEmail()

	var user models.User

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// User doesn't exist → create
			newUser := models.User{
				Email: incomingUser.Email,
				Role:  "user", // default role
			}
			if err := newUser.Save(); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user."})
				return
			}
			user = newUser
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Database error."})
			return
		}
	} else {
		// User exists → use it
		user = *existingUser
	}

	// Step 2: Generate access token
	token, err := utils.GenerateToken(user.Id, user.Email, "", user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token."})
		return
	}

	// Step 3: Generate refresh token (30 days)
	refreshToken, err := models.NewRefreshToken(user.Id, 30)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create refresh token."})
		return
	}
	if err := refreshToken.Save(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save refresh token."})
		return
	}

	// Step 4: Return full user info
	ctx.JSON(http.StatusOK, gin.H{
		"message":      "Authenticated successfully.",
		"id":           user.Id,
		"name":         user.Name,
		"image":        user.Image,
		"email":        user.Email,
		"phone":        user.Phone,
		"token":        token,
		"refreshToken": refreshToken.Token,
		"role":         user.Role,
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
