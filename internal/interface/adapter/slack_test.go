package adapter

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zuxt268/homing/internal/domain"
	"github.com/zuxt268/homing/internal/infrastructure/driver"
	"github.com/zuxt268/homing/internal/interface/dto/external"
)

func TestSlackSendMessage(t *testing.T) {
	t.Skip()

	httpClient := &http.Client{}
	client := driver.NewClient(httpClient)
	s := NewSlack(client)
	err := s.SendMessage(context.Background(), external.SlackRequest{
		Text:      "TEST",
		Username:  "neko",
		IconEmoji: ":cat:",
	})
	assert.NoError(t, err)

	wi := domain.WordpressInstagram{}
	err = s.Alert(context.Background(), errors.New("wao").Error(), wi)
	assert.NoError(t, err)
}
