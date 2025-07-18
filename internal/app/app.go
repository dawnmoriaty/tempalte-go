package app

import (
	"GIN/configs"
	"GIN/db/sqlc"
	"GIN/internal/middleware"
	"GIN/pkg/database" // Import package database của bạn
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
    Config *configs.DatabaseConfig // Thay bằng struct config của bạn
}

func NewApplication(cfg *configs.DatabaseConfig) *Application {
    // Kết nối DB và tạo store
    dbPool := database.Connect(cfg)
    store := db.NewStore(dbPool)

    // Khởi tạo gin engine
    engine := gin.Default()
	// === Cấu hình middleware ===
	engine.Use(middleware.CORSMiddleware())

	// cấu hình endpoint Tự động chuyển hướng nếu URL bị thiếu hoặc thừa dấu gạch chéo ở cuối
	engine.RedirectTrailingSlash = true
    // === Khởi tạo các module ===
    userModule := NewUserModule(store)
    // Thêm các module khác ở đây...

    // === Đăng ký routes từ các module ===
    userModule.Routes.Setup(engine)
    // Đăng ký các routes khác...

    return &Application{
        Engine: engine,
        Config: cfg,
    }
}
// graceful shutdown như bên springboot
func (app *Application) Run() error {
	// Tạo một http.Server từ Gin engine
	srv := &http.Server{
		Addr:    ":" + app.Config.HTTPPort,
		Handler: app.Engine,
	}

	// Chạy server trong một goroutine riêng để không bị block
	go func() {
		zap.S().Infof("Server is running on port %s", app.Config.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.S().Fatalf("listen: %s\n", err)
		}
	}()

	// Tạo một channel để lắng nghe tín hiệu tắt từ OS (Ctrl+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block cho đến khi nhận được tín hiệu
	<-quit
	zap.S().Info("Shutting down server...")

	// Tạo một context với timeout để cho server 5 giây hoàn thành các request hiện tại
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.S().Fatal("Server forced to shutdown:", err)
		return err
	}

	zap.S().Info("Server exiting")
	return nil
}
