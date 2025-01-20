package controllers

import (
	"auth-services/services"
	"auth-services/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

// RegisterUser handles user registration
func RegisterUser(c *gin.Context) {
	var registerData struct {
		FirstName string `json:"firstName" binding:"required"`
		LastName  string `json:"lastName" binding:"required"`
		Username  string `json:"username" binding:"required"`
		Email     string `json:"email" binding:"required"`
		Password  string `json:"password" binding:"required"`
		Role      string `json:"role"`
		Dob       string `json:"dob" binding:"required"`
	}

	if err := c.BindJSON(&registerData); err != nil {
		utils.RespondError(c, 400, "Invalid request data")
		fmt.Errorf("error", err)
		return
	}

	// Parse the date using the format "2006-01-02"
	dob, err := time.Parse("2006-01-02", registerData.Dob)
	if err != nil {
		utils.RespondError(c, 400, "Invalid date format, use YYYY-MM-DD")
		fmt.Println("Date parsing error:", err)
		return
	}

	user, err := services.RegisterUser(registerData.Username, registerData.Email, registerData.Password, registerData.Role, registerData.FirstName, registerData.LastName, dob)
	if err != nil {
		utils.RespondError(c, 400, err.Error())
		return
	}

	utils.Respond(c, 201, "User registered successfully", user)
}

func ResendOTP(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		utils.RespondError(c, 400, "Email is required")
		return
	}

	err := services.ResendOTP(email)
	if err != nil {
		utils.RespondError(c, 400, err.Error())
		return
	}

	utils.Respond(c, 200, "OTP resent successfully", nil)
}

// LoginUser handles user login
func LoginUser(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&loginData); err != nil {
		utils.RespondError(c, 400, "Invalid request data")
		return
	}

	response, err := services.LoginUser(loginData.Email, loginData.Password)
	if err != nil {
		utils.RespondError(c, 401, err.Error())
		return
	}

	utils.Respond(c, 200, "Login successful", response)
}

// VerifyEmail handles email verification using OTP
func VerifyEmail(c *gin.Context) {
	var verifyData struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := c.BindJSON(&verifyData); err != nil {
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
