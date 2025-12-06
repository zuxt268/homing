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
)

func Run() {

	db, err := database.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	httpClient := &http.Client{Timeout: time.Minute * 5}
	httpDriver := driver.NewClient(httpClient)

	// GbpAdapter初期化
	credentialsData, err := os.ReadFile(config.Env.GoogleCredentialPath)
	if err != nil {
		log.Fatal("Failed to read credentials file:", err)
	}
	gbpAdapter, err := di.NewGbpAdapter(credentialsData)
	if err != nil {
		log.Fatal("Failed to initialize GBP adapter:", err)
	}

	e := echo.New()

	// ミドルウェア設定
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())

	// ハンドラー初期化
	apiHandler := di.NewHandler(httpDriver, db, gbpAdapter)

	// Swagger ルート
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	api := e.Group("/api")
	api.GET("/healthcheck", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
	api.POST("/sync", apiHandler.SyncAll)
	api.POST("/sync/:id", apiHandler.SyncOne)
	api.POST("/token", apiHandler.SaveToken)
	api.GET("/token", apiHandler.GetToken)
	api.POST("/token/check", apiHandler.CheckToken)

	api.GET("/wordpress-instagram", apiHandler.GetWordpressInstagramList)
	api.GET("/wordpress-instagram/:id", apiHandler.GetWordpressInstagram)
	api.POST("/wordpress-instagram", apiHandler.CreateWordpressInstagram)
	api.PUT("/wordpress-instagram/:id", apiHandler.UpdateWordpressInstagram)
	api.DELETE("/wordpress-instagram/:id", apiHandler.DeleteWordpressInstagram)

	api.GET("/google-business", apiHandler.GetGoogleBusinessList)
	api.GET("/google-business/fetch", apiHandler.FetchGoogleBusinessList)

	api.POST("/business-instagram", apiHandler.CreateBusinessInstagram)
	api.POST("/business-instagram", apiHandler.CreateBusinessInstagram)
	api.PUT("/business-instagram/:id", apiHandler.UpdateBusinessInstagram)

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
