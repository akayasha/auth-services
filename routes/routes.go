package routes

import (
	"auth-services/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Group authentication routes under /api/auth
	authRoutes := r.Group("/api/auth")
	{
		authRoutes.POST("/register", controllers.RegisterUser)
		authRoutes.POST("/login", controllers.LoginUser)
		authRoutes.POST("/verify-email", controllers.VerifyEmail)
		authRoutes.POST("/resend-otp", controllers.ResendOTP)
	}

	return r
}
