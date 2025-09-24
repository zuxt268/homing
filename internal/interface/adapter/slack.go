package adapter

import (
	"context"
	"fmt"
	"strings"

	"github.com/zuxt268/homing/internal/config"
	"github.com/zuxt268/homing/internal/domain"
	"github.com/zuxt268/homing/internal/infrastructure/driver"
	"github.com/zuxt268/homing/internal/interface/dto/external"
)

type Slack interface {
	Alert(ctx context.Context, msg string, customer domain.Customer) error
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

const template = `[%s]
顧客: id=%d, name=%s`

func (s *slack) Alert(ctx context.Context, msg string, customer domain.Customer) error {
	sb := strings.Builder{}
	sb.WriteString("｀｀｀")
	sb.WriteString("<@U04P797HYPM>\n")
	sb.WriteString(fmt.Sprintf(template, strings.TrimSpace(msg), customer.ID, customer.Name))
	sb.WriteString("｀｀｀")
	return s.SendMessage(ctx, external.SlackRequest{
		Text:      sb.String(),
		Username:  "homing",
		IconEmoji: ":cat:",
	})
}

func (s *slack) SendMessage(ctx context.Context, payload external.SlackRequest) error {
	_, err := s.httpDriver.Post(ctx, config.Env.SlackWebhookUrl, payload, map[string]string{
		"Content-Type": "application/json",
	})
	return err
}
