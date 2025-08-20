package app

import (
	"GIN/configs"
	"GIN/db/sqlc"
	"GIN/internal/middleware"
	"GIN/pkg/database"
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Application struct {
	Engine *gin.Engine
	Config *configs.Config
}

func NewApplication(cfg *configs.Config) *Application {
	dbPool := database.Connect(&cfg.Database)
	store := db.NewStore(dbPool)
	//redis.ConnectRedis()
	engine := gin.Default()

	engine.Use(middleware.CORSMiddleware())
	userModule := NewUserModule(store)
	userModule.Routes.Setup(engine)

	return &Application{
		Engine: engine,
		Config: cfg,
	}
}

func (app *Application) Run() error {

	srv := &http.Server{
		Addr:    ":" + app.Config.Database.HTTPPort,
		Handler: app.Engine,
	}

	go func() {
		zap.S().Infof("Server is running on port %s", app.Config.Database.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.S().Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.S().Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.S().Fatal("Server forced to shutdown:", err)
		return err
	}

	zap.S().Info("Server exiting")
	return nil
}
