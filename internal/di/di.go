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

func NewFileDownloader() adapter.FileDownloader {
	return adapter.NewFileDownloader()
}

func NewCustomerUsecase(httpDriver driver.HttpDriver, db *gorm.DB) usecase.CustomerUsecase {
	return usecase.NewCustomerUsecase(
		NewFileDownloader(),
		NewInstagramAdapter(httpDriver),
		NewSlack(httpDriver),
		NewWordpressAdapter(httpDriver),
		NewPostRepository(db),
		NewWordpressInstagramRepository(db),
		NewTokenRepository(db),
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
		NewInstagramAdapter(httpDriver),
		NewWordpressAdapter(httpDriver),
	)
}

func NewHandler(httpDriver driver.HttpDriver, db *gorm.DB) handler.APIHandler {
	return handler.NewAPIHandler(
		NewCustomerUsecase(httpDriver, db),
		NewTokenUsecase(httpDriver, db),
		NewWordpressInstagramUsecase(httpDriver, db),
	)
}
