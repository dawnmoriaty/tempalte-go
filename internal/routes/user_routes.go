package routes

import (
	"GIN/internal/handler"
	"GIN/internal/middleware"
	"GIN/pkg/token"
	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	handler    handler.UserHandler
	tokenMaker token.TokenMaker
}

func NewUserRoutes(h handler.UserHandler, t token.TokenMaker) UserRoutes {
	return UserRoutes{handler: h, tokenMaker: t}
}

func (r *UserRoutes) Setup(engine *gin.Engine) {

	publicRoutes := engine.Group("/api/v1/users")
	{
		publicRoutes.POST("/register", r.handler.CreateUser)
		publicRoutes.POST("/login", r.handler.Login)
	}
	protectedRoutes := engine.Group("/api/v1/users")
	protectedRoutes.Use(middleware.AuthMiddleware(r.tokenMaker))
	{
		protectedRoutes.GET("/profile", r.handler.GetProfile)
	}
}
