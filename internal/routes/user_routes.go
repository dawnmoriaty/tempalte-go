package routes

import (
	"GIN/internal/handler"
	"GIN/internal/middleware"
	"GIN/pkg/token"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type UserRoutes struct {
	handler    handler.UserHandler
	tokenMaker token.TokenMaker
}

func NewUserRoutes(h handler.UserHandler, t token.TokenMaker) UserRoutes {
	return UserRoutes{handler: h, tokenMaker: t}
}

func (r *UserRoutes) Setup(engine *gin.Engine) {
	// Public routes - không cần authentication
	publicRoutes := engine.Group("/api/v1/auth")
	{
		publicRoutes.POST("/register", r.handler.CreateUser)
		publicRoutes.POST("/login", r.handler.Login)
		publicRoutes.POST("/refresh", r.handler.RefreshToken)
	}

	// Protected routes - cần authentication với Redis validation
	protectedRoutes := engine.Group("/api/v1/users")
	protectedRoutes.Use(middleware.AuthMiddlewareWithRedis(r.tokenMaker))
	{
		protectedRoutes.GET("/profile", r.handler.GetProfile)
		protectedRoutes.POST("/logout", r.handler.Logout)
		protectedRoutes.POST("/logout-all", r.handler.LogoutAll)
		protectedRoutes.GET("/sessions", r.handler.GetActiveSessions)
	}

	// Health check endpoint
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "user-service-with-redis",
		})
	})
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
