package app

import (
	"GIN/configs"
	db "GIN/db/sqlc"
	"GIN/internal/handler"
	"GIN/internal/repository"
	"GIN/internal/routes"
	"GIN/internal/service"
	"GIN/pkg/token"
	"log"
)

type UserModule struct {
	Routes routes.UserRoutes
}

func NewUserModule(store db.Store, config *configs.Config) UserModule {
	tokenMaker, err := token.NewJwtMaker(config.JWT.AccessTokenSecret)
	if err != nil {
		log.Fatalf("cannot create token maker: %v", err)
	}
	userRepo := repository.NewUserRepository(store)
	roleRepo := repository.NewRoleRepository(store)
	userService := service.NewUserService(userRepo, roleRepo, tokenMaker, config)
	userHandler := handler.NewUserHandler(userService)
	userRoutes := routes.NewUserRoutes(userHandler)

	return UserModule{Routes: userRoutes}
}
