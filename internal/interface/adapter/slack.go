package adapter

import (
	"context"

	"github.com/zuxt268/homing/internal/config"
	"github.com/zuxt268/homing/internal/domain/entity"
	"github.com/zuxt268/homing/internal/infrastructure/driver"
	"github.com/zuxt268/homing/internal/interface/dto/external"
)

type Slack interface {
	Alert(ctx context.Context, msg string, customer entity.Customer) error
	SendMessage(ctx context.Context, payload external.SlackRequest) error
}

type slack struct {
	httpDriver driver.HttpDriver
}

func NewSlack(httpDriver driver.HttpDriver) Slack {
	return &slack{
		httpDriver: httpDriver,
	}
}

func (s *slack) Alert(ctx context.Context, msg string, customer entity.Customer) error {
	return nil
}

func (s *slack) SendMessage(ctx context.Context, payload external.SlackRequest) error {
	_, err := s.httpDriver.Post(ctx, config.Env.SlackWebhookUrl, payload, map[string]string{
		"Content-Type": "application/json",
	})
	return err
}
