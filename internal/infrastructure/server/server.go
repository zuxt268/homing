package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zuxt268/homing/internal/config"
	"github.com/zuxt268/homing/internal/di"
	"github.com/zuxt268/homing/internal/interface/handler"
)

func Run() {
	// DI コンテナ初期化
	customerUsecase, err := di.InitializeCustomerUsecase()
	if err != nil {
		log.Fatal("Failed to initialize dependencies:", err)
	}

	e := echo.New()

	// ミドルウェア設定
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())

	// ハンドラー初期化
	apiHandler := handler.NewAPIHandler(customerUsecase)

	// ルーティング
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	// API ルート設定
	api := e.Group("/api")
	api.POST("/sync", apiHandler.SyncAll)
	api.POST("/sync/:customer_id", apiHandler.SyncOne)

	srv := &http.Server{
		Addr:    config.Env.Address,
		Handler: e,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("listen: %s\n", err)
		}
	}()

	fmt.Println()
	fmt.Println("**********************")
	fmt.Println("homing server started!")
	fmt.Println("**********************")
	fmt.Println()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
