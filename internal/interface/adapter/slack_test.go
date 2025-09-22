package adapter

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zuxt268/homing/internal/infrastructure/driver"
	"github.com/zuxt268/homing/internal/interface/dto/external"
)

func TestSlackSendMessage(t *testing.T) {
	httpClient := &http.Client{}
	client := driver.NewClient(httpClient)
	s := NewSlack(client)
	err := s.SendMessage(context.Background(), external.SlackRequest{
		Text:      "TEST",
		Username:  "neko",
		IconEmoji: ":cat:",
	})
	assert.NoError(t, err)
}
