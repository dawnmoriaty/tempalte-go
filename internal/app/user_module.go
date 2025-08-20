package app

import (
	db "GIN/db/sqlc"
	"GIN/internal/handler"
	"GIN/internal/repository"
	"GIN/internal/routes"
	"GIN/internal/service"
)

type UserModule struct {
	Routes routes.UserRoutes
}

func NewUserModule(store db.Store) UserModule {
	userRepo := repository.NewUserRepository(store)
	roleRepo := repository.NewRoleRepository(store)
	userService := service.NewUserService(userRepo, roleRepo)
	userHandler := handler.NewUserHandler(userService)
	userRoutes := routes.NewUserRoutes(userHandler)

	return UserModule{Routes: userRoutes}
}
