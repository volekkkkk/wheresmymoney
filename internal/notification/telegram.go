package notification

import (
	"github.com/mymmrac/telego"
)

type TelegramProvider struct {
	token string
}

func (*TelegramProvider) NewTelegramProvider(token string) (*telego.Bot, error) {
	return telego.NewBot(token, telego.WithDefaultDebugLogger())
}

func (*TelegramProvider) Send(message string) error {
	return nil
}
