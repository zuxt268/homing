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
	Alert(ctx context.Context, msg string, wi domain.WordpressInstagram) error
	SendMessage(ctx context.Context, payload external.SlackRequest) error
	SendTokenExpired(ctx context.Context) error
	Success(ctx context.Context, wi *domain.WordpressInstagram, wordpressUrl, instagramUrl string) error
}

type slack struct {
	httpDriver             driver.HttpDriver
	noticeWebAppChannelUrl string
	prjARootChannelUrl     string
}

func NewSlack(httpDriver driver.HttpDriver) Slack {
	return &slack{
		httpDriver:             httpDriver,
		noticeWebAppChannelUrl: config.Env.NoticeWebAppChannelUrl,
		prjARootChannelUrl:     config.Env.PrjARootChannelUrl,
	}
}

const template = `[%s]
顧客: id=%d, name=%s`

func (s *slack) Alert(ctx context.Context, msg string, wi domain.WordpressInstagram) error {
	sb := strings.Builder{}
	sb.WriteString("｀｀｀")
	sb.WriteString("<@U04P797HYPM>\n")
	sb.WriteString(fmt.Sprintf(template, strings.TrimSpace(msg), wi.ID, wi.Name))
	sb.WriteString("｀｀｀")
	return s.noticeWebAppChannel(ctx, external.SlackRequest{
		Text:      sb.String(),
		Username:  "homing",
		IconEmoji: ":cat:",
	})
}

func (s *slack) SendMessage(ctx context.Context, payload external.SlackRequest) error {
	return s.noticeWebAppChannel(ctx, payload)
}

func (s *slack) SendTokenExpired(ctx context.Context) error {
	return s.prjARootChannel(ctx, external.SlackRequest{
		Text:      "トークンの有効期限が近づいています",
		Username:  "homing",
		IconEmoji: ":heavy_exclamation:",
	})
}

func (s *slack) noticeWebAppChannel(ctx context.Context, payload external.SlackRequest) error {
	_, err := s.httpDriver.Post(ctx, config.Env.NoticeWebAppChannelUrl, payload, map[string]string{
		"Content-Type": "application/json",
	})
	return err
}

func (s *slack) prjARootChannel(ctx context.Context, payload external.SlackRequest) error {
	_, err := s.httpDriver.Post(ctx, config.Env.PrjARootChannelUrl, payload, map[string]string{
		"Content-Type": "application/json",
	})
	return err
}

const templateSuccess = `[SYSTEM USER]
id: %d
name: %s
wordpress: %s
instagram: %s
`

func (s *slack) Success(ctx context.Context, wi *domain.WordpressInstagram, wordpressUrl, instagramUrl string) error {
	sb := strings.Builder{}
	sb.WriteString("｀｀｀")
	sb.WriteString(fmt.Sprintf(templateSuccess, wi.ID, wi.Name, wordpressUrl, instagramUrl))
	sb.WriteString("｀｀｀")
	return s.noticeWebAppChannel(ctx, external.SlackRequest{
		Text:      sb.String(),
		Username:  "homing",
		IconEmoji: ":cat:",
	})
}
