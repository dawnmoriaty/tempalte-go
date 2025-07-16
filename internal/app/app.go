package app

import (
	"GIN/configs"
	"GIN/db/sqlc"
	"GIN/pkg/database" // Import package database của bạn
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
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
		log.Printf("Server is running on port %s", app.Config.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Tạo một channel để lắng nghe tín hiệu tắt từ OS (Ctrl+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block cho đến khi nhận được tín hiệu
	<-quit
	log.Println("Shutting down server...")

	// Tạo một context với timeout để cho server 5 giây hoàn thành các request hiện tại
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Gọi Shutdown để tắt server một cách duyên dáng
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
		return err
	}

	log.Println("Server exiting")
	return nil
}
