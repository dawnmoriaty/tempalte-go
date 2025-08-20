package routes

import (
	"GIN/internal/handler"
	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	handler handler.UserHandler
}

func NewUserRoutes(h handler.UserHandler) UserRoutes {
	return UserRoutes{handler: h}
}

func (r *UserRoutes) Setup(engine *gin.Engine) {
	userGroup := engine.Group("/api/v1/users")
	{
		userGroup.POST("/", r.handler.CreateUser)
		userGroup.POST("/login", r.handler.Login)
	}
}
