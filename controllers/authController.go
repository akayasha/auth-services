package controllers

import (
	"auth-services/services"
	"auth-services/utils"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// RegisterUser handles user registration
func RegisterUser(c *gin.Context) {
	var registerData struct {
		FirstName string `json:"firstName" binding:"required"`
		LastName  string `json:"lastName" binding:"required"`
		Username  string `json:"username" binding:"required"`
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required"`
		Role      string `json:"role"`
		Dob       string `json:"dob" binding:"required"`
	}

	if err := c.ShouldBindJSON(&registerData); err != nil {
		fmt.Println("Register bind error:", err)
		utils.RespondError(c, 400, "Invalid request data")
		return
	}

	dob, err := time.Parse("2006-01-02", registerData.Dob)
	if err != nil {
		utils.RespondError(c, 400, "Invalid date format, use YYYY-MM-DD")
		return
	}

	user, err := services.RegisterUser(
		registerData.Username,
		registerData.Email,
		registerData.Password,
		registerData.Role,
		registerData.FirstName,
		registerData.LastName,
		dob,
	)
	if err != nil {
		utils.RespondError(c, 400, err.Error())
		return
	}

	utils.Respond(c, 201, "User registered successfully", user)
}

// Resend OTP
func ResendOTP(c *gin.Context) {
	var resendData struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&resendData); err != nil {
		utils.RespondError(c, 400, "Email is required")
		return
	}

	if err := services.ResendOTP(resendData.Email); err != nil {
		utils.RespondError(c, 400, err.Error())
		return
	}

	utils.Respond(c, 200, "OTP resent successfully", nil)
}

// Login
func LoginUser(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		utils.RespondError(c, 400, "Invalid login request")
		return
	}

	response, err := services.LoginUser(loginData.Email, loginData.Password)
	if err != nil {
		utils.RespondError(c, 401, err.Error())
		return
	}

	utils.Respond(c, 200, "Login successful", response)
}

// Verify Email (OTP)
func VerifyEmail(c *gin.Context) {
	var verifyData struct {
		Email string `json:"email" binding:"required,email"`
		OTP   string `json:"otp" binding:"required"`
	}

	if err := c.ShouldBindJSON(&verifyData); err != nil {
		utils.RespondError(c, 400, "Invalid request data")
		return
	}

	message, err := services.VerifyEmail(verifyData.Email, verifyData.OTP)
	if err != nil {
		utils.RespondError(c, 400, err.Error())
		return
	}

	utils.Respond(c, 200, message, nil)
}

func RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, 400, "Refresh token required")
		return
	}

	tokens, err := services.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		utils.RespondError(c, 401, err.Error())
		return
	}

	utils.Respond(c, 200, "Token refreshed", tokens)
}
