package routes

import (
	"auth-services/controllers"
	"auth-services/middlewares"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Group authentication routes under /api/auth
	authRoutes := r.Group("/api/auth")
	{
		authRoutes.POST("/register", controllers.RegisterUser)
		authRoutes.POST(
			"/login",
			middlewares.RateLimit(5, 1*time.Minute),
			controllers.LoginUser,
		)

		authRoutes.POST(
			"/verify-email",
			middlewares.RateLimit(3, 5*time.Minute),
			controllers.VerifyEmail,
		)
		authRoutes.POST("/resend-otp", controllers.ResendOTP)
		authRoutes.POST("/auth/refresh", controllers.RefreshToken)
	}

	return r
}
