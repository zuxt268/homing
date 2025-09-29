package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/docker/go-connections/nat"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/zuxt268/homing/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB() (*gorm.DB, error) {
	// 設定が存在しない場合はモックDBを返す（開発用）
	if config.Env.DBHost == "" {
		log.Println("Warning: Database configuration not found, using mock connection")
		return nil, fmt.Errorf("database configuration not available")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Env.DBUser,
		config.Env.DBPassword,
		config.Env.DBHost,
		config.Env.DBPort,
		config.Env.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	return db, nil
}

func NewTestContainerDBClient() (*gorm.DB, func()) {
	ctx := context.Background()

	// MySQL testcontainer を起動
	req := testcontainers.ContainerRequest{
		Image:        "mysql:8.0",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "testpass",
			"MYSQL_DATABASE":      "testdb",
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("3306/tcp"),
			wait.ForSQL("3306/tcp", "mysql", func(host string, port nat.Port) string {
				return fmt.Sprintf("root:testpass@tcp(%s:%s)/testdb?charset=utf8mb4&parseTime=True&loc=Local", host, port.Port())
			}).WithStartupTimeout(120*time.Second).WithPollInterval(2*time.Second),
		),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to start MySQL container: %v", err))
	}

	// コンテナの接続情報を取得
	host, err := container.Host(ctx)
	if err != nil {
		panic(fmt.Sprintf("Failed to get container host: %v", err))
	}

	port, err := container.MappedPort(ctx, "3306")
	if err != nil {
		panic(fmt.Sprintf("Failed to get container port: %v", err))
	}

	// DSN作成
	dsn := fmt.Sprintf("root:testpass@tcp(%s:%s)/testdb?charset=utf8mb4&parseTime=True&loc=Local", host, port.Port())

	// GORM DB接続
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		_ = container.Terminate(ctx)
		panic(fmt.Sprintf("Failed to connect to test database: %v", err))
	}

	// マイグレーション実行
	sqlDB, err := db.DB()
	if err != nil {
		_ = container.Terminate(ctx)
		panic(fmt.Sprintf("Failed to get SQL DB: %v", err))
	}

	migrationDir := filepath.Join(GetProjectRoot(), "migrations")
	migrations := &migrate.FileMigrationSource{
		Dir: migrationDir,
	}
	_, err = migrate.Exec(sqlDB, "mysql", migrations, migrate.Up)
	if err != nil {
		_ = container.Terminate(ctx)
		panic(fmt.Sprintf("Failed to run migrations: %v", err))
	}

	fmt.Println("TestContainer MySQL initialized with migrations")

	// クリーンアップ関数を返す
	cleanup := func() {
		_ = container.Terminate(ctx)
	}

	return db, cleanup
}

func GetProjectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	// go.modが見つかるまで親ディレクトリに遡る
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break // ルートディレクトリに到達
		}
		dir = parent
	}
	return ""
}
