package app

import (
	"GIN/configs"
	db "GIN/db/sqlc"
	"GIN/internal/handler"
	"GIN/internal/repository"
	"GIN/internal/routes"
	"GIN/internal/service"
	"GIN/pkg/token"
)

type UserModule struct {
	Routes routes.UserRoutes
}

func NewUserModule(store db.Store, config *configs.Config, tokenMaker token.TokenMaker) UserModule {

	userRepo := repository.NewUserRepository(store)
	roleRepo := repository.NewRoleRepository(store)
	userService := service.NewUserService(userRepo, roleRepo, tokenMaker, config)
	userHandler := handler.NewUserHandler(userService)
	userRoutes := routes.NewUserRoutes(userHandler, tokenMaker)
	return UserModule{Routes: userRoutes}
}
