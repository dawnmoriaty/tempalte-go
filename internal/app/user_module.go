package app

import (
	"GIN/configs"
	db "GIN/db/sqlc"
	"GIN/internal/handler"
	"GIN/internal/repository"
	"GIN/internal/routes"
	"GIN/internal/service"
)

type UserModule struct {
    Routes *routes.UserRoutes
}

func NewUserModule(store db.Store, cfg *configs.Config) *UserModule {
    userRepo := repository.NewUserRepository(store)
    userService := service.NewUserService(userRepo)
    userHandler := handler.NewUserHandler(userService)
    userRoutes := routes.NewUserRoutes(userHandler)

    return &UserModule{Routes: userRoutes}
}