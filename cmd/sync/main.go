package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/zuxt268/homing/internal/di"
	"github.com/zuxt268/homing/internal/infrastructure/database"
	"github.com/zuxt268/homing/internal/infrastructure/driver"
)

func main() {
	db, err := database.NewDB()
	if err != nil {
		panic(err)
	}
	httpClient := &http.Client{Timeout: time.Second * 10}
	httpDriver := driver.NewClient(httpClient)

	customerUsecase := di.NewCustomerUsecase(httpDriver, db)

	err = customerUsecase.SyncAll(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully synced all")
}
