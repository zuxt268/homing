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

	// S3Adapter初期化
	s3Adapter, err := di.NewS3Adapter(config.Env.S3Bucket, config.Env.S3Region, config.Env.S3Prefix)
	if err != nil {
		log.Fatal("Failed to initialize S3 adapter:", err)
	}

	e := echo.New()

	// ミドルウェア設定
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())

	// ハンドラー初期化
	apiHandler := di.NewHandler(httpDriver, db, gbpAdapter, s3Adapter)

	// Swagger ルート
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	api := e.Group("/api")
	api.GET("/healthcheck", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
	api.POST("/sync/wordpress-instagram", apiHandler.SyncAllWordpressInstagram)
	api.POST("/sync/wordpress-instagram/:id", apiHandler.SyncOneWordpressInstagram)

	api.POST("/sync/business-instagram", apiHandler.SyncAllGoogleBusinessInstagram)
	api.POST("/sync/business-instagram/:id", apiHandler.SyncOneGoogleBusinessInstagram)

	api.POST("/token", apiHandler.SaveToken)
	api.GET("/token", apiHandler.GetToken)
	api.POST("/token/check", apiHandler.CheckToken)

	api.GET("/wordpress-instagram/count", apiHandler.GetWordpressInstagramCount)
	api.GET("/wordpress-instagram", apiHandler.GetWordpressInstagramList)
	api.GET("/wordpress-instagram/:id", apiHandler.GetWordpressInstagram)
	api.POST("/wordpress-instagram", apiHandler.CreateWordpressInstagram)
	api.PUT("/wordpress-instagram/:id", apiHandler.UpdateWordpressInstagram)
	api.DELETE("/wordpress-instagram/:id", apiHandler.DeleteWordpressInstagram)

	api.GET("/google-business", apiHandler.GetGoogleBusinessList)
	api.POST("/google-business/fetch", apiHandler.FetchGoogleBusinessList)

	api.GET("/business-instagram", apiHandler.GetBusinessInstagramList)
	api.GET("/business-instagram/:id", apiHandler.GetBusinessInstagram)
	api.POST("/business-instagram", apiHandler.CreateBusinessInstagram)
	api.PUT("/business-instagram/:id", apiHandler.UpdateBusinessInstagram)
	api.DELETE("/business-instagram/:id", apiHandler.DeleteBusinessInstagram)

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
