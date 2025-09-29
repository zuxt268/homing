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
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/zuxt268/homing/docs" // Swaggerドキュメント用
	"github.com/zuxt268/homing/internal/config"
	"github.com/zuxt268/homing/internal/di"
	"github.com/zuxt268/homing/internal/infrastructure/database"
	"github.com/zuxt268/homing/internal/infrastructure/driver"
	"github.com/zuxt268/homing/internal/interface/handler"
)

func Run() {
	
	db, err := database.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	httpClient := &http.Client{Timeout: time.Second * 10}
	httpDriver := driver.NewClient(httpClient)

	customerUsecase := di.NewCustomerUsecase(httpDriver, db)

	e := echo.New()

	// ミドルウェア設定
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())

	// ハンドラー初期化
	apiHandler := handler.NewAPIHandler(customerUsecase)

	// Swagger ルート
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// API ルート設定
	api := e.Group("/api")
	api.POST("/sync", apiHandler.SyncAll)
	api.POST("/sync/:customer_id", apiHandler.SyncOne)
	api.GET("/customers/:customer_id", apiHandler.GetCustomer)

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
