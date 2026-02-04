package di

import (
	"github.com/zuxt268/homing/internal/infrastructure/driver"
	"github.com/zuxt268/homing/internal/interface/adapter"
	"github.com/zuxt268/homing/internal/interface/handler"
	"github.com/zuxt268/homing/internal/interface/repository"
	"github.com/zuxt268/homing/internal/usecase"
	"gorm.io/gorm"
)

func NewCustomerRepository(db *gorm.DB) repository.CustomerRepository {
	return repository.NewCustomerRepository(db)
}

func NewPostRepository(db *gorm.DB) repository.PostRepository {
	return repository.NewPostRepository(db)
}

func NewInstagramAdapter(httpClient driver.HttpDriver) adapter.InstagramAdapter {
	return adapter.NewInstagramAdapter(httpClient)
}

func NewSlack(httpDriver driver.HttpDriver) adapter.Slack {
	return adapter.NewSlack(httpDriver)
}

func NewWordpressAdapter(httpDriver driver.HttpDriver) adapter.WordpressAdapter {
	return adapter.NewWordpressAdapter(httpDriver)
}

func NewWordpressInstagramRepository(db *gorm.DB) repository.WordpressInstagramRepository {
	return repository.NewWordpressInstagramRepository(db)
}

func NewTokenRepository(db *gorm.DB) repository.TokenRepository {
	return repository.NewTokenRepository(db)
}

func NewGoogleBusinessRepository(db *gorm.DB) repository.GoogleBusinessRepository {
	return repository.NewGoogleBusinessRepository(db)
}

func NewBusinessInstagramRepository(db *gorm.DB) repository.BusinessInstagramRepository {
	return repository.NewBusinessInstagramRepository(db)
}

func NewGooglePostRepository(db *gorm.DB) repository.GooglePostRepository {
	return repository.NewGooglePostRepository(db)
}

func NewGbpAdapter(credentialsData []byte) (adapter.GbpAdapter, error) {
	return adapter.NewGbpAdapter(credentialsData)
}

func NewFileDownloader() adapter.FileDownloader {
	return adapter.NewFileDownloader()
}

func NewS3Adapter(bucket, region, prefix string) (adapter.S3Adapter, error) {
	return adapter.NewS3Adapter(bucket, region, prefix)
}

func NewCustomerUsecase(httpDriver driver.HttpDriver, db *gorm.DB, gbpAdapter adapter.GbpAdapter, s3Adapter adapter.S3Adapter) usecase.CustomerUsecase {
	return usecase.NewCustomerUsecase(
		NewInstagramAdapter(httpDriver),
		NewSlack(httpDriver),
		NewWordpressAdapter(httpDriver),
		gbpAdapter,
		NewPostRepository(db),
		NewWordpressInstagramRepository(db),
		NewTokenRepository(db),
		NewBusinessInstagramRepository(db),
		NewGooglePostRepository(db),
		s3Adapter,
	)
}

func NewTokenUsecase(httpDriver driver.HttpDriver, db *gorm.DB) usecase.TokenUsecase {
	return usecase.NewTokenUsecase(
		NewInstagramAdapter(httpDriver),
		NewSlack(httpDriver),
		NewTokenRepository(db),
	)
}

func NewWordpressInstagramUsecase(httpDriver driver.HttpDriver, db *gorm.DB) usecase.WordpressInstagramUsecase {
	return usecase.NewWordpressInstagramUsecase(
		NewWordpressInstagramRepository(db),
		NewTokenRepository(db),
		NewPostRepository(db),
		NewInstagramAdapter(httpDriver),
		NewWordpressAdapter(httpDriver),
	)
}

func NewBusinessInstagramUsecase(httpDriver driver.HttpDriver, db *gorm.DB, gbpAdapter adapter.GbpAdapter) usecase.BusinessInstagramUsecase {
	return usecase.NewBusinessInstagramUsecase(
		NewGoogleBusinessRepository(db),
		NewTokenRepository(db),
		NewBusinessInstagramRepository(db),
		NewGooglePostRepository(db),
		NewInstagramAdapter(httpDriver),
		gbpAdapter,
	)
}

func NewHandler(httpDriver driver.HttpDriver, db *gorm.DB, gbpAdapter adapter.GbpAdapter, s3Adapter adapter.S3Adapter) handler.APIHandler {
	return handler.NewAPIHandler(
		NewCustomerUsecase(httpDriver, db, gbpAdapter, s3Adapter),
		NewTokenUsecase(httpDriver, db),
		NewWordpressInstagramUsecase(httpDriver, db),
		NewBusinessInstagramUsecase(httpDriver, db, gbpAdapter),
	)
}
