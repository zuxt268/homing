//go:build wireinject
// +build wireinject

package di

import (
	"net/http"

	"github.com/google/wire"
	"github.com/zuxt268/homing/internal/infrastructure/driver"
	"github.com/zuxt268/homing/internal/interface/adapter"
	"github.com/zuxt268/homing/internal/interface/repository"
	"github.com/zuxt268/homing/internal/usecase"
)

// Provider sets
var DriverSet = wire.NewSet(
	wire.Struct(new(http.Client)),
	driver.NewClient,
)

var AdapterSet = wire.NewSet(
	adapter.NewInstagramAdapter,
	adapter.NewWordpressAdapter,
	adapter.NewFileDownloader,
	adapter.NewSlack,
)

var RepositorySet = wire.NewSet(
	repository.NewCustomerRepository,
	repository.NewPostRepository,
)

var UsecaseSet = wire.NewSet(
	usecase.NewCustomerUsecase,
)

func InitializeCustomerUsecase() (usecase.CustomerUsecase, error) {
	wire.Build(
		DriverSet,
		AdapterSet,
		RepositorySet,
		UsecaseSet,
	)
	return nil, nil
}
